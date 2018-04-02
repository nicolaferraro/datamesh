package controller

import (
	"testing"
	"github.com/nicolaferraro/datamesh/log"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/nicolaferraro/datamesh/projection"
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/golang/protobuf/ptypes/struct"
)

const testDir = "../.testdata/log"

func TestBasicController(t *testing.T) {
	os.RemoveAll(testDir)
	eventLog, err := log.NewLog(testDir)
	assert.Nil(t, err)
	prj := projection.NewProjection()
	ctrl := NewController(prj, eventLog)

	evt := protobuf.Event{
		ClientIdentifier: "ID",
		Name: "evt",
	}

	assert.Nil(t, eventLog.Consume(&evt))

	world := &structpb.Struct{map[string]*structpb.Value{
		"Wor": {Kind: &structpb.Value_StringValue{"ld!"}},
	}}

	transaction := protobuf.Transaction{
		Event: &evt,
		Operations: []*protobuf.Operation{
			{&protobuf.Operation_Upsert{&protobuf.UpsertOperation{&protobuf.Data{
				Path: &protobuf.Path{"hello", 1},
				Content: world,
			}}}},
		},
	}

	assert.Nil(t, ctrl.Apply(&transaction))

	//a,_ := prj.Get("hello")
	//println(a)

	// TODO better conversion of protobuf objects

}