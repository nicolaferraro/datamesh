package transaction

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/nicolaferraro/datamesh/projection"
	"github.com/nicolaferraro/datamesh/protobuf"
	"encoding/json"
	"time"
	"context"
	"github.com/nicolaferraro/datamesh/notification"
)


func initTransactionManager() *TransactionManager {
	ctx := context.Background()
	prj := projection.NewProjection()
	bus := notification.NewNotificationBus(ctx)
	tx := NewTransactionManager(ctx, prj, bus)
	return tx
}


func TestBasicController(t *testing.T) {
	tx := initTransactionManager()

	evt := protobuf.Event{
		ClientIdentifier: "ID",
		Name: "evt",
		Version: 1,
	}
	tx.bus.Notify(notification.NewEventAppendedNotification(&evt, false))

	data := []byte(`{
		"name": "Hello",
		"value": "World!"
	}`)

	transaction := protobuf.Transaction{
		Event: &evt,
		Operations: []*protobuf.Operation{
			protobuf.NewUpsertOperation("x", 1, data),
		},
	}

	tx.bus.Notify(notification.NewTransactionReceivedNotification(&transaction))

	var serData interface{}
	assert.Nil(t, json.Unmarshal(data, &serData))

	time.Sleep(50 * time.Millisecond)
	_, retrData, err := tx.projection.Get("x")
	assert.Nil(t, err)
	assert.Equal(t, serData, retrData)
}

func TestOperationCombo(t *testing.T) {
	tx := initTransactionManager()

	evt := protobuf.Event{
		ClientIdentifier: "ID",
		Name: "evt",
		Version: 1,
	}

	tx.bus.Notify(notification.NewEventAppendedNotification(&evt, false))


	data := []byte(`{
		"name": "Hello",
		"value": "World!"
	}`)

	finalData := []byte(`{
		"value": "World!"
	}`)

	transaction := protobuf.Transaction{
		Event: &evt,
		Operations: []*protobuf.Operation{
			protobuf.NewUpsertOperation("x", 1, data),
			protobuf.NewDeleteOperation("x.name", 2),
		},
	}

	tx.bus.Notify(notification.NewTransactionReceivedNotification(&transaction))

	var serData interface{}
	assert.Nil(t, json.Unmarshal(finalData, &serData))

	time.Sleep(50 * time.Millisecond)
	_, retrData, err := tx.projection.Get("x")
	assert.Nil(t, err)
	assert.Equal(t, serData, retrData)
}

func TestOperationOverride(t *testing.T) {
	tx := initTransactionManager()

	evt := protobuf.Event{
		ClientIdentifier: "ID",
		Name: "evt",
		Version: 1,
	}

	tx.bus.Notify(notification.NewEventAppendedNotification(&evt, false))


	data1 := []byte(`{
		"name": "1",
		"surname": "2",
		"value": "3"
	}`)

	data2 := []byte(`{
		"name": "11",
		"surname": "12",
		"value": "13",
		"other": "14"
	}`)

	data3 := []byte(`{
		"name": "Hello",
		"value": "World!"
	}`)

	finalData := []byte(`{
		"value": "World!"
	}`)

	transaction := protobuf.Transaction{
		Event: &evt,
		Operations: []*protobuf.Operation{
			protobuf.NewUpsertOperation("x", 1, data1),
			protobuf.NewUpsertOperation("x", 1, data2),
			protobuf.NewUpsertOperation("x", 1, data3),
			protobuf.NewDeleteOperation("x.name", 2),
		},
	}

	tx.bus.Notify(notification.NewTransactionReceivedNotification(&transaction))

	var serData interface{}
	assert.Nil(t, json.Unmarshal(finalData, &serData))

	time.Sleep(50 * time.Millisecond)
	_, retrData, err := tx.projection.Get("x")
	assert.Nil(t, err)
	assert.Equal(t, serData, retrData)
}

func TestDelayedEvent(t *testing.T) {
	tx := initTransactionManager()

	evt := protobuf.Event{
		ClientIdentifier: "ID",
		Name: "evt",
	}
	pEvt := protobuf.Event{
		ClientIdentifier: evt.ClientIdentifier,
		Name: evt.Name,
		Version: 1,
	}

	data := []byte(`{
		"name": "Hello",
		"value": "World!"
	}`)

	transaction := protobuf.Transaction{
		Event: &evt,
		Operations: []*protobuf.Operation{
			protobuf.NewUpsertOperation("x", 1, data),
		},
	}

	tx.bus.Notify(notification.NewTransactionReceivedNotification(&transaction))
	time.Sleep(50 * time.Millisecond)

	// Now push the event
	tx.bus.Notify(notification.NewEventAppendedNotification(&pEvt, false))

	var serData interface{}
	assert.Nil(t, json.Unmarshal(data, &serData))

	time.Sleep(50 * time.Millisecond)
	_, retrData, err := tx.projection.Get("x")
	assert.Nil(t, err)
	assert.Equal(t, serData, retrData)
}
