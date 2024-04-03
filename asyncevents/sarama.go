package asyncevents

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

// SaramaAsyncEventsConsumer consumes Kafka messages for asynchronous event
// handling.
type SaramaAsyncEventsConsumer struct {
	handler       sarama.ConsumerGroupHandler
	consumerGroup sarama.ConsumerGroup
	topics        []string

	// stopTimeout controls the time that Stop waits for the consumer to start
	// before aborting.
	stopTimeout time.Duration
	// cancelFunc ensures that a consumer isn't stopped before it starts.
	//
	// Calling Stop will read from this channel, which must not have anything
	// to read until Start is called.
	cancelFunc chan context.CancelFunc
	startOnce  sync.Once
}

const (
	DefaultSaramaAsyncEventConsumerStopTimeout = time.Second
)

func NewSaramaEventsConsumer(consumerGroup sarama.ConsumerGroup,
	handler sarama.ConsumerGroupHandler, topics ...string) *SaramaAsyncEventsConsumer {

	return &SaramaAsyncEventsConsumer{
		cancelFunc:    make(chan context.CancelFunc, 1),
		consumerGroup: consumerGroup,
		handler:       handler,
		stopTimeout:   DefaultSaramaAsyncEventConsumerStopTimeout,
		topics:        topics,
	}
}

// Start consuming Kafka messages.
//
// If Start begins consuming Kafka messages, then it will block until Stop is
// called. If the consumer is already running, Start returns
// ErrAlreadyStarted.
//
// The caller is responsible for calling Stop() to cleanup resources.
func (p *SaramaAsyncEventsConsumer) Start(ctx context.Context) (err error) {
	ctx, cancel, err := p.start(ctx)
	if err != nil {
		return err
	}
	defer cancel()

	defer func() {
		if errors.Is(err, context.Canceled) {
			err = nil
		}
	}()

	for {
		//log.Printf("Start(): loop")
		err := p.consumerGroup.Consume(ctx, p.topics, p.handler)
		if err != nil {
			return err
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
	}
}

// start ensures that the events consumer is started only once.
func (p *SaramaAsyncEventsConsumer) start(ctx context.Context) (context.Context, context.CancelFunc, error) {
	var newCtx context.Context
	var cancel context.CancelFunc
	var err error = ErrAlreadyStarted

	p.startOnce.Do(func() {
		newCtx, cancel = context.WithCancel(ctx)
		select {
		case p.cancelFunc <- cancel:
			close(p.cancelFunc)
			err = nil
		case <-time.After(time.Second):
			err = ErrStartBlocked
		}
	})

	return newCtx, cancel, err
}

var (
	ErrAlreadyStarted  = errors.New("already started")
	ErrStartBlocked    = errors.New("Start() is blocked, this should _never_ happen")
	ErrStopBeforeStart = errors.New("Stop() called before Start()")
)

// Stop consuming Kafka messages.
//
// It returns an error if Stop is called before Start(). This should probably
// take a context.
//
// TODO: this should probably take a context.
func (p *SaramaAsyncEventsConsumer) Stop() error {
	select {
	case cancel := <-p.cancelFunc:
		if cancel != nil {
			cancel()
		}
		return nil
	case <-time.After(p.stopTimeout):
		return ErrStopBeforeStart
	}
}

// saramaConsumerGroupHandler implements sarama.ConsumerGroupHandler.
type saramaConsumerGroupHandler struct {
	consumer        SaramaMessageConsumer
	consumerTimeout time.Duration
}

func NewSaramaConsumerGroupHandler(consumer SaramaMessageConsumer, timeout time.Duration) *saramaConsumerGroupHandler {
	if timeout == 0 {
		timeout = DefaultMessageConsumptionTimeout
	}
	return &saramaConsumerGroupHandler{
		consumer:        consumer,
		consumerTimeout: timeout,
	}
}

// Setup implements sarama.ConsumerGroupHandler.
func (h *saramaConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup implements sarama.ConsumerGroupHandler.
func (h *saramaConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

var (
	// DefaultMessageConsumptionTimeout is the default time to allow
	// SaramaMessageConsumer.Consume to work before aborting.
	DefaultMessageConsumptionTimeout = time.Minute
)

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (h *saramaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		err := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), h.consumerTimeout)
			defer cancel()
			return h.consumer.Consume(ctx, session, message)
		}()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				// TODO consult some sort of error handler for what to do next
				log.Printf("timed out")
			}
			return err
		}
	}
	return nil
}

// SaramaMessageConsumer processes Kafka messages.
type SaramaMessageConsumer interface {
	// Consume should process a message.
	//
	// Consume is responsible for marking the message consumed, unless the
	// context is canceled, in which case the caller should retry, or mark the
	// message as appropriate.
	Consume(ctx context.Context, session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) error
}
