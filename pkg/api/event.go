package api

import "github.com/nicolaferraro/event-db/pkg/api/eventtype"

type Event interface {

	Type()		eventtype.EventType

}