package events

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/kelseyhightower/envconfig"
)

type CloudEventsConfig struct {
	KafkaBroker        string `envconfig:"KAFKA_BROKERS" required:"true"`
	KafkaPrefix        string `envconfig:"KAFKA_PREFIX" required:"true"`
	KafkaBaseTopic     string `envconfig:"KAFKA_TOPIC" default:"events"`
	KafkaConsumerGroup string `envconfig:"KAFKA_CONSUMER_GROUP" required:"false"`
	EventSource        string `envconfig:"CLOUD_EVENTS_SOURCE" required:"false"`
	SaramaConfig       *sarama.Config
}

func (k *CloudEventsConfig) LoadFromEnv() error {
	if err := envconfig.Process("", k); err != nil {
		return err
	}
	return nil
}

func (k *CloudEventsConfig) GetTopic() string {
	return k.KafkaPrefix + k.KafkaBaseTopic
}

func validateProducerConfig(config *CloudEventsConfig) error {
	if config.KafkaConsumerGroup == "" {
		return errors.New("consumer group cannot be empty")
	}
	if config.EventSource == "" {
		return errors.New("event source cannot be empty")
	}
	return nil
}

