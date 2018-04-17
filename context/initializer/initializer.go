package initializer

import (
	"context"
	"github.com/nicolaferraro/datamesh/eventlog"
	"github.com/nicolaferraro/datamesh/notification"
	"log"
	"github.com/nicolaferraro/datamesh/common"
	"golang.org/x/sync/semaphore"
)

const (
	PrefetchSize = 20
)

type Initializer struct {
	ctx				context.Context
	started			bool
	terminated		bool
	eventlog		*eventlog.EventLog
	bus				*notification.NotificationBus
	semaphore		*semaphore.Weighted
	targetVersion	uint64
	currentVersion	uint64
	serializer		*common.Serializer
}

func NewInitializer(ctx context.Context, eventlog *eventlog.EventLog, bus *notification.NotificationBus) *Initializer {
	init := Initializer{
		ctx: ctx,
		eventlog: eventlog,
		bus: bus,
	}
	init.serializer = common.NewSerializer(ctx, &init)
	bus.Connect(&init)

	init.semaphore = semaphore.NewWeighted(PrefetchSize)

	return &init
}

func (init *Initializer) OnNotification(n notification.Notification) {
	if init.terminated {
		return
	} else if n.MeshStartNotification != nil {
		init.serializer.Push(n)
	} else if n.TransactionProcessedNotification != nil {
		init.semaphore.Release(1)
		init.serializer.Push(n)
	}
}

func (init *Initializer) ExecuteSerially(value interface{}) (bool, error) {
	if n, ok := value.(notification.Notification); ok {
		if n.MeshStartNotification != nil && !init.started {
			init.started = true

			log.Println("data mesh projection initialization started")

			reader, err := init.eventlog.NewReader()
			if err != nil {
				// TODO handle
				return false, err
			}

			init.targetVersion = reader.Top
			if init.targetVersion == 0 {
				init.initialized()
			} else {
				for v := uint64(1); v <= init.targetVersion; v++ {
					init.semaphore.Acquire(init.ctx, 1)
					evt, err := reader.Next()
					if err != nil {
						log.Println("Fatal error: ", err) // TODO manage
					}
					if evt.Version == 0 {
						evt.Version = v
					}
					init.bus.Notify(notification.NewEventAppendedNotification(evt, true))
				}
			}

		} else if n.TransactionProcessedNotification != nil {
			ver := n.TransactionProcessedNotification.Version
			init.currentVersion = ver

			if ver == init.targetVersion {
				// End
				init.initialized()
			}
		}
	}
	return true, nil
}

func (init *Initializer) initialized() {
	if !init.terminated {
		init.terminated = true
		log.Printf("data mesh projection initialized at version %d\n", init.targetVersion)
		init.bus.Notify(notification.NewMeshInitializedNotification())
	}
}
