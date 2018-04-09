package eventlog

import (
	"testing"
	"os"
	"strconv"
	"reflect"
	"context"
	"github.com/nicolaferraro/datamesh/notification"
)

const testDir = "../.testdata/log"


func TestLogBasicUsage(t *testing.T) {
	os.RemoveAll(testDir)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := notification.NewNotificationBus(ctx)

	eventLog, err := NewEventLog(ctx, testDir, bus)
	if err != nil {
		t.Fatal("Cannot create log:", err)
	}

	for i:=1; i<=20; i++ {
		record := []byte("Record" + strconv.Itoa(i))
		num, err := eventLog.AppendRaw(record)
		if err != nil {
			t.Fatal("Cannot append:", err)
		}
		if num != uint64(i) {
			t.Fatal("Wrong number of entries. Expected:", i, "Got:", num)
		}

		if !reflect.DeepEqual(record, []byte("Record" + strconv.Itoa(i))) {
			t.Fatal("Record changed")
		}
	}

	if err := eventLog.Sync(); err != nil {
		t.Fatal("Cannot fsync:", err)
	}

}