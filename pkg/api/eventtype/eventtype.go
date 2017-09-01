package eventtype

//go:generate stringer -type=EventType

// Event types
type EventType int

const (
	Add		EventType = iota
	Remove
	Noop
)

type Event interface {

	Type()		EventType

}