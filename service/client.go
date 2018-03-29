package service

import (
	"google.golang.org/grpc"
	"strconv"
)

func NewDataMeshClientConnection(host string, port int) (DataMeshClient, error) {
	conn, err := grpc.Dial(host + ":" + strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return NewDataMeshClient(conn), nil
}
