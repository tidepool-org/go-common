package events

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.uber.org/fx"
)

//ExponentialRetry implements exponential backoff
type ExponentialRetry struct {
	attempt      int
	maxAttempts  int
	delay        time.Duration
	maxDelay     time.Duration
	initialDelay time.Duration
	multiplier   float32
	adder        time.Duration
}

//RetryAlgorithm implements a retry mechanism
type RetryAlgorithm interface {
	Retry() bool // returns true after possible delay if should attempt retry
}

//NewExponentialRetry creates a new exponential backoff mechanism
// retries will not be attempted if the delay exceeds maxDelay, so if initialDelay exceeds maxDelay, then no retries will be attempted
// further retries will not be attempted if maxAttempts have already been tried, so if you set maxAttempts to 0, then no retries will be attempted
func NewExponentialRetry(config *CloudEventsConfig) *ExponentialRetry {
	return &ExponentialRetry{
		attempt:      0,
		delay:        0,
		maxDelay:     config.RetryMaxDelay,
		initialDelay: config.RetryInitialDelay,
		multiplier:   config.RetryFactor,
		adder:        config.RetryAddend,
		maxAttempts:  config.RetryMaxAttempts,
	}
}

// Retry implements retry algorithm
func (r *ExponentialRetry) Retry(ctx context.Context) bool {
	r.attempt = r.attempt + 1
	if r.attempt > r.maxAttempts {
		return false
	}
	if r.attempt == 1 {
		r.delay = r.initialDelay
	} else {
		r.delay = time.Duration(int64(r.multiplier * float32(r.delay)))
		r.delay = r.delay + r.adder
	}
	if r.delay > r.maxDelay {
		return false
	}

	select {
	case <-time.After(r.delay):
		return true
	case <-ctx.Done():
		return false
	}
}

type SaramaConsumer struct {
	config             *CloudEventsConfig
	consumerGroup      sarama.ConsumerGroup
	ready              chan bool
	topic              string
	handlers           []EventHandler
	deadLetterProducer *KafkaCloudEventsProducer
	retry              *ExponentialRetry
}

//CascadingEventConsumer an event consumer that cascaded failures
type CascadingEventConsumer struct {
	Consumers []EventConsumer
}

// NewCascadingCloudEventsConsumer create a cascading events consumer
// For every delay in config.CascadeDelays, a new consumer in addition to the one for the primary topic
func NewCascadingCloudEventsConsumer(config *CloudEventsConfig) (EventConsumer, error) {
	topic := config.KafkaTopic
	delay := config.KafkaDelay
	var consumers []EventConsumer
	for i, nextDelay := range config.CascadeDelays {
		newconfig := *config
		newconfig.KafkaDelay = delay
		newconfig.KafkaTopic = topic
		newconfig.KafkaDeadLettersTopic = fmt.Sprintf(config.CascadePattern, config.KafkaTopic, i, nextDelay)
		consumer, err := NewSaramaCloudEventsConsumer(&newconfig)
		if err != nil {
			return nil, err
		}
		consumers = append(consumers, consumer)
		topic = newconfig.KafkaDeadLettersTopic
		delay = nextDelay
	}
	newconfig := *config
	newconfig.KafkaDelay = delay
	newconfig.KafkaTopic = topic
	consumer, err := NewSaramaCloudEventsConsumer(&newconfig)
	if err != nil {
		return nil, err
	}
	consumers = append(consumers, consumer)

	return &CascadingEventConsumer{
		Consumers: consumers,
	}, nil
}

//RegisterHandler registers the handler with all consumers in the cascade
func (c *CascadingEventConsumer) RegisterHandler(handler EventHandler) {
	for _, consumer := range c.Consumers {
		consumer.RegisterHandler(handler)
	}
}

