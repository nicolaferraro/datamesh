package service

import (
	"path"
	"github.com/nicolaferraro/datamesh/notification"
	"context"
	"github.com/nicolaferraro/datamesh/eventlog"
	meshcontext "github.com/nicolaferraro/datamesh/context"
	"sync"
	"strconv"
	"github.com/golang/glog"
)

const (
	LogSubdir	= "log"
)

type Mesh struct {
	ctx			context.Context
	dir			string
	eventLog	*eventlog.EventLog
	server		*DefaultDataMeshServer
	mainBus		*notification.NotificationBus
	contexts	map[string]*meshcontext.MeshContext
	mutex		sync.Mutex
}

func NewMesh(ctx context.Context, dir string, port int) (*Mesh, error) {
	mainBus := notification.NewNotificationBus(ctx)
	log, err := eventlog.NewEventLog(ctx, path.Join(dir, LogSubdir), mainBus)
	if err != nil {
		return nil, err
	}

	mesh := Mesh{
		ctx: ctx,
		dir: dir,
		eventLog: log,
		mainBus: mainBus,
		contexts: make(map[string]*meshcontext.MeshContext),
	}

	srv := NewDefaultDataMeshServer(port, log, &mesh)
	mesh.server = srv

	mainBus.Connect(&mesh)

	return &mesh, nil
}

func (msh *Mesh) OnNotification(n notification.Notification) {
	// Forward notifications to all context buses
	contexts := msh.contexts
	for _, ctx := range contexts {
		ctx.Notify(n)
	}
}

func (msh *Mesh) GetContext(name string, revision uint64) *meshcontext.MeshContext {
	id := contextId(name, revision)
	ctx := msh.contexts[id]
	if ctx != nil {
		return ctx
	}

	// I hope double check is not buggy also in go...
	msh.mutex.Lock()
	defer msh.mutex.Unlock()

	ctx = msh.contexts[id]
	if ctx != nil {
		return ctx
	}

	glog.Infof("Context %s is not present. Creating...", id)
	meshContext := meshcontext.NewMeshContext(msh.ctx, msh.eventLog, id)
	newContextMap := make(map[string]*meshcontext.MeshContext)
	for k,v := range msh.contexts {
		newContextMap[k] = v
	}
	newContextMap[id] = meshContext

	// Replace the mesh context map and start
	msh.contexts = newContextMap
	meshContext.Start()
	return meshContext
}

func (msh *Mesh) Start() error {
	return msh.server.Start()
}

func contextId(name string, revision uint64) string {
	return name + "/" + strconv.FormatUint(revision, 10)
}