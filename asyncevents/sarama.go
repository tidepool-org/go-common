package asyncevents

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math"
	"time"

	"github.com/IBM/sarama"
)

// SaramaEventsConsumer consumes Kafka messages for asynchronous event
// handling.
type SaramaEventsConsumer struct {
	Handler       sarama.ConsumerGroupHandler
	ConsumerGroup sarama.ConsumerGroup
	Topics        []string
}

func NewSaramaEventsConsumer(consumerGroup sarama.ConsumerGroup,
	handler sarama.ConsumerGroupHandler, topics ...string) *SaramaEventsConsumer {

	return &SaramaEventsConsumer{
		ConsumerGroup: consumerGroup,
		Handler:       handler,
		Topics:        topics,
	}
}

// Run the consumer, to begin consuming Kafka messages.
//
// Run is stopped by its context being canceled. When its context is canceled,
// it returns nil.
func (p *SaramaEventsConsumer) Run(ctx context.Context) (err error) {
	defer canceledContextReturnsNil(&err)

	for {
		err := p.ConsumerGroup.Consume(ctx, p.Topics, p.Handler)
		if err != nil {
			return err
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
	}
}

// canceledContextReturnsNil checks for a context.Canceled error, and when
// found, returns nil instead.
//
// It is meant to be called via defer.
func canceledContextReturnsNil(err *error) {
	if err != nil && errors.Is(*err, context.Canceled) {
		*err = nil
	}
}

// SaramaConsumerGroupHandler implements sarama.ConsumerGroupHandler.
type SaramaConsumerGroupHandler struct {
	Consumer        SaramaMessageConsumer
	ConsumerTimeout time.Duration
}

// NewSaramaConsumerGroupHandler builds a consumer group handler.
//
// A timeout of 0 will use DefaultMessageConsumptionTimeout.
func NewSaramaConsumerGroupHandler(consumer SaramaMessageConsumer, timeout time.Duration) *SaramaConsumerGroupHandler {
	if timeout == 0 {
		timeout = DefaultMessageConsumptionTimeout
	}
	return &SaramaConsumerGroupHandler{
		Consumer:        consumer,
		ConsumerTimeout: timeout,
	}
}

const (
	// DefaultMessageConsumptionTimeout is the default time to allow
	// SaramaMessageConsumer.Consume to work before canceling.
	DefaultMessageConsumptionTimeout = 30 * time.Second
)

// Setup implements sarama.ConsumerGroupHandler.
func (h *SaramaConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return h.Consumer.Setup(session)
}

// Cleanup implements sarama.ConsumerGroupHandler.
func (h *SaramaConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return h.Consumer.Cleanup(session)
}

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (h *SaramaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		err := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), h.ConsumerTimeout)
			defer cancel()
			return h.Consumer.Consume(ctx, session, message)
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

// Close implements sarama.ConsumerGroupHandler.
func (h *SaramaConsumerGroupHandler) Close() error { return nil }

// SaramaMessageConsumer processes Kafka messages.
type SaramaMessageConsumer interface {
	// Consume should process a message.
	//
	// Consume is responsible for marking the message consumed, unless the
	// context is canceled, in which case the caller should retry, or mark the
	// message as appropriate.
	Consume(ctx context.Context, session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) error

	// Setup rolls up to SaramaConsumerGroupHandler to implement [sarama.ConsumerGroupHandler].
	Setup(session sarama.ConsumerGroupSession) error

	// Cleanup rolls up to SaramaConsumerGroupHandler to implement [sarama.ConsumerGroupHandler].
	Cleanup(session sarama.ConsumerGroupSession) error
}

var ErrRetriesLimitExceeded = errors.New("retry limit exceeded")

