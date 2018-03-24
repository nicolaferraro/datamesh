package log

import (
	"testing"
	"os"
	"strconv"
)

const testDir = "../.testdata/log"


func TestLogBasicUsage(t *testing.T) {
	os.RemoveAll(testDir)
	eventLog, err := NewLog(testDir)
	if err != nil {
		t.Fatal("Cannot create log:", err)
	}

	for i:=1; i<=20; i++ {
		num, err := eventLog.Append([]byte("Record" + strconv.Itoa(i)))
		if err != nil {
			t.Fatal("Cannot append:", err)
		}
		if num != uint64(i) {
			t.Fatal("Wrong number of entries. Expected:", i, "Got:", num)
		}
	}

	if err := eventLog.Sync(); err != nil {
		t.Fatal("Cannot fsync:", err)
	}

}