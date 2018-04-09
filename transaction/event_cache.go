package transaction

import (
	"github.com/nicolaferraro/datamesh/protobuf"
)

type EventCache struct {
	events	[]*protobuf.Event
}

func NewEventCache() *EventCache {
	return &EventCache{}
}

func (cache *EventCache) Register(evt *protobuf.Event) error {
	cache.events = append(cache.events, evt)
	return nil
}

func (cache *EventCache) Prune(version uint64) {
	var index int
	events := cache.events

	for index = 0; index < len(events) && events[index].Version < version; index++ {
	}

	cache.events = events[index:]
}

func (cache *EventCache) Get(clientIdentifier string) *protobuf.Event {
	for _, evt := range cache.events {
		if evt.ClientIdentifier == clientIdentifier {
			return evt
		}
	}
	return nil
}
