package asyncevents

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/IBM/sarama"
)

var errTest error = errors.New("test error")

func TestSaramaAsyncEventsConsumerLifecycle(s *testing.T) {
	consumerGroup := &nullSaramaConsumerGroup{}
	topics := []string{"test"}

	s.Run("successful start and stop", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		handler := &nullSaramaConsumerGroupHandler{}
		eventsConsumer := NewSaramaConsumerGroupManager(consumerGroup, handler, topics...)
		err := launchStart(ctx, t, eventsConsumer)
		if !errors.Is(err, nil) {
			t.Errorf("expected nil error, got %v", err)
		}
	})

	s.Run("reports errors (that aren't context.Canceled)", func(t *testing.T) {
		consumerGroup := &erroringSaramaConsumerGroup{err: errTest}
		handler := &nullSaramaConsumerGroupHandler{}
		eventsConsumer := NewSaramaConsumerGroupManager(consumerGroup, handler, topics...)
		err := launchStart(context.Background(), t, eventsConsumer)
		if !errors.Is(err, errTest) {
			t.Errorf("expected %s, got %v", errTest, err)
		}
	})
}

func TestSaramaConsumerGroupHandler(s *testing.T) {
	s.Run("works as expected", func(t *testing.T) {
		testConsumer := &sleepingSaramaMessageConsumer{time.Nanosecond, nil}
		testSession := &nullSaramaConsumerGroupSession{}
		messages := make(chan *sarama.ConsumerMessage)
		testClaim := newTestSaramaConsumerGroupClaim(messages)
		go func() { messages <- &sarama.ConsumerMessage{}; close(messages) }()
		handler := NewSaramaConsumerGroupHandler(slog.Default(), testConsumer, 0)

		err := handler.ConsumeClaim(testSession, testClaim)
		if !errors.Is(err, nil) {
			t.Errorf("expected ConsumeClaim to return nil, got %v", err)
		}
	})

	s.Run("returns errors", func(t *testing.T) {
		testConsumer := newCountingSaramaMessageConsumer(errTest)
		testSession := &nullSaramaConsumerGroupSession{}
		messages := make(chan *sarama.ConsumerMessage)
		testClaim := newTestSaramaConsumerGroupClaim(messages)
		go func() { messages <- &sarama.ConsumerMessage{}; close(messages) }()
		handler := NewSaramaConsumerGroupHandler(slog.Default(), testConsumer, time.Nanosecond)

		err := handler.ConsumeClaim(testSession, testClaim)
		if !errors.Is(err, errTest) {
			t.Errorf("expected ConsumeClaim to return %s, got %v", errTest, err)
		}
	})

	s.Run("enforces a deadline", func(t *testing.T) {
		testConsumer := &sleepingSaramaMessageConsumer{time.Second, nil}
		testSession := &nullSaramaConsumerGroupSession{}
		messages := make(chan *sarama.ConsumerMessage)
		testClaim := newTestSaramaConsumerGroupClaim(messages)
		go func() { messages <- &sarama.ConsumerMessage{}; close(messages) }()
		handler := NewSaramaConsumerGroupHandler(slog.Default(), testConsumer, time.Nanosecond)

		err := handler.ConsumeClaim(testSession, testClaim)
		if !errors.Is(testConsumer.consumeError, context.DeadlineExceeded) {
			t.Errorf("expected ConsumeClaim to return %s, got %v", context.DeadlineExceeded, err)
		}
	})
}

