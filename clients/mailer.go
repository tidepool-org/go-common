package clients

import (
	"context"
	"github.com/tidepool-org/go-common/events"
)

const MailerTopic = "emails"

type MailerClient struct {
	producer *events.KafkaCloudEventsProducer
}

func NewMailerClient() (*MailerClient, error) {
	config := &events.CloudEventsConfig{}
	if err := config.LoadFromEnv(); err != nil {
		return nil, err
	}

	config.KafkaTopic = MailerTopic
	producer, err := events.NewKafkaCloudEventsProducer(config)
	if err != nil {
		return nil, err
	}

	return &MailerClient{producer}, nil
}

func (m *MailerClient) SendEmailTemplate(ctx context.Context, event events.SendEmailTemplateEvent) error {
	return m.producer.Send(ctx, event)
}
