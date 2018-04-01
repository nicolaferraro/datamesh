package service

import (
	"context"
	"google.golang.org/grpc"
	"fmt"
	"net"
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nicolaferraro/datamesh/common"
)


type DefaultDataMeshServer struct {
	port     		int
	consumer		common.EventConsumer
	grpcServer 		*grpc.Server
}

func NewDefaultDataMeshServer(port int, consumer common.EventConsumer) *DefaultDataMeshServer {
	return &DefaultDataMeshServer{
		port: port,
		consumer: consumer,
	}
}

func (srv *DefaultDataMeshServer) Push(ctx context.Context, evt *protobuf.Event) (*empty.Empty, error) {
	if err := srv.consumer.Consume(evt); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}


func (srv *DefaultDataMeshServer) FastProcess(context.Context, *protobuf.Transaction) (*empty.Empty, error) {
	// TBD
	return &empty.Empty{}, nil
}


func (srv *DefaultDataMeshServer) Process(protobuf.DataMesh_ProcessServer) error {
	// TBD
	return nil
}

func (srv *DefaultDataMeshServer) Read(context.Context, *protobuf.Path) (*protobuf.Data, error) {
	// TBD
	return nil, nil
}

func (srv *DefaultDataMeshServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", srv.port))
	if err != nil {
		return err
	}

	srv.grpcServer = grpc.NewServer()
	protobuf.RegisterDataMeshServer(srv.grpcServer, srv)
	srv.grpcServer.Serve(lis)
	return nil
}
