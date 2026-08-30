package asyncevents

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"
)

// SaramaRunner interfaces between [events.Runner] and go-common's
// [SaramaEventsConsumer].
//
// This means providing Initialize(), Run(), and Terminate() to satisfy events.Runner, while
// under the hood calling SaramaEventConsumer's Run(), and canceling its Context as
// appropriate.
type SaramaRunner struct {
	eventsRunner SaramaEventsRunner
	cancelCtx    context.CancelFunc
	cancelMu     sync.Mutex
}

func NewSaramaRunner(eventsRunner SaramaEventsRunner) *SaramaRunner {
	return &SaramaRunner{
		eventsRunner: eventsRunner,
	}
}

// SaramaEventsRunner is implemented by go-common's [SaramaEventsRunner].
type SaramaEventsRunner interface {
	Run(ctx context.Context) error
}

// SaramaRunnerConfig collects values needed to initialize a SaramaRunner.
//
// This provides isolation for the SaramaRunner from ConfigReporter,
// envconfig, or any of the other options in platform for reading config
// values.
type SaramaRunnerConfig struct {
	Brokers         []string
	GroupID         string
	Topics          []string
	MessageConsumer SaramaMessageConsumer

	Sarama *sarama.Config
}

func (r *SaramaRunner) Initialize() error { return nil }

// Run adapts platform's event.Runner to work with go-common's
// SaramaEventsConsumer.
func (r *SaramaRunner) Run() error {
	if r.eventsRunner == nil {
		return errors.New("unable to run SaramaRunner, eventsRunner is nil")
	}

	r.cancelMu.Lock()
	ctx, err := func() (context.Context, error) {
		defer r.cancelMu.Unlock()
		if r.cancelCtx != nil {
			return nil, errors.New("unable to Run SaramaRunner, it's already initialized")
		}
		var ctx context.Context
		ctx, r.cancelCtx = context.WithCancel(context.Background())
		return ctx, nil
	}()
	if err != nil {
		return err
	}
	if err := r.eventsRunner.Run(ctx); err != nil {
		return fmt.Errorf("unable to Run SaramaRunner: %w", err)
	}
	return nil
}

// Terminate adapts platform's event.Runner to work with go-common's
// SaramaEventsConsumer.
func (r *SaramaRunner) Terminate() error {
	r.cancelMu.Lock()
	defer r.cancelMu.Unlock()
	if r.cancelCtx == nil {
		return errors.New("unable to Terminate SaramaRunner, it's not running")
	}
	r.cancelCtx()
	return nil
}

// CappedExponentialBinaryDelay builds delay functions that use exponential
// binary backoff with a maximum duration.
func CappedExponentialBinaryDelay(cap time.Duration) func(int) time.Duration {
	return func(tries int) time.Duration {
		b := DelayExponentialBinary(tries)
		if b > cap {
			return cap
		}
		return b
	}
}

// CascadingSaramaEventsRunner manages multiple sarama consumer groups to execute a
// topic-cascading retry process.
//
// The topic names are generated from Config.Topics combined with Delays. If given a single
// topic "updates", and delays: 0s, 1s, and 5s, then the following topics will be consumed:
// updates, updates-retry-1s, updates-retry-5s. The consumer of the updates-retry-5s topic
// will write failed messages to updates-dead.
//
// The inspiration for this system was drawn from
// https://www.uber.com/blog/reliable-reprocessing/
type CascadingSaramaEventsRunner struct {
	Config             SaramaRunnerConfig
	ConsumptionTimeout time.Duration
	Delays             []time.Duration
	Logger             Logger
	SaramaBuilders     SaramaBuilders
}

func NewCascadingSaramaEventsRunner(config SaramaRunnerConfig, logger Logger,
	delays []time.Duration, consumptionTimeout time.Duration) *CascadingSaramaEventsRunner {

	return &CascadingSaramaEventsRunner{
		Config:             config,
		Delays:             delays,
		Logger:             logger,
		SaramaBuilders:     DefaultSaramaBuilders{},
		ConsumptionTimeout: consumptionTimeout,
	}
}

