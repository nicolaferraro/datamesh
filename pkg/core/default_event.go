package core

import (
	"github.com/nicolaferraro/event-db/pkg/api"
	"github.com/nicolaferraro/event-db/pkg/api/eventtype"
)


// Default struct for Event
type DefaultEvent struct {
	eventType	eventtype.EventType
}

func NewEvent(eventType eventtype.EventType) api.Event {
	event := DefaultEvent{}
	event.eventType = eventType
	return event
}


// Interface methods
func (event DefaultEvent) Type() eventtype.EventType {
	return event.eventType
}