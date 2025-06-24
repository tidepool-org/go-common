package asyncevents

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

// CascadingSaramaMessageConsumer cascades messages that failed to be consumed to another
// topic. It is an implementation of [SaramaMessageConsumer].
//
// It also sets an adjustable delay via the "not-before" and "failures" headers so that as
// the message moves from topic to topic, the time between processing is increased according
// to [FailuresToDelay].
type CascadingSaramaMessageConsumer struct {
	Consumer  SaramaMessageConsumer
	NextTopic string
	Producer  LimitedAsyncProducer
	Logger    Logger
}

// Consume implements [SaramaMessageConsumer].
func (c *CascadingSaramaMessageConsumer) Consume(ctx context.Context,
	session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) (err error) {

	if err := c.Consumer.Consume(ctx, session, msg); err != nil {
		txnErr := c.withTxn(ctx, func() error {
			select {
			case <-ctx.Done():
				if ctxErr := ctx.Err(); !errors.Is(ctxErr, context.Canceled) {
					return ctxErr
				}
				return nil
			case c.Producer.Input() <- c.cascadeMessage(ctx, msg):
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
func (c *CascadingSaramaMessageConsumer) withTxn(ctx context.Context, f func() error) (err error) {
	if err := c.Producer.BeginTxn(); err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer func(err *error) {
		if err != nil && *err != nil {
			if abortErr := c.Producer.AbortTxn(); abortErr != nil {
				c.Logger.Log(ctx, slog.LevelInfo, "Unable to abort transaction", "error", abortErr)
			}
			return
		}
		if commitErr := c.Producer.CommitTxn(); commitErr != nil {
			c.Logger.Log(ctx, slog.LevelInfo, "Unable to commit transaction", "error", commitErr)
		}
	}(&err)
	return f()
}

// cascadeMessage to the next topic.
func (c *CascadingSaramaMessageConsumer) cascadeMessage(ctx context.Context,
	msg *sarama.ConsumerMessage) *sarama.ProducerMessage {

	pHeaders := make([]sarama.RecordHeader, len(msg.Headers))
	for idx, header := range msg.Headers {
		pHeaders[idx] = *header
	}
	return &sarama.ProducerMessage{
		Key:     sarama.ByteEncoder(msg.Key),
		Value:   sarama.ByteEncoder(msg.Value),
		Topic:   c.NextTopic,
		Headers: c.updateCascadeHeaders(ctx, pHeaders),
	}
}

// updateCascadeHeaders calculates not before and failures header values.
//
// Existing not before and failures headers will be dropped in place of the new ones.
func (c *CascadingSaramaMessageConsumer) updateCascadeHeaders(ctx context.Context,
	headers []sarama.RecordHeader) []sarama.RecordHeader {

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
				c.Logger.Log(ctx, slog.LevelInfo, "Unable to parse consumption failures count", "error", err)
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

// CascadingSaramaEventsManagerConfig for a [CascadingSaramaEventsManager].
type CascadingSaramaEventsManagerConfig struct {
	Consumer SaramaMessageConsumer

	Brokers            []string
	GroupID            string
	Topics             []string
	ConsumptionTimeout time.Duration
	Delays             []time.Duration
	Logger             Logger
	SaramaBuilders     SaramaBuilders
	Sarama             *sarama.Config
}

// CascadingSaramaEventsManager manages multiple Sarama consumer groups to execute a
// topic-cascading retry process. It coordinates multiple [SaramaConsumerGroupManager]
// instances to achieve this.
//
// The topics' names are generated from a combination of the configured topics and
// configured delays. For example, if configured with a topic "updates", and delays: 0s, 1s,
// and 5s, then the following topics will be consumed: updates, updates-retry-1s,
// updates-retry-5s. The consumer of the updates-retry-5s topic will write failed messages
// to updates-dead.
//
// The inspiration for this system was drawn from
// https://www.uber.com/blog/reliable-reprocessing/
type CascadingSaramaEventsManager struct {
	CascadingSaramaEventsManagerConfig
}

func NewCascadingSaramaEventsManager(config CascadingSaramaEventsManagerConfig) *CascadingSaramaEventsManager {
	if config.SaramaBuilders == nil {
		config.SaramaBuilders = &DefaultSaramaBuilders{}
	}
	return &CascadingSaramaEventsManager{
		CascadingSaramaEventsManagerConfig: config,
	}
}

func (c *CascadingSaramaEventsManager) Run(ctx context.Context) error {
	if len(c.Topics) == 0 {
		return errors.New("no topics")
	}
	if len(c.Delays) == 0 {
		return errors.New("no delays")
	}

	producersCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	errs := make(chan error, len(c.Topics)*len(c.Delays))
	defer func() {
		c.Logger.Log(ctx, slog.LevelDebug, "CascadingSaramaEventsManager: waiting for managers")
		wg.Wait()
		c.Logger.Log(ctx, slog.LevelDebug, "CascadingSaramaEventsManager: all managers returned")
		close(errs)
	}()

	for _, topic := range c.Topics {
		for idx, delay := range c.Delays {
			producerCfg := c.producerConfig(idx, delay)
			// The producer is built here rather than in buildManager() to control when
			// producer is closed. Were the producer to be closed before manager.Run()
			// returns, it would be possible for manager to write to the producer's
			// Inputs() channel, which if closed, would cause a panic.
			producer, err := c.SaramaBuilders.NewAsyncProducer(c.Brokers, producerCfg)
			if err != nil {
				return fmt.Errorf("unable to build async producer %s: %w", c.GroupID, err)
			}

			manager, err := c.buildManager(producersCtx, idx, producer, delay, topic)
			if err != nil {
				return err
			}

			wg.Add(1)
			go func(topic string) {
				defer func() {
					closeErr := producer.Close()
					if closeErr != nil {
						c.Logger.Log(producersCtx, slog.LevelInfo, "CascadingSaramaEventsManager: unable to close producer", "error", closeErr)
					}
					wg.Done()
				}()
				if err := manager.Run(producersCtx); err != nil {
					errs <- fmt.Errorf("topics[%q]: %s", topic, err)
				}
				c.Logger.Log(ctx, slog.LevelDebug, "CascadingSaramaEventsManager: manager go proc returning", "topic", topic)
			}(topic)
		}
	}

	select {
	case <-ctx.Done():
		c.Logger.Log(ctx, slog.LevelDebug, "CascadingSaramaEventsManager: context is done")
		return nil
	case err := <-errs:
		c.Logger.Log(ctx, slog.LevelDebug, "CascadingSaramaEventsManager: Run(): error from manager", "error", err)
		return err
	}
}

func (c *CascadingSaramaEventsManager) producerConfig(idx int, delay time.Duration) *sarama.Config {
	uniqueConfig := *c.Sarama
	hostID := os.Getenv("HOSTNAME") // set by default in kubernetes pods
	if hostID == "" {
		hostID = fmt.Sprintf("%d-%d", time.Now().UnixNano()/int64(time.Second), os.Getpid())
	}
	txnID := fmt.Sprintf("%s-%s-%d-%s", c.GroupID, delay.String(), idx, hostID)
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

// buildManager returns a [SaramaConsumerGroupManager] that manages multiple, composed,
// [SaramaMessageConsumer]s according to its topics and delays configuration.
func (c *CascadingSaramaEventsManager) buildManager(ctx context.Context, idx int,
	producer LimitedAsyncProducer, delay time.Duration, baseTopic string) (
	*SaramaConsumerGroupManager, error) {

	groupID := c.GroupID
	if delay > 0 {
		groupID += "-retry-" + delay.String()
	}
	group, err := c.SaramaBuilders.NewConsumerGroup(c.Brokers, groupID,
		c.Sarama)
	if err != nil {
		return nil, fmt.Errorf("unable to build sarama consumer group %s: %w", groupID, err)
	}

	var consumer SaramaMessageConsumer = c.Consumer
	if len(c.Delays) > 0 {
		nextTopic := baseTopic + "-dead"
		if idx+1 < len(c.Delays) {
			nextTopic = baseTopic + "-retry-" + c.Delays[idx+1].String()
		}
		consumer = &CascadingSaramaMessageConsumer{
			Consumer:  consumer,
			NextTopic: nextTopic,
			Producer:  producer,
			Logger:    c.Logger,
		}
	}
	if delay > 0 {
		consumer = &NotBeforeConsumer{
			Consumer: consumer,
			Logger:   c.Logger,
		}
	}

	handler := NewSaramaConsumerGroupHandler(c.Logger, consumer, c.ConsumptionTimeout)
	topic := baseTopic
	if delay > 0 {
		topic += "-retry-" + delay.String()
	}
	c.Logger.Log(ctx, slog.LevelDebug, "creating consumer", "topic", topic)

	return NewSaramaConsumerGroupManager(group, handler, topic), nil
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

// LimitedAsyncProducer restricts the [sarama.AsyncProducer] interface to ensure that its
// recipient isn't able to call Close(), thereby opening the potential for a panic when
// writing to a closed channel.
type LimitedAsyncProducer interface {
	AbortTxn() error
	BeginTxn() error
	CommitTxn() error
	Input() chan<- *sarama.ProducerMessage
}
