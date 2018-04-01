package common

type Observer interface {
	OnNotification(Notification) error
}

type Observable struct {
	observers	[]Observer
}

func (obs *Observable) Listen(observer Observer) {
	obs.observers = append(obs.observers, observer)
}

func (obs *Observable) Notify(notification Notification) error {
	for _, observer := range obs.observers {
		if err := observer.OnNotification(notification); err != nil {
			return err
		}
	}
	return nil
}