func TestNTimesRetryingConsumer(s *testing.T) {
	var testTimes = 3

	s.Run("retries N times", func(t *testing.T) {
		testConsumer := newCountingSaramaMessageConsumer(errTest)
		c := &NTimesRetryingConsumer{
			Times:    testTimes,
			Consumer: testConsumer,
			Delay:    DelayNone,
		}
		ctx := context.Background()
		err := c.Consume(ctx, nil, nil)
		if !errors.Is(err, ErrRetriesLimitExceeded) {
			t.Errorf("expected %s, got %v", ErrRetriesLimitExceeded, err)
		}
		if testConsumer.Count != testTimes {
			t.Errorf("expected %d tries, got %d", testTimes, testConsumer.Count)
		}
	})

	s.Run("retries when the context deadline is exceeded", func(t *testing.T) {
		testConsumer := newCountingSaramaMessageConsumer(context.DeadlineExceeded)
		c := &NTimesRetryingConsumer{
			Times:    testTimes,
			Consumer: testConsumer,
			Delay:    DelayNone,
		}
		ctx := context.Background()
		err := c.Consume(ctx, nil, nil)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("expected %s, got %v", context.DeadlineExceeded, err)
		}
		if testConsumer.Count != testTimes {
			t.Errorf("expected %d tries, got %d", testTimes, testConsumer.Count)
		}
	})

	s.Run("returns nil when the context is canceled", func(t *testing.T) {
		testConsumer := newCountingSaramaMessageConsumer(context.Canceled)
		c := &NTimesRetryingConsumer{
			Times:    testTimes,
			Consumer: testConsumer,
			Delay:    DelayNone,
		}
		ctx := context.Background()
		err := c.Consume(ctx, nil, nil)
		if !errors.Is(err, nil) {
			t.Errorf("expected nil error, got %s", err)
		}
		if testConsumer.Count >= testTimes {
			t.Errorf("expected < %d tries, got %d", testTimes, testConsumer.Count)
		}
	})

	s.Run("doesn't retry on successful consumption", func(t *testing.T) {
		testConsumer := newCountingSaramaMessageConsumer(nil)
		c := &NTimesRetryingConsumer{
			Times:    testTimes,
			Consumer: testConsumer,
			Delay:    DelayNone,
		}
		ctx := context.Background()
		err := c.Consume(ctx, nil, nil)
		if !errors.Is(err, nil) {
			t.Errorf("expected nil error, got %v", err)
		}
		if testConsumer.Count != 1 {
			t.Errorf("expected 1 try, got %d", testConsumer.Count)
		}
	})
}

func ExampleDelayNone() {
	fmt.Println(
		DelayNone(0),
		DelayNone(10),
		DelayNone(20),
		DelayNone(100),
		DelayNone(1024),
	)
	// Output: 0s 0s 0s 0s 0s
}

func ExampleDelayExponentialBinary() {
	fmt.Println(
		DelayExponentialBinary(0),
		DelayExponentialBinary(1),
		DelayExponentialBinary(2),
		DelayExponentialBinary(3),
		DelayExponentialBinary(4),
		DelayExponentialBinary(5),
		DelayExponentialBinary(6),
		DelayExponentialBinary(7),
		DelayExponentialBinary(8),
	)
	// Output: 1s 2s 4s 8s 16s 32s 1m4s 2m8s 4m16s
}

func ExampleFib() {
	fmt.Println(
		Fib(0),
		Fib(1),
		Fib(2),
		Fib(10),
		Fib(20),
	)
	// Output: 0 1 1 55 6765
}

// launchStart is a helper for calling Start in a goroutine.
//
// It uses a channel to know that the goroutine has seen some amount of CPU
// time, which isn't guaranteed to alleviate the race of calling Start, but in
// practice seems to be sufficient. Running with -count 10000 had 0 failures.
func launchStart(ctx context.Context, t testing.TB, ec *SaramaConsumerGroupManager) (err error) {
	t.Helper()
	runReturned := make(chan error)
	go func() {
		defer close(runReturned)
		runReturned <- ec.Run(ctx)
	}()
	return <-runReturned
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

type erroringSaramaConsumerGroup struct {
	nullSaramaConsumerGroup
	err error
}

func (g *erroringSaramaConsumerGroup) Consume(_ context.Context, _ []string, _ sarama.ConsumerGroupHandler) error {
	return g.err
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
	err      error
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
	if s.context == nil {
		s.context = context.Background()
	}
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
	consumeError  error
}

func (c *sleepingSaramaMessageConsumer) Consume(ctx context.Context,
	session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) (err error) {
	defer func(err *error) {
		if err != nil {
			c.consumeError = *err
		}
	}(&err)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(c.sleepDuration):
		return nil
	}
}

type countingSaramaMessageConsumer struct {
	err   error
	Count int
}

func newCountingSaramaMessageConsumer(err error) *countingSaramaMessageConsumer {
	return &countingSaramaMessageConsumer{
		err: err,
	}
}

func (c *countingSaramaMessageConsumer) Consume(ctx context.Context,
	_ sarama.ConsumerGroupSession, _ *sarama.ConsumerMessage) error {
	c.Count++
	return c.err
}
