package service

import (
	"google.golang.org/grpc"
	"strconv"
	"github.com/nicolaferraro/datamesh/protobuf"
)

func NewDataMeshClientConnection(host string, port int) (protobuf.DataMeshClient, error) {
	conn, err := grpc.Dial(host + ":" + strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return protobuf.NewDataMeshClient(conn), nil
}
