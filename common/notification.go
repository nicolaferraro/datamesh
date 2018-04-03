package common

type NotificationType int

const (
	NotificationTypeProjectionVersion	NotificationType = iota
	NotificationTypeEventPushed			NotificationType = iota
)

type Notification struct {
	Type	NotificationType
	Payload	interface{}
}
