package common

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

type mockApplier struct {
	flag 	bool
	counter	int
}

func (m *mockApplier) ApplyValue(v interface{}) (bool, error) {
	if m.flag {
		m.counter++
	}
	return m.flag, nil
}

func TestSerializer(t *testing.T) {
	mock := mockApplier{}
	serializer := NewSerializer(&mock)

	mock.flag = true

	serializer.Push("")
	time.Sleep(30 * time.Millisecond)

	assert.Equal(t, 1, mock.counter)
	mock.flag = false
	for i:=0; i<100; i++ {
		serializer.Push("")
	}
	time.Sleep(30 * time.Millisecond)

	assert.Equal(t, 1, mock.counter)
	mock.flag = true
	serializer.OnNotification(Notification{
		Type: NotificationTypeEventPushed,
	})
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 101, mock.counter)
}

