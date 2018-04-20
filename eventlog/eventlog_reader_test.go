package eventlog

import (
	"testing"
	"os"
	"strconv"
	"context"
	"github.com/nicolaferraro/datamesh/notification"
	"github.com/stretchr/testify/assert"
	"github.com/nicolaferraro/datamesh/protobuf"
	"time"
)


func TestLogReader(t *testing.T) {
	os.RemoveAll(testDir)
	ctx := context.Background()
	bus := notification.NewNotificationBus(ctx)

	eventLog, err := NewEventLog(ctx, testDir, bus)
	assert.Nil(t, err)

	num := 20
	for i:=1; i<=num; i++ {
		record := []byte("Record" + strconv.Itoa(i))
		evt := protobuf.Event{
			Name: "Record",
			ClientIdentifier: strconv.Itoa(i),
			Payload: record,
		}

		assert.Nil(t, eventLog.Consume(&evt))
	}

	time.Sleep(250 * time.Millisecond)
	assert.Nil(t, eventLog.Sync())

	reader, err := eventLog.NewReader()
	assert.Nil(t, err)
	assert.Equal(t, num, int(reader.Top))

	for i:=1; i<=int(reader.Top); i++ {
		evt, err := reader.Next()
		assert.Nil(t, err)

		expected := []byte("Record" + strconv.Itoa(i))
		assert.Equal(t, expected, evt.Payload)
	}

}
