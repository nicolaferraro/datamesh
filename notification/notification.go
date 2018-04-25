package notification

import (
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/nicolaferraro/datamesh/context/projection"
	"github.com/nicolaferraro/datamesh/common"
)

// Container for all notifications
type Notification struct {
	MeshContextStartNotification		*MeshContextStartNotification
	MeshContextInitializedNotification	*MeshContextInitializedNotification
	EventAppendedNotification			*EventAppendedNotification
	TransactionReceivedNotification		*TransactionReceivedNotification
	TransactionProcessedNotification	*TransactionProcessedNotification
	TransactionFailedNotification		*TransactionFailedNotification
	ClientConnectedNotification			*ClientConnectedNotification
	ClientDisconnectedNotification		*ClientDisconnectedNotification
}

// Signals that all pieces of a mesh context have been connected
type MeshContextStartNotification struct {
}

func NewMeshContextStartNotification() Notification {
	return Notification{
		MeshContextStartNotification: &MeshContextStartNotification{},
	}
}

// Signals that a mesh context has been completely initialized (e.g. projections aligned)
type MeshContextInitializedNotification struct {
}

func NewMeshContextInitializedNotification() Notification {
	return Notification{
		MeshContextInitializedNotification: &MeshContextInitializedNotification{},
	}
}

// Signals that a event has been appended to the event log
type EventAppendedNotification struct {
	Event	*protobuf.Event
	Replay	bool
}

func NewEventAppendedNotification(event *protobuf.Event, replay bool) Notification {
	return Notification{
		EventAppendedNotification: &EventAppendedNotification{
			Event: event,
			Replay: replay,
		},
	}
}

// Signals that a certain transaction has been received
type TransactionReceivedNotification struct {
	Transaction	*protobuf.Transaction
}

func NewTransactionReceivedNotification(transaction *protobuf.Transaction) Notification {
	return Notification{
		TransactionReceivedNotification: &TransactionReceivedNotification{
			Transaction: transaction,
		},
	}
}

// Signals that a certain transaction number has been processed correctly
type TransactionProcessedNotification struct {
	Projection	*projection.Projection
	Version		uint64
	Error		error
}

func NewTransactionProcessedNotification(projection *projection.Projection, version uint64) Notification {
	return Notification{
		TransactionProcessedNotification: &TransactionProcessedNotification{
			Projection: projection,
			Version: version,
		},
	}
}

// Signals that a certain transaction number has failed
type TransactionFailedNotification struct {
	Projection	*projection.Projection
	Event		*protobuf.Event
	Error		error
}

func NewTransactionFailedNotification(projection *projection.Projection, event *protobuf.Event, error error) Notification {
	return Notification{
		TransactionFailedNotification: &TransactionFailedNotification{
			Projection: projection,
			Event: event,
			Error: error,
		},
	}
}

// Signal a client connection
type ClientConnectedNotification struct {
	Client common.EventConsumer
}

func NewClientConnectedNotification(client common.EventConsumer) Notification {
	return Notification{
		ClientConnectedNotification: &ClientConnectedNotification{
			Client: client,
		},
	}
}

// Signal a client disconnection
type ClientDisconnectedNotification struct {
	Client common.EventConsumer
}

func NewClientDisconnectedNotification(client common.EventConsumer) Notification {
	return Notification{
		ClientDisconnectedNotification: &ClientDisconnectedNotification{
			Client: client,
		},
	}
}
