package events

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type EventHandler interface {
	CanHandle(ce cloudevents.Event) bool
	Handle(ce cloudevents.Event) error
}

type EventConsumer interface {
	RegisterHandler(handler EventHandler)
	Start() error
	Stop(ctx context.Context) error
}
