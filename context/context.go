package context

import (
	"github.com/nicolaferraro/datamesh/context/projection"
	"github.com/nicolaferraro/datamesh/context/processor"
	"github.com/nicolaferraro/datamesh/context/transaction"
	"github.com/nicolaferraro/datamesh/context/initializer"
	"github.com/nicolaferraro/datamesh/notification"
	"github.com/nicolaferraro/datamesh/eventlog"
	"context"
	"github.com/nicolaferraro/datamesh/common"
)

type MeshContextStatus int
const (
	Initializing	MeshContextStatus	= iota
	Ready
)

type MeshContext struct {
	id          string
	eventLog	*eventlog.EventLog
	initializer *initializer.Initializer
	projection  *projection.Projection
	processor   *processor.EventProcessor
	tx          *transaction.TransactionManager
	contextBus  *notification.NotificationBus
	status		MeshContextStatus
}

func NewMeshContext(ctx context.Context, eventLog *eventlog.EventLog, id string) *MeshContext {
	contextBus := notification.NewNotificationBus(ctx)

	proc := processor.NewEventProcessor(ctx, id, contextBus)

	prj := projection.NewProjection(id)
	tx := transaction.NewTransactionManager(ctx, id, prj, contextBus)

	init := initializer.NewInitializer(ctx, id, eventLog, contextBus)

	meshContext := MeshContext{
		id: id,
		eventLog: eventLog,
		initializer: init,
		projection: prj,
		processor: proc,
		tx: tx,
		contextBus: contextBus,
		status: Initializing,
	}

	contextBus.Connect(&meshContext)

	return &meshContext
}

func (ctx *MeshContext) Notify(n notification.Notification) {
	ctx.contextBus.Notify(n)
}

func (ctx *MeshContext) GetDataRetriever() common.DataRetriever {
	return ctx.projection
}

func (ctx *MeshContext) Start() {
	ctx.contextBus.Notify(notification.NewMeshContextStartNotification())
}

func (ctx *MeshContext) OnNotification(n notification.Notification) {
	if n.MeshContextInitializedNotification != nil {
		ctx.status = Ready
	}
}

func (ctx *MeshContext) Status() MeshContextStatus {
	return ctx.status
}