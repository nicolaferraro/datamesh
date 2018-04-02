package mesh

import (
	"github.com/nicolaferraro/datamesh/log"
	"github.com/nicolaferraro/datamesh/projection"
	"path"
	"github.com/nicolaferraro/datamesh/service"
	"github.com/nicolaferraro/datamesh/controller"
)

const (
	LogSubdir	= "log"
)

type Mesh struct {
	dir			string
	log			*log.Log
	projection	*projection.Projection
	server		*service.DefaultDataMeshServer
}

func NewMesh(dir string, port int) (*Mesh, error) {
	eventLog, err := log.NewLog(path.Join(dir, LogSubdir))
	if err != nil {
		return nil, err
	}
	prj := projection.NewProjection()

	ctrl := controller.NewController(prj, eventLog)

	srv := service.NewDefaultDataMeshServer(port, eventLog, ctrl, prj)

	return &Mesh{
		dir: dir,
		log: eventLog,
		projection: prj,
		server: srv,
	}, nil
}

func (msh *Mesh) Start() error {
	return msh.server.Start()
}