package asyncevents

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

// SaramaAsyncEventsConsumer consumes Kafka messages for asynchronous event
// handling.
type SaramaAsyncEventsConsumer struct {
	handler       sarama.ConsumerGroupHandler
	consumerGroup sarama.ConsumerGroup
	topics        []string
}

func NewSaramaEventsConsumer(consumerGroup sarama.ConsumerGroup,
	handler sarama.ConsumerGroupHandler, topics ...string) *SaramaAsyncEventsConsumer {

	return &SaramaAsyncEventsConsumer{
		consumerGroup: consumerGroup,
		handler:       handler,
		topics:        topics,
	}
}

// Run consuming Kafka messages.
//
// Run is stopped by its context being canceled, so when that happens, it
// returns nil.
func (p *SaramaAsyncEventsConsumer) Run(ctx context.Context) (err error) {
	defer cancelReturnsNil(&err)

	for {
		err := p.consumerGroup.Consume(ctx, p.topics, p.handler)
		if err != nil {
			return err
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
	}
}

// cancelReturnsNil helps return nil no matter from where in a function defer
// is called. It is meant to be called via defer.
func cancelReturnsNil(err *error) {
	if err != nil && errors.Is(*err, context.Canceled) {
		*err = nil
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

const (
	// DefaultMessageConsumptionTimeout is the default time to allow
	// SaramaMessageConsumer.Consume to work before aborting.
	DefaultMessageConsumptionTimeout = time.Minute
)

// Setup implements sarama.ConsumerGroupHandler.
func (h *saramaConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup implements sarama.ConsumerGroupHandler.
func (h *saramaConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (h *saramaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		err := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), h.consumerTimeout)
			defer cancel()
			return h.consumer.Consume(ctx, session, message)
		}()
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			log.Print(err)
		case !errors.Is(err, nil):
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

var ErrRetriesLimitExceeded = errors.New("retry limit exceeded")

// NTimesRetryingConsumer enhances a SaramaMessageConsumer with a finite
// number of immediate retries.
type NTimesRetryingConsumer struct {
	Times    int
	Consumer SaramaMessageConsumer
}

func (c *NTimesRetryingConsumer) Consume(ctx context.Context,
	session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) error {

	var joinedErrors error
	var tries = 0
	for tries < c.Times {
		err := c.Consumer.Consume(ctx, session, message)
		switch {
		case errors.Is(err, context.Canceled),
			errors.Is(err, context.DeadlineExceeded):
			return err
		case errors.Is(err, nil):
			return nil
		default:
			tries++
			joinedErrors = errors.Join(joinedErrors, err)
			log.Printf("retrying after consumer error: %s", err)
		}
	}

	return errors.Join(joinedErrors, c.retryLimitError())
}

func (c *NTimesRetryingConsumer) retryLimitError() error {
	return fmt.Errorf("%w (%d)", ErrRetriesLimitExceeded, c.Times)
}