// LimitedAsyncProducer restricts the [sarama.AsyncProducer] interface to ensure that its
// recipient isn't able to call Close(), thereby opening the potential for a panic when
// writing to a closed channel.
type LimitedAsyncProducer interface {
	AbortTxn() error
	BeginTxn() error
	CommitTxn() error
	Input() chan<- *sarama.ProducerMessage
}

func (r *CascadingSaramaEventsRunner) Run(ctx context.Context) error {
	if len(r.Config.Topics) == 0 {
		return errors.New("no topics")
	}
	if len(r.Delays) == 0 {
		return errors.New("no delays")
	}

	producersCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	errs := make(chan error, len(r.Config.Topics)*len(r.Delays))
	defer func() {
		r.Logger.Log(ctx, slog.LevelDebug, "CascadingSaramaEventsRunner: waiting for consumers")
		wg.Wait()
		r.Logger.Log(ctx, slog.LevelDebug, "CascadingSaramaEventsRunner: all consumers returned")
		close(errs)
	}()

	for _, topic := range r.Config.Topics {
		for idx, delay := range r.Delays {
			producerCfg := r.producerConfig(idx, delay)
			// The producer is built here rather than in buildConsumer() to control when
			// producer is closed. Were the producer to be closed before consumer.Run()
			// returns, it would be possible for consumer to write to the producer's
			// Inputs() channel, which if closed, would cause a panic.
			producer, err := r.SaramaBuilders.NewAsyncProducer(r.Config.Brokers, producerCfg)
			if err != nil {
				return fmt.Errorf("unable to build async producer %s: %w", r.Config.GroupID, err)
			}

			consumer, err := r.buildConsumer(producersCtx, idx, producer, delay, topic)
			if err != nil {
				return err
			}

			wg.Add(1)
			go func(topic string) {
				defer func() {
					closeErr := producer.Close()
					if closeErr != nil {
						r.Logger.Log(producersCtx, slog.LevelInfo, "CascadingSaramaEventsRunner: unable to close producer", "error", closeErr)
					}
					wg.Done()
				}()
				if err := consumer.Run(producersCtx); err != nil {
					errs <- fmt.Errorf("topics[%q]: %s", topic, err)
				}
				r.Logger.Log(ctx, slog.LevelDebug, "CascadingSaramaEventsRunner: consumer go proc returning", "topic", topic)
			}(topic)
		}
	}

	select {
	case <-ctx.Done():
		r.Logger.Log(ctx, slog.LevelDebug, "CascadingSaramaEventsRunner: context is done")
		return nil
	case err := <-errs:
		r.Logger.Log(ctx, slog.LevelDebug, "CascadingSaramaEventsRunner: Run(): error from consumer", "error", err)
		return err
	}
}

func (r *CascadingSaramaEventsRunner) producerConfig(idx int, delay time.Duration) *sarama.Config {
	uniqueConfig := *r.Config.Sarama
	hostID := os.Getenv("HOSTNAME") // set by default in kubernetes pods
	if hostID == "" {
		hostID = fmt.Sprintf("%d-%d", time.Now().UnixNano()/int64(time.Second), os.Getpid())
	}
	txnID := fmt.Sprintf("%s-%s-%d-%s", r.Config.GroupID, delay.String(), idx, hostID)
	uniqueConfig.Producer.Transaction.ID = txnID
	uniqueConfig.Producer.Idempotent = true
	uniqueConfig.Producer.RequiredAcks = sarama.WaitForAll
	uniqueConfig.Net.MaxOpenRequests = 1
	uniqueConfig.Consumer.IsolationLevel = sarama.ReadCommitted
	return &uniqueConfig
}

// SaramaBuilders allows tests to inject mock objects.
type SaramaBuilders interface {
	NewAsyncProducer([]string, *sarama.Config) (sarama.AsyncProducer, error)
	NewConsumerGroup([]string, string, *sarama.Config) (sarama.ConsumerGroup, error)
}

// DefaultSaramaBuilders implements SaramaBuilders for normal, non-test use.
type DefaultSaramaBuilders struct{}

func (DefaultSaramaBuilders) NewAsyncProducer(brokers []string, config *sarama.Config) (
	sarama.AsyncProducer, error) {

	return sarama.NewAsyncProducer(brokers, config)
}

