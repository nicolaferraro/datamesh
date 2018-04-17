package context

import (
	"github.com/nicolaferraro/datamesh/context/initializer"
	"github.com/nicolaferraro/datamesh/context/processor"
	"github.com/nicolaferraro/datamesh/context/projection"
	"github.com/nicolaferraro/datamesh/protobuf"
)

type MeshContext struct {
	init	*initializer.Initializer
	proc	*processor.EventProcessor
	prj		*projection.Projection
	tx		*protobuf.Transaction
}