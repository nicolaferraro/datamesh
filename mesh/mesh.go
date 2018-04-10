package mesh

import (
	"github.com/nicolaferraro/datamesh/projection"
	"path"
	"github.com/nicolaferraro/datamesh/service"
	"github.com/nicolaferraro/datamesh/notification"
	"context"
	"github.com/nicolaferraro/datamesh/eventlog"
	"github.com/nicolaferraro/datamesh/processor"
	"github.com/nicolaferraro/datamesh/transaction"
	"github.com/nicolaferraro/datamesh/initializer"
)

const (
	LogSubdir	= "log"
)

type Mesh struct {
	dir			string
	eventLog	*eventlog.EventLog
	projection	*projection.Projection
	processor   *processor.EventProcessor
	tx			*transaction.TransactionManager
	server		*service.DefaultDataMeshServer
	initializer	*initializer.Initializer
	bus			*notification.NotificationBus
}

func NewMesh(ctx context.Context, dir string, port int) (*Mesh, error) {
	bus := notification.NewNotificationBus(ctx)
	log, err := eventlog.NewEventLog(ctx, path.Join(dir, LogSubdir), bus)
	if err != nil {
		return nil, err
	}

	proc := processor.NewEventProcessor(ctx, bus)

	prj := projection.NewProjection()
	tx := transaction.NewTransactionManager(ctx, prj, bus)

	init := initializer.NewInitializer(ctx, log, bus)

	srv := service.NewDefaultDataMeshServer(port, bus, log, prj)

	mesh := Mesh{
		dir: dir,
		eventLog: log,
		processor: proc,
		projection: prj,
		tx: tx,
		server: srv,
		initializer: init,
		bus: bus,
	}

	return &mesh, nil
}

func (msh *Mesh) Start() error {
	msh.bus.Notify(notification.NewMeshStartNotification())
	return msh.server.Start()
}