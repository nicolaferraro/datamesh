package mesh

import (
	"github.com/nicolaferraro/datamesh/log"
	"github.com/nicolaferraro/datamesh/projection"
	"path"
)

const (
	LogSubdir	= "log"
)

type Mesh struct {
	dir			string
	log			*log.Log
	projection	*projection.Projection
}

func NewMesh(dir string) (*Mesh, error) {
	eventLog, err := log.NewLog(path.Join(dir, LogSubdir))
	if err != nil {
		return nil, err
	}
	prj := projection.NewProjection()

	return &Mesh{
		dir: dir,
		log: eventLog,
		projection: prj,
	}, nil
}