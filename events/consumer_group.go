package events

import (
	"context"
	"log"
	"sync"

	"github.com/Shopify/sarama"

	"github.com/tidepool-org/go-common/errors"
)

var ErrConsumerStopped = errors.New("consumer has been stopped")

type SaramaConsumerGroup struct {
	config        *CloudEventsConfig
	consumerGroup sarama.ConsumerGroup
	consumer      MessageConsumer
	topic         string

	cancelFuncMu sync.Mutex
	cancelFunc   context.CancelFunc
}

func NewSaramaConsumerGroup(config *CloudEventsConfig, consumer MessageConsumer) (EventConsumer, error) {
	if err := validateConsumerConfig(config); err != nil {
		return nil, err
	}

	return &SaramaConsumerGroup{
		config:   config,
		consumer: consumer,
		topic:    config.GetPrefixedTopic(),
	}, nil
}

func validateConsumerConfig(config *CloudEventsConfig) error {
	if config.KafkaConsumerGroup == "" {
		return errors.New("consumer group cannot be empty")
	}
	return nil
}

func (s *SaramaConsumerGroup) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s *SaramaConsumerGroup) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s *SaramaConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := s.consumer.HandleKafkaMessage(message); err != nil {
			log.Printf("failed to process kafka message: %v", err)
			return err
		}
		session.MarkMessage(message, "")
	}

	return nil
}

func (s *SaramaConsumerGroup) Start() error {
	if err := s.initialize(); err != nil {
		return err
	}

	ctx, cancel := s.newContext()
	defer cancel()

	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := s.consumerGroup.Consume(ctx, []string{s.topic}, s); err != nil {
			log.Printf("Error from consumer: %v", err)
			if err == context.Canceled {
				return ErrConsumerStopped
			}
			return err
		}
	}
}

func (s *SaramaConsumerGroup) Stop() error {
	s.cancelFuncMu.Lock()
	defer s.cancelFuncMu.Unlock()

	if s.cancelFunc != nil {
		s.cancelFunc()
		s.cancelFunc = nil
	}
	return nil
}

func (s *SaramaConsumerGroup) newContext() (context.Context, context.CancelFunc) {
	s.cancelFuncMu.Lock()
	defer s.cancelFuncMu.Unlock()
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel
	return ctx, cancel
}

func (s *SaramaConsumerGroup) initialize() error {
	cg, err := sarama.NewConsumerGroup(
		s.config.KafkaBrokers,
		s.config.KafkaConsumerGroup,
		s.config.SaramaConfig,
	)
	if err != nil {
		return err
	}

	if err := s.consumer.Initialize(s.config); err != nil {
		return err
	}

	s.consumerGroup = cg
	return nil
}
