package service

import (
	"context"
	"google.golang.org/grpc"
	"fmt"
	"net"
	"github.com/nicolaferraro/datamesh/common"
	"github.com/golang/protobuf/proto"
)

/*
type DataMeshServer interface {
	Push(context.Context, *Event) (*Empty, error)
}
*/


type DefaultDataMeshServer struct {
	port     		int
	observer		common.MessageObserver
	grpcServer 		*grpc.Server
}

func NewDefaultDataMeshServer(port int, observer common.MessageObserver) *DefaultDataMeshServer {
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
	if err = srv.observer.Accept(msg); err != nil {
		return nil, err
	}
	return &Empty{}, nil
}

func (srv *DefaultDataMeshServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", srv.port))
	if err != nil {
		return err
	}

	srv.grpcServer = grpc.NewServer()
	RegisterDataMeshServer(srv.grpcServer, srv)
	srv.grpcServer.Serve(lis)
	return nil
}
