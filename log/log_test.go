package log

import (
	"testing"
	"os"
	"strconv"
	"reflect"
)

const testDir = "../.testdata/log"


func TestLogBasicUsage(t *testing.T) {
	os.RemoveAll(testDir)
	eventLog, err := NewLog(testDir)
	if err != nil {
		t.Fatal("Cannot create log:", err)
	}

	for i:=1; i<=20; i++ {
		record := []byte("Record" + strconv.Itoa(i))
		num, err := eventLog.Append(record)
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