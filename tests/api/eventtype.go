package api

//go:generate stringer -type=EventType

// Event types
type EventType int

const (
	EventAdd		EventType = iota
	EventRemove
	EventNoop
)