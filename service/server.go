package service

import (
	"context"
	"google.golang.org/grpc"
	"fmt"
	"net"
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/nicolaferraro/datamesh/common"
	"encoding/json"
)


type DefaultDataMeshServer struct {
	port     			int
	consumer			common.EventConsumer
	executor			common.TransactionExecutor
	retriever			common.DataRetriever
	consumerController	common.EventConsumerController
	grpcServer 			*grpc.Server
}

func NewDefaultDataMeshServer(port int, consumer common.EventConsumer, executor common.TransactionExecutor, retriever common.DataRetriever, consumerController common.EventConsumerController) *DefaultDataMeshServer {
	return &DefaultDataMeshServer{
		port: port,
		consumer: consumer,
		executor: executor,
		retriever: retriever,
		consumerController: consumerController,
	}
}

func (srv *DefaultDataMeshServer) Push(ctx context.Context, evt *protobuf.Event) (*protobuf.Empty, error) {
	if err := srv.consumer.Consume(evt); err != nil {
		return nil, err
	}
	return &protobuf.Empty{}, nil
}


func (srv *DefaultDataMeshServer) Process(ctx context.Context, transaction *protobuf.Transaction) (*protobuf.Empty, error) {
	if err := srv.executor.Apply(transaction); err != nil {
		return nil, err
	}
	return &protobuf.Empty{}, nil
}


func (srv *DefaultDataMeshServer) ProcessQueue(empty *protobuf.Empty, server protobuf.DataMesh_ProcessQueueServer) error {
	srv.consumerController.ConnectEventConsumer(protobuf.NewProcessQueueConsumer(server))
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