//Start starts an async consumer goproc (that may sleep)
func (c *CascadingEventConsumer) Start(ctx context.Context) error {
	for _, consumer := range c.Consumers {
		err := consumer.Start(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// eventConsumerWrapper provides support for Open Telemetry tracing of Kafka
type eventConsumerWrapper struct {
	consumer *SaramaConsumer
	handler  sarama.ConsumerGroupHandler
}

func newEventConsumerWrapper(consumer *SaramaConsumer) *eventConsumerWrapper {
	return &eventConsumerWrapper{
		consumer,
		otelsarama.WrapConsumerGroupHandler(consumer),
	}
}

func (e *eventConsumerWrapper) Setup(s sarama.ConsumerGroupSession) error {
	return e.handler.Setup(s)
}

func (e *eventConsumerWrapper) Cleanup(s sarama.ConsumerGroupSession) error {
	return e.handler.Cleanup(s)
}

func (e *eventConsumerWrapper) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	return e.handler.ConsumeClaim(s, c)
}

func (e *eventConsumerWrapper) RegisterHandler(handler EventHandler) {
	e.consumer.RegisterHandler(handler)
}

func (e *eventConsumerWrapper) Start(ctx context.Context) error {
	return e.consumer.Start(ctx)
}

//CloudEventsConfigProvider provides a cloud events config
func CloudEventsConfigProvider() (*CloudEventsConfig, error) {
	cfg := NewConfig()
	if err := cfg.LoadFromEnv(); err != nil {
		return nil, err
	}
	return cfg, nil
}

//CloudEventsConsumerProvider provides a cloud events consumer
func CloudEventsConsumerProvider(config *CloudEventsConfig, handler EventHandler) (EventConsumer, error) {
	consumer, err := NewSaramaCloudEventsConsumer(config)
	if err != nil {
		return nil, err
	}
	consumer.RegisterHandler(handler)
	return consumer, nil
}

//NewSaramaCloudEventsConsumer creates a new cloud events consumer
func NewSaramaCloudEventsConsumer(config *CloudEventsConfig) (EventConsumer, error) {
	if err := validateConsumerConfig(config); err != nil {
		return nil, err
	}

	consumer := &SaramaConsumer{
		config:   config,
		ready:    make(chan bool),
		topic:    config.GetPrefixedTopic(),
		handlers: make([]EventHandler, 0),
		retry:    NewExponentialRetry(config),
	}

	handler := newEventConsumerWrapper(consumer)
	return handler, nil
}

//StartEventConsumer starts an event consumer
func StartEventConsumer(consumer EventConsumer, lifecycle fx.Lifecycle) {
	consumerCtx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{}, 1)
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				// blocks until context is terminated
				err := consumer.Start(consumerCtx)
				if err != nil {
					log.Printf("Unable to start cloud events consumer: %v", err)
				}
				done <- struct{}{}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			cancel()
			<-done
			return nil
		},
	})
}

//Setup marks a consumer to be ready
func (s *SaramaConsumer) Setup(session sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(s.ready)
	return nil
}

//Cleanup frees any resources
func (s *SaramaConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

//ConsumeClaim synchronously consumes all the messages in a claim, optionally delaying consumption a fixed amount of time
func (s *SaramaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	delay := time.Second * time.Duration(s.config.KafkaDelay)
	for message := range claim.Messages() {
		m := kafka_sarama.NewMessageFromConsumerMessage(message)
		timestamp := message.Timestamp
		scheduledTime := timestamp.Add(delay)
		remaining := time.Now().Sub(scheduledTime)
		if delay == 0 || remaining > 0 {
			timer := time.NewTimer(remaining)
			select {
			case <-timer.C:
				break
			case <-session.Context().Done():
				timer.Stop()
				return nil
			}
		}
		// just ignore non-cloud event messages
		if rs, rserr := binding.ToEvent(context.Background(), m); rserr == nil {
			err := s.handleCloudEvent(session.Context(), *rs)
			if err != nil {
				log.Printf("failed to process event %s", err)
				return err
			}
		}
		session.MarkMessage(message, "")
	}

	return nil
}

func (s *SaramaConsumer) handleCloudEvent(ctx context.Context, ce cloudevents.Event) error {
	var errors []error
	for _, handler := range s.handlers {
		if handler.CanHandle(ce) {
			if err := handler.Handle(ce); err != nil {
				errors = append(errors, err)
			}
		}
	}
	if len(errors) != 0 {
		if s.retry.Retry(ctx) {
			return s.handleCloudEvent(ctx, ce)
		}
		if s.config.IsDeadLettersEnabled() {
			return s.sendToDeadLetterTopic(ce)
		}
		return fmt.Errorf("no dead letter topic, skipping event %v, %v", ce.ID(), errors)
	}
	return nil
}

func (s *SaramaConsumer) sendToDeadLetterTopic(ce cloudevents.Event) error {
	if err := s.deadLetterProducer.SendCloudEvent(context.Background(), ce); err != nil {
		log.Printf("Failed to send event %v to dead-letter topic: %v", ce, err)
		return err
	}
	return nil
}

//RegisterHandler register a handler to process events
func (s *SaramaConsumer) RegisterHandler(handler EventHandler) {
	s.handlers = append(s.handlers, handler)
}

//Start starts an async consumer goproc (that may sleep)
func (s *SaramaConsumer) Start(ctx context.Context) error {
	if err := s.initialize(); err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := s.consumerGroup.Consume(ctx, []string{s.topic}, s); err != nil {
				log.Printf("Error from consumer: %v", err)
				return
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			s.ready = make(chan bool)
		}
	}()

	wg.Wait()
	return s.consumerGroup.Close()
}

func (s *SaramaConsumer) initialize() error {
	cg, err := sarama.NewConsumerGroup(
		s.config.KafkaBrokers,
		s.config.KafkaConsumerGroup,
		s.config.SaramaConfig,
	)
	if err != nil {
		return err
	}

	if s.config.IsDeadLettersEnabled() {
		s.deadLetterProducer, err = NewKafkaCloudEventsProducerForDeadLetters(s.config)
		if err != nil {
			return err
		}
	}

	s.consumerGroup = cg
	return nil
}
