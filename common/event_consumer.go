package common

import "github.com/nicolaferraro/datamesh/protobuf"

type EventConsumer interface {

	Consume(*protobuf.Event) 	error

}

type CloseableEventConsumer interface {
	EventConsumer
	Close()
	IsClosed() bool
}

type EventConsumerController interface {
	ConnectEventConsumer(CloseableEventConsumer)
}