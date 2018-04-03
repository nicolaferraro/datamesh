package service

import (
	"context"
	"google.golang.org/grpc"
	"fmt"
	"net"
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nicolaferraro/datamesh/common"
	"encoding/json"
)


type DefaultDataMeshServer struct {
	port     		int
	consumer		common.EventConsumer
	executor		common.TransactionExecutor
	retriever		common.DataRetriever
	grpcServer 		*grpc.Server
}

func NewDefaultDataMeshServer(port int, consumer common.EventConsumer, executor common.TransactionExecutor, retriever common.DataRetriever) *DefaultDataMeshServer {
	return &DefaultDataMeshServer{
		port: port,
		consumer: consumer,
		executor: executor,
		retriever: retriever,
	}
}

func (srv *DefaultDataMeshServer) Push(ctx context.Context, evt *protobuf.Event) (*empty.Empty, error) {
	if err := srv.consumer.Consume(evt); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}


func (srv *DefaultDataMeshServer) FastProcess(ctx context.Context, transaction *protobuf.Transaction) (*empty.Empty, error) {
	if err := srv.executor.Apply(transaction); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}


func (srv *DefaultDataMeshServer) Process(protobuf.DataMesh_ProcessServer) error {
	// TBD
	return nil
}

func (srv *DefaultDataMeshServer) Read(ctx context.Context, path *protobuf.Path) (*protobuf.Data, error) {
	data, err := srv.retriever.Get(path.Location)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &protobuf.Data{
		Path: path,
		Content: jsonData,
	}, nil
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
