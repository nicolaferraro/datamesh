package controller

import (
	"testing"
	"github.com/nicolaferraro/datamesh/log"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/nicolaferraro/datamesh/projection"
	"github.com/nicolaferraro/datamesh/protobuf"
	"encoding/json"
	"time"
)

const testDir = "../.testdata/log"

func initController(t *testing.T) *Controller {
	os.RemoveAll(testDir)
	eventLog, err := log.NewLog(testDir)
	assert.Nil(t, err)
	prj := projection.NewProjection()
	notifier := NewNotifier()
	return NewController(prj, eventLog, notifier)
}

func TestBasicController(t *testing.T) {
	ctrl := initController(t)

	evt := protobuf.Event{
		ClientIdentifier: "ID",
		Name: "evt",
	}

	assert.Nil(t, ctrl.log.Consume(&evt))


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

	assert.Nil(t, ctrl.Apply(&transaction))

	var serData interface{}
	assert.Nil(t, json.Unmarshal(data, &serData))

	time.Sleep(50 * time.Millisecond)
	_, retrData, err := ctrl.projection.Get("x")
	assert.Nil(t, err)
	assert.Equal(t, serData, retrData)
}

func TestOperationCombo(t *testing.T) {
	ctrl := initController(t)

	evt := protobuf.Event{
		ClientIdentifier: "ID",
		Name: "evt",
		Version: 1,
	}

	assert.Nil(t, ctrl.log.Consume(&evt))


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

	assert.Nil(t, ctrl.Apply(&transaction))

	var serData interface{}
	assert.Nil(t, json.Unmarshal(finalData, &serData))

	time.Sleep(50 * time.Millisecond)
	_, retrData, err := ctrl.projection.Get("x")
	assert.Nil(t, err)
	assert.Equal(t, serData, retrData)
}

func TestOperationOverride(t *testing.T) {
	ctrl := initController(t)

	evt := protobuf.Event{
		ClientIdentifier: "ID",
		Name: "evt",
		Version: 1,
	}

	assert.Nil(t, ctrl.log.Consume(&evt))


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

	assert.Nil(t, ctrl.Apply(&transaction))

	var serData interface{}
	assert.Nil(t, json.Unmarshal(finalData, &serData))

	time.Sleep(50 * time.Millisecond)
	_, retrData, err := ctrl.projection.Get("x")
	assert.Nil(t, err)
	assert.Equal(t, serData, retrData)
}

func TestDelayedEvent(t *testing.T) {
	ctrl := initController(t)

	evt := protobuf.Event{
		ClientIdentifier: "ID",
		Name: "evt",
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

	assert.Nil(t, ctrl.Apply(&transaction))

	// Now push the event
	assert.Nil(t, ctrl.log.Consume(&evt))

	var serData interface{}
	assert.Nil(t, json.Unmarshal(data, &serData))

	time.Sleep(50 * time.Millisecond)
	_, retrData, err := ctrl.projection.Get("x")
	assert.Nil(t, err)
	assert.Equal(t, serData, retrData)
}
