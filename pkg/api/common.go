package api

// Common interfaces

type Named interface {
	Name()	string
}


type EventReceiver interface {
	EventInputChannel()		chan<- Event
}

type EventProducer interface {
	EventOutputChannel()	<-chan Event
}