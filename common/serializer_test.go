package common

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"context"
)

type mockApplier struct {
	flag 	bool
	counter	int
}

func (m *mockApplier) ExecuteSerially(v interface{}) (bool, error) {
	if m.flag {
		m.counter++
	}
	return m.flag, nil
}

func TestSerializer(t *testing.T) {
	ctx := context.Background()
	mock := mockApplier{}
	serializer := NewSerializer(ctx, &mock)

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
	serializer.OnNotification("example")
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 101, mock.counter)
}

