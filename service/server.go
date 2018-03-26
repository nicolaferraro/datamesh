package service

import (
	"context"
	"google.golang.org/grpc"
	"fmt"
	"net"
	"log"
	"github.com/nicolaferraro/datamesh/utils"
	"github.com/golang/protobuf/proto"
)

/*
type DataMeshServer interface {
	Push(context.Context, *Event) (*Empty, error)
}
*/


type DefaultDataMeshServer struct {
	port		int
	observer	utils.MessageObserver
}

func NewDefaultDataMeshServer(port int, observer utils.MessageObserver) *DefaultDataMeshServer {
	return &DefaultDataMeshServer{
		port: port,
		observer: observer,
	}
}

func (srv *DefaultDataMeshServer) Push(ctx context.Context, evt *Event) (*Empty, error) {
	msg, err := proto.Marshal(evt)
	if err != nil {
		return nil, err
	}
	return nil, srv.observer.Accept(msg)
}

func (srv *DefaultDataMeshServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", srv.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	RegisterDataMeshServer(grpcServer, srv)
	return grpcServer.Serve(lis)
}