// NTimesRetryingConsumer enhances a SaramaMessageConsumer with a finite
// number of immediate retries.
//
// The delay between each retry can be controlled via the Delay property. If
// no Delay property is specified, a delay based on the Fibonacci sequence is
// used.
//
// Logger is intentionally minimal. The slog.Log function is used by default.
type NTimesRetryingConsumer struct {
	Times    int
	Consumer SaramaMessageConsumer
	Delay    func(tries int) time.Duration
	Logger   func(ctx context.Context) Logger
}

// Logger is an intentionally minimal interface for basic logging.
//
// It matches the signature of slog.Log.
type Logger interface {
	Log(ctx context.Context, level slog.Level, msg string, args ...any)
}

func (c *NTimesRetryingConsumer) Consume(ctx context.Context,
	session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) error {

	var joinedErrors error
	var tries int = 0
	var delay time.Duration = 0
	if c.Delay == nil {
		c.Delay = DelayFibonacci
	}
	if c.Logger == nil {
		c.Logger = func(_ context.Context) Logger { return slog.Default() }
	}
	logger := c.Logger(ctx)
	done := ctx.Done()
	for tries < c.Times {
		select {
		case <-done:
			if ctxErr := ctx.Err(); ctxErr != nil {
				return ctxErr
			}
			return nil
		case <-time.After(delay):
			err := c.Consumer.Consume(ctx, session, message)
			if err == nil {
				return nil
			}
			if c.isContextErr(err) {
				return err
			} else if errors.Is(err, nil) {
				return nil
			}
			delay = c.Delay(tries)
			logger.Log(ctx, slog.LevelInfo, "failure consuming Kafka message, will retry",
				slog.Attr{Key: "tries", Value: slog.IntValue(tries)},
				slog.Attr{Key: "times", Value: slog.IntValue(c.Times)},
				slog.Attr{Key: "delay", Value: slog.DurationValue(delay)},
				slog.Attr{Key: "err", Value: slog.AnyValue(err)},
			)
			joinedErrors = errors.Join(joinedErrors, err)
			tries++
		}
	}

	return errors.Join(joinedErrors, c.retryLimitError())
}

func (c *NTimesRetryingConsumer) Setup(session sarama.ConsumerGroupSession) error {
	return c.Consumer.Setup(session)
}
func (c *NTimesRetryingConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return c.Consumer.Cleanup(session)
}

func (c *NTimesRetryingConsumer) isContextErr(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}

func (c *NTimesRetryingConsumer) retryLimitError() error {
	return fmt.Errorf("%w (%d)", ErrRetriesLimitExceeded, c.Times)
}

// DelayNone is a function returning a constant "no delay" of 0 seconds.
var DelayNone = func(_ int) time.Duration { return DelayConstant(0) }

// DelayConstant is a function returning a constant number of seconds.
func DelayConstant(n int) time.Duration { return time.Duration(n) * time.Second }

// DelayExponentialBinary returns a binary exponential delay.
//
// The delay is 2**tries seconds.
func DelayExponentialBinary(tries int) time.Duration {
	return time.Second * time.Duration(math.Pow(2, float64(tries)))
}

// DelayFibonacci returns a delay based on the Fibonacci sequence.
func DelayFibonacci(tries int) time.Duration {
	return time.Second * time.Duration(Fib(tries))
}

// Fib returns the nth number in the Fibonacci sequence.
func Fib(n int) int {
	if n == 0 {
		return 0
	} else if n < 3 {
		return 1
	}

	n1, n2 := 1, 1
	for i := 3; i <= n; i++ {
		n1, n2 = n1+n2, n1
	}

	return n1
}

// TopicShiftingRetryingConsumer retries by moving failing messages between topics.
//
// Inspired by https://www.uber.com/blog/reliable-reprocessing/
type TopicShiftingRetryingConsumer struct {
	// Delay before processing a message.
	//
	// The default, 0, will process immediately, and is suitable for use in the first
	// consumer of a sequence.
	Delay time.Duration
	// NextTopic is the target topic of failing messages.
	NextTopic string
	// Consumer processes messages, those that fail will be shifted to the next topic.
	Consumer SaramaMessageConsumer
	// Producer shifts the failed messages to the next topic.
	Producer sarama.AsyncProducer
	// GroupID for the sarama consumer group.
	GroupID string
}

