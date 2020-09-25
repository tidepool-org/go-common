package events

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/tidepool-org/go-common/clients/shoreline"
)

const (
	DeleteUserEventType = "users:delete"
	UpdateUserEventType = "users:update"
	CreateUserEventType = "users:create"
)

type Event interface {
	GetEventType() string
}

var _ Event = DeleteUserEvent{}

type DeleteUserEvent struct {
	shoreline.UserData `json:",inline"`
}

func (d DeleteUserEvent) GetEventType() string {
	return DeleteUserEventType
}

var _ Event = CreateUserEvent{}

type CreateUserEvent struct {
	shoreline.UserData `json:",inline"`
}

func (d CreateUserEvent) GetEventType() string {
	return CreateUserEventType
}

var _ Event = UpdateUserEvent{}

type UpdateUserEvent struct {
	Original shoreline.UserData `json:"original"`
	Updated  shoreline.UserData `json:"updated"`
}

func (d UpdateUserEvent) GetEventType() string {
	return UpdateUserEventType
}

type UserEventsHandler interface {
	HandleUpdateUserEvent(payload UpdateUserEvent)
	HandleCreateUserEvent(payload CreateUserEvent)
	HandleDeleteUserEvent(payload DeleteUserEvent)
}

var _ EventHandler = &DelegatingUserEventsHandler{}

type DelegatingUserEventsHandler struct {
	delegate UserEventsHandler
}

func (d *DelegatingUserEventsHandler) CanHandle(ce cloudevents.Event) bool {
	switch ce.Type() {
	case CreateUserEventType, UpdateUserEventType, DeleteUserEventType:
		return true
	default:
		return false
	}
}

func (d *DelegatingUserEventsHandler) Handle(ce cloudevents.Event) error {
	switch ce.Type() {
	case CreateUserEventType:
		payload := CreateUserEvent{}
		if err := ce.DataAs(&payload); err != nil {
			return err
		}
		d.delegate.HandleCreateUserEvent(payload)
	case UpdateUserEventType:
		payload := UpdateUserEvent{}
		if err := ce.DataAs(&payload); err != nil {
			return err
		}
		d.delegate.HandleUpdateUserEvent(payload)
	case DeleteUserEventType:
		payload := DeleteUserEvent{}
		if err := ce.DataAs(&payload); err != nil {
			return err
		}
		d.delegate.HandleDeleteUserEvent(payload)
	}
	return nil
}

type NoopUserEventsHandler struct{}

var _ UserEventsHandler = &NoopUserEventsHandler{}

func (d *NoopUserEventsHandler) HandleUpdateUserEvent(payload UpdateUserEvent) {}
func (d *NoopUserEventsHandler) HandleCreateUserEvent(payload CreateUserEvent) {}
func (d *NoopUserEventsHandler) HandleDeleteUserEvent(payload DeleteUserEvent) {}