func (DefaultSaramaBuilders) NewConsumerGroup(brokers []string, groupID string,
	config *sarama.Config) (sarama.ConsumerGroup, error) {

	return sarama.NewConsumerGroup(brokers, groupID, config)
}

func (r *CascadingSaramaEventsRunner) buildConsumer(ctx context.Context, idx int,
	producer LimitedAsyncProducer, delay time.Duration, baseTopic string) (
	*SaramaEventsConsumer, error) {

	groupID := r.Config.GroupID
	if delay > 0 {
		groupID += "-retry-" + delay.String()
	}
	group, err := r.SaramaBuilders.NewConsumerGroup(r.Config.Brokers, groupID,
		r.Config.Sarama)
	if err != nil {
		return nil, fmt.Errorf("unable to build sarama consumer group %s: %w", groupID, err)
	}

	var consumer = r.Config.MessageConsumer
	if len(r.Delays) > 0 {
		nextTopic := baseTopic + "-dead"
		if idx+1 < len(r.Delays) {
			nextTopic = baseTopic + "-retry-" + r.Delays[idx+1].String()
		}
		consumer = &CascadingConsumer{
			Consumer:  consumer,
			NextTopic: nextTopic,
			Producer:  producer,
			Logger:    r.Logger,
		}
	}
	if delay > 0 {
		consumer = &NotBeforeConsumer{
			Consumer: consumer,
			Logger:   r.Logger,
		}
	}

	handler := NewSaramaConsumerGroupHandler(r.Logger, consumer, r.ConsumptionTimeout)
	topic := baseTopic
	if delay > 0 {
		topic += "-retry-" + delay.String()
	}
	r.Logger.Log(ctx, slog.LevelDebug, "creating consumer", "topic", topic)

	return NewSaramaEventsConsumer(group, handler, topic), nil
}

// NotBeforeConsumer delays consumption until a specified time.
type NotBeforeConsumer struct {
	Consumer SaramaMessageConsumer
	Logger   Logger
}

func (c *NotBeforeConsumer) Consume(ctx context.Context, session sarama.ConsumerGroupSession,
	msg *sarama.ConsumerMessage) error {

	notBefore, err := c.notBeforeFromMsgHeaders(msg)
	if err != nil {
		c.Logger.Log(ctx, slog.LevelInfo, "Unable to parse kafka header not-before value", "error", err)
	}
	delay := time.Until(notBefore)

	select {
	case <-ctx.Done():
		if ctxErr := ctx.Err(); !errors.Is(ctxErr, context.Canceled) {
			return ctxErr
		}
		return nil
	case <-time.After(time.Until(notBefore)):
		if !notBefore.IsZero() {
			c.Logger.Log(ctx, slog.LevelDebug, "delayed", "topic", msg.Topic, "not-before", notBefore, "delay", delay)
		}
		return c.Consumer.Consume(ctx, session, msg)
	}
}

// HeaderNotBefore tells consumers not to consume a message before a certain time.
var HeaderNotBefore = []byte("x-tidepool-not-before")

// NotBeforeTimeFormat specifies the [time.Parse] format to use for HeaderNotBefore.
var NotBeforeTimeFormat = time.RFC3339Nano

// HeaderFailures counts the number of failures encountered trying to consume the message.
var HeaderFailures = []byte("x-tidepool-failures")

// FailuresToDelay maps the number of consumption failures to the next delay.
//
// Rather than using a failures header, the name of the topic could be used as a lookup, if
// so desired.
var FailuresToDelay = map[int]time.Duration{
	0: 0,
	1: 1 * time.Second,
	2: 2 * time.Second,
	3: 3 * time.Second,
	4: 5 * time.Second,
}

func (c *NotBeforeConsumer) notBeforeFromMsgHeaders(msg *sarama.ConsumerMessage) (
	time.Time, error) {

	for _, header := range msg.Headers {
		if bytes.Equal(header.Key, HeaderNotBefore) {
			notBefore, err := time.Parse(NotBeforeTimeFormat, string(header.Value))
			if err != nil {
				return time.Time{}, fmt.Errorf("parsing not before header: %s", err)
			} else {
				return notBefore, nil
			}
		}
	}
	return time.Time{}, fmt.Errorf("header not found: x-tidepool-not-before")
}