func (c *TopicShiftingRetryingConsumer) Setup(session sarama.ConsumerGroupSession) error {
	// TODO start retry topic consumers? Or is that handled somewhere else?
	// if this consumer only had a single "next" topic... that might do the thing..?
	// TODO ensure retry consumers are in "read committed" mode
	return c.Consumer.Setup(session)
}

func (c *TopicShiftingRetryingConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	// TODO stop retry topic consumers
	return c.Consumer.Cleanup(session)
}

func (c *TopicShiftingRetryingConsumer) Consume(ctx context.Context,
	session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) error {

	select {
	case <-time.After(c.Delay):
		// no op
	case <-ctx.Done():
		return nil
	}

	err := c.Consumer.Consume(ctx, session, msg)
	if err != nil {
		if shiftErr := c.shiftMessage(msg); shiftErr != nil {
			return errors.Join(err, shiftErr)
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

func (c *TopicShiftingRetryingConsumer) shiftMessage(msg *sarama.ConsumerMessage) error {
	if txnErr := c.Producer.BeginTxn(); txnErr != nil {
		return fmt.Errorf("beginning transaction: %s", txnErr)
	}
	if txnErr := c.Producer.AddMessageToTxn(msg, c.GroupID, nil); txnErr != nil {
		return fmt.Errorf("adding message to transaction: %s", txnErr)
	}
	pHeaders := make([]sarama.RecordHeader, len(msg.Headers))
	for i := range msg.Headers {
		pHeaders[i] = *msg.Headers[i]
	}
	pMsg := &sarama.ProducerMessage{
		Topic:   c.NextTopic,
		Key:     sarama.ByteEncoder(msg.Key),
		Value:   sarama.ByteEncoder(msg.Value),
		Headers: pHeaders,
	}
	c.Producer.Input() <- pMsg
	if txnErr := c.Producer.CommitTxn(); txnErr != nil {
		return fmt.Errorf("committing transaction: %s", txnErr)
	}
	return nil
}

// SaramaClientConfig joins together the necessary configuration for a [sarama.Client].
type SaramaClientConfig struct {
	Brokers []string
	GroupID string
	Topics  []string

	SaramConfig *sarama.Config
}

type FanOutConsumer struct {
	Consumers []SaramaMessageConsumer
}

func NewFanOutConsumer2(consumers []SaramaMessageConsumer) *FanOutConsumer {
	return &FanOutConsumer{
		Consumers: consumers,
	}
}

func NewFanOutConsumer(topicBase string, producer sarama.AsyncProducer, groupID string,
	delays []time.Duration, consumer SaramaMessageConsumer) *FanOutConsumer {

	consumers := make([]SaramaMessageConsumer, 0, len(delays))
	for idx, delay := range delays {
		topicSuffix := "dead"
		if len(delays) > idx+1 {
			topicSuffix = delays[idx+1].String()
		}
		consumers = append(consumers, &TopicShiftingRetryingConsumer{
			NextTopic: topicBase + topicSuffix,
			Delay:     delay,
			Producer:  producer,
			Consumer:  consumer,
		})
	}

	return &FanOutConsumer{
		Consumers: consumers,
	}
}

func (c *FanOutConsumer) Consume(ctx context.Context, session sarama.ConsumerGroupSession,
	msg *sarama.ConsumerMessage) error {

	return c.Consumers[0].Consume(ctx, session, msg)
}

func (c *FanOutConsumer) Setup(session sarama.ConsumerGroupSession) error {
	var err error
	for _, consumer := range c.Consumers {
		err = errors.Join(err, consumer.Setup(session))
	}
	return err
}

func (c *FanOutConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	var err error
	for _, consumer := range c.Consumers {
		err = errors.Join(err, consumer.Cleanup(session))
	}
	return err
}
