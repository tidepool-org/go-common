package asyncevents

import (
	"context"
	stderrors "errors"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

// SaramaReceiver implements asynchronous event processing for messages
// received via Kafka.
type SaramaProcessor struct {
	cancel     context.CancelFunc
	config     SaramaConfig
	processors []SaramaConsumerMessageProcessor
}

type SaramaConsumerMessageProcessor interface {
	Process(ctx context.Context, msg *sarama.ConsumerMessage) error
	Topics() []string
}

type SubProcessor struct{}

func (p *SubProcessor) Process(ctx context.Context, msg *sarama.ConsumerMessage) error { return nil }

func (p *SubProcessor) Topics() []string { return nil }

func NewSaramaProcessor(config SaramaConfig, processors []SaramaConsumerMessageProcessor) *SaramaProcessor {
	return &SaramaProcessor{
		config:     config,
		processors: processors,
	}
}

// Start processing messages.
//
// The caller must call Stop() lest memory be leaked.
func (p *SaramaProcessor) Start(ctx context.Context) error {
	consumerGroup, err := p.config.NewConsumerGroup()
	if err != nil {
		return err
	}

	ctx, cancel := NewContextWithSafeCancel(ctx)
	p.cancel = cancel // Stop() calls p.cancel()

	handler := &SaramaConsumerGroupHandler{p.processors}
	topics := p.uniqueTopics()
	for {
		err := consumerGroup.Consume(ctx, topics, handler)
		if err != nil {
			if stderrors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
	}
}

// uniqueTopics builds a list of unique topics by querying each processors'
// topics.
func (p *SaramaProcessor) uniqueTopics() []string {
	topicsMap := map[string]struct{}{}
	topics := []string{}
	for _, processor := range p.processors {
		for _, topic := range processor.Topics() {
			if _, seen := topicsMap[topic]; !seen {
				topicsMap[topic] = struct{}{}
				topics = append(topics, topic)
			}
		}
	}
	return topics
}

// Stop processing.
func (p *SaramaProcessor) Stop() error {
	p.cancel()
	return nil
}

// SaramaConfig wraps up the various config needed to create a
// sarama.ConsumerGroup.
type SaramaConfig struct {
	Addrs   []string
	GroupID string
	Config  *sarama.Config
}

// NewConsumerGroup is a helper for creating a sarama.ConsumerGroup.
func (c SaramaConfig) NewConsumerGroup() (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroup(c.Addrs, c.GroupID, c.Config)
}

// SaramaConsumerGroupHandler implements sarama.ConsumerGroupHandler.
type SaramaConsumerGroupHandler struct {
	processors []SaramaConsumerMessageProcessor
}

// Setup implements sarama.ConsumerGroupHandler.
func (h *SaramaConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup implements sarama.ConsumerGroupHandler.
func (h *SaramaConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

var (
	MaxMessageProcessingDuration = time.Minute
)

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (h *SaramaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		err := h.handleMessage(message)
		if err != nil {
			return err
		}
		// success
		session.MarkMessage(message, "")
	}

	return nil
}

func (h *SaramaConsumerGroupHandler) handleMessage(message *sarama.ConsumerMessage) error {
	msgErrs := []error{}
	for _, proc := range h.processors {
		err := h.processWithTimeout(proc, message, MaxMessageProcessingDuration)
		if err != nil {
			msgErrs = append(msgErrs, err)
		}
	}
	if len(msgErrs) > 0 {
		// TODO dispatch some error handling, and probably pass in msgErrs... Or a map from processor to error?
		return stderrors.Join(msgErrs...)
	}

	return nil
}

// processWithTimeout is a helper that wraps a call to Process with
// timeout-enabled context.
func (h *SaramaConsumerGroupHandler) processWithTimeout(proc SaramaConsumerMessageProcessor,
	message *sarama.ConsumerMessage, timeout time.Duration) error {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return proc.Process(ctx, message)
}

// NewContextWithSafeCancel wraps a context.Context with a thread-safe cancel
// function.
func NewContextWithSafeCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	cancelable, cancel := context.WithCancel(ctx)
	var once sync.Once
	return cancelable, func() { once.Do(cancel) }
}
