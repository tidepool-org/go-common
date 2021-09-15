package clients

import (
	"context"
	"github.com/tidepool-org/go-common/events"
)

const MailerTopic = "emails"

type MailerClient struct {
	producer *events.KafkaCloudEventsProducer
}

func NewMailerClient(config *events.CloudEventsConfig) (*MailerClient, error) {
	producer, err := events.NewKafkaCloudEventsProducer(config)
	if err != nil {
		return nil, err
	}

	return &MailerClient{producer}, nil
}

func (m *MailerClient) SendEmailTemplate(ctx context.Context, event events.SendEmailTemplateEvent) error {
	return m.producer.Send(ctx, event)
}
