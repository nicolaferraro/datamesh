package api


type Event interface {

	Type()		EventType

}


type EventReceiver interface {
	EventInputChannel()		chan<- Event
}

type EventProducer interface {
	EventOutputChannel()	<-chan Event
}