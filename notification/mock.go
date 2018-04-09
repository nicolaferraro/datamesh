package notification

import (
	"time"
	"errors"
)

type mockNotificationReceiver struct {
	notifications	[]Notification
	gotData			chan bool
}

func newMockNotificationReceiver() *mockNotificationReceiver {
	return &mockNotificationReceiver{
		gotData: make(chan bool, 1),
	}
}

func (mock *mockNotificationReceiver) OnNotification(n Notification) {
	mock.notifications = append(mock.notifications, n)
	select {
	case mock.gotData <- true:
	default:
	}
}

func (mock *mockNotificationReceiver) Wait(num int) error {
	for {
		if len(mock.notifications) == num {
			return nil
		}
		select {
			case <- mock.gotData:
			case <- time.After(5 * time.Second):
				return errors.New("Timeout")
		}
	}
}
