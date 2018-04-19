package initializer

import (
	"context"
	"github.com/nicolaferraro/datamesh/eventlog"
	"github.com/nicolaferraro/datamesh/notification"
	"golang.org/x/sync/semaphore"
	"github.com/golang/glog"
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
}

func NewInitializer(ctx context.Context, eventlog *eventlog.EventLog, bus *notification.NotificationBus) *Initializer {
	init := Initializer{
		ctx: ctx,
		eventlog: eventlog,
		bus: bus,
	}
	bus.Connect(&init)

	init.semaphore = semaphore.NewWeighted(PrefetchSize)

	return &init
}

func (init *Initializer) OnNotification(n notification.Notification) {
	if init.terminated {
		return
	} else if n.MeshStartNotification != nil {
		if !init.started {
			init.started = true
			go init.run()
		}
	} else if n.TransactionProcessedNotification != nil {
		init.semaphore.Release(1)
		ver := n.TransactionProcessedNotification.Version
		init.currentVersion = ver

		if ver == init.targetVersion {
			// End
			go init.initialized()
		}
	}
}

func (init *Initializer) run() {
	glog.Info("Data Mesh projection initialization started")

	reader, err := init.eventlog.NewReader()
	if err != nil {
		glog.Error("Error during projection initialization: ", err)
		return
		// TODO handle better
	}

	init.targetVersion = reader.Top
	if init.targetVersion == 0 {
		init.initialized()
	} else {
		for v := uint64(1); v <= init.targetVersion; v++ {
			init.semaphore.Acquire(init.ctx, 1)
			evt, err := reader.Next()
			if err != nil {
				glog.Error("Error while replaying event log initialization: ", err)
				return
				// TODO handle better
			}
			if evt.Version == 0 {
				evt.Version = v
			}
			init.bus.Notify(notification.NewEventAppendedNotification(evt, true))
		}
	}
}

func (init *Initializer) initialized() {
	if !init.terminated {
		init.terminated = true
		glog.Infof("Data Mesh projection initialized at version %d\n", init.targetVersion)
		init.bus.Notify(notification.NewMeshInitializedNotification())
	}
}
