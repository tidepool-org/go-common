package asyncevents

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

func TestSaramaAsyncEventsConsumerLifecycle(s *testing.T) {
	consumerGroup := &nullSaramaConsumerGroup{}
	handler := &nullSaramaConsumerGroupHandler{}
	topics := []string{"test"}

	s.Run("successful start and stop", func(t *testing.T) {
		ctx := context.Background()
		eventsConsumer := NewSaramaEventsConsumer(consumerGroup, handler, topics...)

		launchStart(ctx, t, eventsConsumer)

		if err := eventsConsumer.Stop(); !errors.Is(err, nil) {
			t.Errorf("expected Start() to return nil, got %v", err)
		}
	})

	s.Run("calling Stop before Start is an error", func(t *testing.T) {
		eventsConsumer := NewSaramaEventsConsumer(consumerGroup, handler, topics...)
		eventsConsumer.stopTimeout = time.Nanosecond

		if err := eventsConsumer.Stop(); !errors.Is(err, ErrStopBeforeStart) {
			t.Errorf("expected Stop() to return %s, got %v", ErrStopBeforeStart, err)
		}
	})

	s.Run("calling Stop multiple times is not an error", func(t *testing.T) {
		ctx := context.Background()
		eventsConsumer := NewSaramaEventsConsumer(consumerGroup, handler, topics...)

		launchStart(ctx, t, eventsConsumer)

		if err := eventsConsumer.Stop(); !errors.Is(err, nil) {
			t.Errorf("expected Stop() to return nil, got %v", err)
		}
		if err := eventsConsumer.Stop(); !errors.Is(err, nil) {
			t.Errorf("expected Stop() to return nil, got %v", err)
		}
	})

	s.Run("calling Start more than once returns an error", func(t *testing.T) {
		ctx := context.Background()
		eventsConsumer := NewSaramaEventsConsumer(consumerGroup, handler, topics...)

		launchStart(ctx, t, eventsConsumer)

		if err := eventsConsumer.Start(ctx); !errors.Is(err, ErrAlreadyStarted) {
			t.Errorf("expected Start() to return %s, got %v", ErrAlreadyStarted, err)
		}
	})
}

func TestSaramaConsumerGroupHandler(s *testing.T) {
	s.Run("works as expected", func(t *testing.T) {
		testConsumer := &sleepingSaramaMessageConsumer{time.Nanosecond}
		testSession := &nullSaramaConsumerGroupSession{}
		messages := make(chan *sarama.ConsumerMessage)
		testClaim := newTestSaramaConsumerGroupClaim(messages)
		go func() { messages <- &sarama.ConsumerMessage{}; close(messages) }()
		handler := NewSaramaConsumerGroupHandler(testConsumer, 0)

		err := handler.ConsumeClaim(testSession, testClaim)
		if !errors.Is(err, nil) {
			t.Errorf("expected ConsumeClaim to return nil, got %v", err)
		}
	})

	s.Run("enforces its timeout", func(t *testing.T) {
		testConsumer := &sleepingSaramaMessageConsumer{time.Second}
		testSession := &nullSaramaConsumerGroupSession{}
		messages := make(chan *sarama.ConsumerMessage)
		testClaim := newTestSaramaConsumerGroupClaim(messages)
		go func() { messages <- &sarama.ConsumerMessage{}; close(messages) }()
		handler := NewSaramaConsumerGroupHandler(testConsumer, time.Nanosecond)

		err := handler.ConsumeClaim(testSession, testClaim)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("expected ConsumeClaim to return %s, got %v", context.DeadlineExceeded, err)
		}
	})
}

// launchStart is a helper for calling Start in a goroutine.
//
// It uses a channel to know that the goroutine has seen some amount of CPU
// time, which isn't guaranteed to alleviate the race of calling Start, but in
// practice seems to be sufficient. Running with -count 10000 had 0 failures.
func launchStart(ctx context.Context, t testing.TB, ec *SaramaAsyncEventsConsumer) {
	t.Helper()
	launched := make(chan struct{}, 1)
	t.Cleanup(func() { close(launched) })
	go func() {
		launched <- struct{}{}
		if err := ec.Start(ctx); !errors.Is(err, nil) {
			t.Errorf("expected Start() to return nil; got %v", err)
		}
	}()
	<-launched
}

// nullSaramaConsumerGroup is a null/no-op base from which to mock test behavior.
type nullSaramaConsumerGroup struct{}

func (g *nullSaramaConsumerGroup) Errors() <-chan error {
	return nil
}

func (g *nullSaramaConsumerGroup) Close() error {
	return nil
}

func (g *nullSaramaConsumerGroup) Pause(partitions map[string][]int32) {}

func (g *nullSaramaConsumerGroup) Resume(partitions map[string][]int32) {}

func (g *nullSaramaConsumerGroup) PauseAll() {}

func (g *nullSaramaConsumerGroup) ResumeAll() {}

func (g *nullSaramaConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	return nil
}

// nullSaramaConsumerGroupHandler is a no-op base from which to mock test
// behavior.
type nullSaramaConsumerGroupHandler struct {
	t testing.TB
}

func (h *nullSaramaConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *nullSaramaConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *nullSaramaConsumerGroupHandler) ConsumeClaim(_ sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if h.t != nil {
			h.t.Logf("consuming message: %+v", msg)
		}
	}
	return nil
}

type testSaramaConsumerGroupHandler struct {
	*nullSaramaConsumerGroupHandler
	consumed []*sarama.ConsumerMessage
}

func (h *testSaramaConsumerGroupHandler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	for msg := range c.Messages() {
		h.consumed = append(h.consumed, msg)
	}
	return nil
}

type nullSaramaConsumerGroupSession struct {
	context context.Context
}

func (s *nullSaramaConsumerGroupSession) Claims() map[string][]int32 {
	return map[string][]int32{}
}

func (s *nullSaramaConsumerGroupSession) MemberID() string {
	return ""
}

func (s *nullSaramaConsumerGroupSession) GenerationID() int32 {
	return 0
}

func (s *nullSaramaConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
}

func (s *nullSaramaConsumerGroupSession) Commit() {}

func (s *nullSaramaConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
}

func (s *nullSaramaConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {}

func (s *nullSaramaConsumerGroupSession) Context() context.Context {
	return s.context
}

func (c *testSaramaConsumerGroupClaim) Topic() string {
	return ""
}

func (c *testSaramaConsumerGroupClaim) Partition() int32 {
	return 0
}

func (c *testSaramaConsumerGroupClaim) InitialOffset() int64 {
	return 0
}

func (c *testSaramaConsumerGroupClaim) HighWaterMarkOffset() int64 {
	return 0
}

type testSaramaConsumerGroupClaim struct {
	messages <-chan *sarama.ConsumerMessage
}

func newTestSaramaConsumerGroupClaim(messages <-chan *sarama.ConsumerMessage) *testSaramaConsumerGroupClaim {
	return &testSaramaConsumerGroupClaim{
		messages: messages,
	}
}

func (c *testSaramaConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	return c.messages
}

type sleepingSaramaMessageConsumer struct {
	sleepDuration time.Duration
}

func (c *sleepingSaramaMessageConsumer) Consume(ctx context.Context,
	session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(c.sleepDuration):
		return nil
	}
}
