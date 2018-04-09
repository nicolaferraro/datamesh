package notification

import (
	"testing"
	"context"
	"github.com/stretchr/testify/assert"
)

func TestSimpleBus(t *testing.T) {

	ctx := context.Background()
	bus := NewNotificationBus(ctx)

	mock1 := newMockNotificationReceiver()
	bus.Connect(mock1)

	mock2 := newMockNotificationReceiver()
	bus.Connect(mock2)

	bus.Notify(NewMeshStartNotification())
	bus.Notify(NewMeshInitializedNotification())

	assert.Nil(t, mock1.Wait(2))
	assert.Nil(t, mock2.Wait(2))

	for _, mock := range []*mockNotificationReceiver{mock1, mock2} {
		assert.Equal(t, 2, len(mock.notifications))
		assert.Equal(t, NewMeshStartNotification(), mock.notifications[0])
		assert.Equal(t, NewMeshInitializedNotification(), mock.notifications[1])
	}

}