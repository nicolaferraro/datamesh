package common

import "github.com/nicolaferraro/datamesh/protobuf"

type EventObserver interface {

	Accept(*protobuf.Event) 	error

}