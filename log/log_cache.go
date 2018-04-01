package log

import (
	"github.com/nicolaferraro/datamesh/protobuf"
)

type LogCache struct {
	events	[]*protobuf.Event
}

func NewLogCache() *LogCache {
	return &LogCache{}
}

/*
 * implements common.MessageObserver
 */
func (cache *LogCache) Accept(evt *protobuf.Event) error {
	cache.events = append(cache.events, evt)
	return nil
}

func (cache *LogCache) Prune(version int64) {
	var index int
	events := cache.events

	for index = 0; index < len(events) && events[index].Version < version; index++ {
	}

	cache.events = events[index:]
}

func (cache *LogCache) Get(clientIdentifier string) *protobuf.Event {
	for _, evt := range cache.events {
		if evt.ClientIdentifier == clientIdentifier {
			return evt
		}
	}
	return nil
}
