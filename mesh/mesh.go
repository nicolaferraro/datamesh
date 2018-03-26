package mesh

import (
	"github.com/nicolaferraro/datamesh/log"
	"github.com/nicolaferraro/datamesh/projection"
	"path"
	"github.com/nicolaferraro/datamesh/service"
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
	srv := service.NewDefaultDataMeshServer(port, eventLog)

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