package common

import "github.com/nicolaferraro/datamesh/protobuf"

type EventConsumer interface {

	Consume(*protobuf.Event) 	error

}
