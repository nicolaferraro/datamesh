package core

import (
	"github.com/nicolaferraro/event-db/deleteme/api"
)


// Default struct for Event
type DefaultEvent struct {
	eventType	api.EventType
}

func NewEvent(eventType api.EventType) api.Event {
	event := DefaultEvent{}
	event.eventType = eventType
	return event
}


// Interface methods
func (event DefaultEvent) Type() api.EventType {
	return event.eventType
}