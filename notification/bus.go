package notification

import "context"

const NotificationBusBufferSize = 512

// Bridges publishers and observers inside the application
type NotificationBus struct {
	OutOfOrder		bool
	listeners		[]NotificationBusListener
	notifications	chan Notification
}

type NotificationBusListener interface {
	OnNotification(Notification)
}

func NewNotificationBus(ctx context.Context) *NotificationBus {
	bus := NotificationBus{
		notifications: make(chan Notification, NotificationBusBufferSize),
	}

	go func() {

		for {
			select {
				case <- ctx.Done():
					return
				case notification := <- bus.notifications:
					for _, listener := range bus.listeners {
						listener.OnNotification(notification)
					}
			}
		}

	}()

	return &bus
}

func (bus *NotificationBus) Connect(listener NotificationBusListener) {
	bus.listeners = append(bus.listeners, listener)
}

func (bus *NotificationBus) Notify(notification Notification) {
	if bus.OutOfOrder {
		select {
		case bus.notifications <- notification:
		default:
			// buffer full. Do it asynchronously
			go func() {
				bus.notifications <- notification
			}()
		}
	} else {
		bus.notifications <- notification
	}

}
