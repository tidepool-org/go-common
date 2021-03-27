package events

import (
	"errors"
	"time"

	"github.com/Shopify/sarama"
	"github.com/kelseyhightower/envconfig"
)

const DeadLetterSuffix = "-dead-letters"

// CloudEventsConfig describes the configuration for a consumer of cloud events from a Kafka topic
type CloudEventsConfig struct {
	EventSource           string          `envconfig:"CLOUD_EVENTS_SOURCE" required:"true"`
	KafkaBrokers          []string        `envconfig:"KAFKA_BROKERS" required:"true"`
	KafkaConsumerGroup    string          `envconfig:"KAFKA_CONSUMER_GROUP" required:"false"`
	KafkaTopic            string          `envconfig:"KAFKA_TOPIC" default:"events"`
	KafkaDeadLettersTopic string          `envconfig:"KAFKA_DEAD_LETTERS_TOPIC"`
	KafkaTopicPrefix      string          `envconfig:"KAFKA_TOPIC_PREFIX" required:"true"`
	KafkaRequireSSL       bool            `envconfig:"KAFKA_REQUIRE_SSL" required:"true"`
	KafkaVersion          string          `envconfig:"KAFKA_VERSION" required:"true"`
	KafkaUsername         string          `envconfig:"KAFKA_USERNAME" required:"false"`
	KafkaPassword         string          `envconfig:"KAFKA_PASSWORD" required:"false"`
	KafkaDelay            time.Duration   `envconfig:"KAFKA_DELAY" default:"0s"`
	CascadeDelays         []time.Duration `envconfig:"KAFKA_CASCADE_DELAYS" default:""`
	CascadePattern        string          `envconfig:"KAFKA_CASCADE_PATTERN" default:"%s-delay-%d"`
	RetryInitialDelay     time.Duration   `envconfig:"KAFKA_RETRY_INITIAL_DELAY" default:"5s"`
	RetryMaxDelay         time.Duration   `envconfig:"KAFKA_RETRY_MAX_DELAY" default:"24h"`
	RetryFactor           float32         `envconfig:"KAFKA_RETRY_FACTOR" default:"1.5"`
	RetryAddend           time.Duration   `envconfig:"KAFKA_RETRY_ADDEND" default:"0s"`
	RetryMaxAttempts      int             `envconfig:"KAFKA_RETRY_MAX_ATTEMPTS" default:"0"`
	SaramaConfig          *sarama.Config
}

//NewConfig creates a new cloud events config
func NewConfig() *CloudEventsConfig {
	cfg := &CloudEventsConfig{}
	cfg.SaramaConfig = sarama.NewConfig()
	cfg.SaramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	return cfg
}

//LoadFromEnv loads the configuration from environment variables
func (k *CloudEventsConfig) LoadFromEnv() error {
	if err := envconfig.Process("", k); err != nil {
		return err
	}
	version, err := sarama.ParseKafkaVersion(k.KafkaVersion)
	if err != nil {
		return err
	}
	k.SaramaConfig.Version = version
	if k.KafkaRequireSSL {
		k.SaramaConfig.Net.TLS.Enable = true
		// Use the root CAs of the host
		k.SaramaConfig.Net.TLS.Config.RootCAs = nil
	}
	if k.KafkaUsername != "" && k.KafkaPassword != "" {
		k.SaramaConfig.Net.SASL.Enable = true
		k.SaramaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		k.SaramaConfig.Net.SASL.User = k.KafkaUsername
		k.SaramaConfig.Net.SASL.Password = k.KafkaPassword
		k.SaramaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &XDGSCRAMClient{HashGeneratorFcn: SHA512}
		}
	}

	return nil
}

//GetPrefixedTopic returns the fully resolved topic name
func (k *CloudEventsConfig) GetPrefixedTopic() string {
	return k.KafkaTopicPrefix + k.KafkaTopic
}

//GetDeadLettersTopic returns the fully resolved dead letter topic name or an empty string if there is no dead letter topic
func (k *CloudEventsConfig) GetDeadLettersTopic() string {
	if k.KafkaDeadLettersTopic == "" {
		return k.KafkaDeadLettersTopic
	}
	return k.KafkaTopicPrefix + k.KafkaDeadLettersTopic
}

//IsDeadLettersEnabled returns true iff there is a dead letter topic
func (k *CloudEventsConfig) IsDeadLettersEnabled() bool {
	return k.GetDeadLettersTopic() != ""
}

func validateConsumerConfig(config *CloudEventsConfig) error {
	if config.KafkaConsumerGroup == "" {
		return errors.New("consumer group cannot be empty")
	}
	return nil
}