// CascadingConsumer cascades messages that failed to be consumed to another topic.
//
// It also sets an adjustable delay via the "not-before" and "failures" headers so that as
// the message moves from topic to topic, the time between processing is increased according
// to [FailuresToDelay].
type CascadingConsumer struct {
	Consumer  SaramaMessageConsumer
	NextTopic string
	Producer  LimitedAsyncProducer
	Logger    Logger
}

func (c *CascadingConsumer) Consume(ctx context.Context, session sarama.ConsumerGroupSession,
	msg *sarama.ConsumerMessage) (err error) {

	if err := c.Consumer.Consume(ctx, session, msg); err != nil {
		txnErr := c.withTxn(func() error {
			select {
			case <-ctx.Done():
				if ctxErr := ctx.Err(); !errors.Is(ctxErr, context.Canceled) {
					return ctxErr
				}
				return nil
			case c.Producer.Input() <- c.cascadeMessage(msg):
				c.Logger.Log(ctx, slog.LevelInfo, "cascaded", "from", msg.Topic, "to", c.NextTopic)
				return nil
			}
		})
		if txnErr != nil {
			c.Logger.Log(ctx, slog.LevelInfo, "Unable to complete cascading transaction", "error", err)
			return err
		}
	}
	return nil
}

// withTxn wraps a function with a transaction that is aborted if an error is returned.
func (c *CascadingConsumer) withTxn(f func() error) (err error) {
	if err := c.Producer.BeginTxn(); err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer func(err *error) {
		if err != nil && *err != nil {
			if abortErr := c.Producer.AbortTxn(); abortErr != nil {
				c.Logger.Log(nil, slog.LevelInfo, "Unable to abort transaction", "error", abortErr)
			}
			return
		}
		if commitErr := c.Producer.CommitTxn(); commitErr != nil {
			c.Logger.Log(nil, slog.LevelInfo, "Unable to commit transaction", "error", commitErr)
		}
	}(&err)
	return f()
}

// cascadeMessage to the next topic.
func (c *CascadingConsumer) cascadeMessage(msg *sarama.ConsumerMessage) *sarama.ProducerMessage {
	pHeaders := make([]sarama.RecordHeader, len(msg.Headers))
	for idx, header := range msg.Headers {
		pHeaders[idx] = *header
	}
	return &sarama.ProducerMessage{
		Key:     sarama.ByteEncoder(msg.Key),
		Value:   sarama.ByteEncoder(msg.Value),
		Topic:   c.NextTopic,
		Headers: c.updateCascadeHeaders(pHeaders),
	}
}

// updateCascadeHeaders calculates not before and failures header values.
//
// Existing not before and failures headers will be dropped in place of the new ones.
func (c *CascadingConsumer) updateCascadeHeaders(headers []sarama.RecordHeader) []sarama.RecordHeader {
	failures := 0
	notBefore := time.Now()

	keep := make([]sarama.RecordHeader, 0, len(headers))
	for _, header := range headers {
		switch {
		case bytes.Equal(header.Key, HeaderNotBefore):
			continue // Drop this header, we'll add a new version below.
		case bytes.Equal(header.Key, HeaderFailures):
			parsed, err := strconv.ParseInt(string(header.Value), 10, 32)
			if err != nil {
				c.Logger.Log(nil, slog.LevelInfo, "Unable to parse consumption failures count", "error", err)
			} else {
				failures = int(parsed)
				notBefore = notBefore.Add(FailuresToDelay[failures])
			}
			continue // Drop this header, we'll add a new version below.
		}
		keep = append(keep, header)
	}

	keep = append(keep, sarama.RecordHeader{
		Key:   HeaderNotBefore,
		Value: []byte(notBefore.Format(NotBeforeTimeFormat)),
	})
	keep = append(keep, sarama.RecordHeader{
		Key:   HeaderFailures,
		Value: []byte(strconv.Itoa(failures + 1)),
	})

	return keep
}
