package service

import (
	"context"
	"google.golang.org/grpc"
	"fmt"
	"net"
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/nicolaferraro/datamesh/common"
	"encoding/json"
	"log"
	"github.com/nicolaferraro/datamesh/notification"
)


type DefaultDataMeshServer struct {
	port     			int
	consumer			common.EventConsumer
	bus					*notification.NotificationBus
	retriever			common.DataRetriever
	grpcServer 			*grpc.Server
}

func NewDefaultDataMeshServer(port int, bus *notification.NotificationBus, consumer common.EventConsumer, retriever common.DataRetriever) *DefaultDataMeshServer {
	return &DefaultDataMeshServer{
		port: port,
		consumer: consumer,
		bus: bus,
		retriever: retriever,
	}
}

func (srv *DefaultDataMeshServer) Push(ctx context.Context, evt *protobuf.Event) (*protobuf.Empty, error) {
	if err := srv.consumer.Consume(evt); err != nil {
		return nil, err
	}
	return &protobuf.Empty{}, nil
}


func (srv *DefaultDataMeshServer) Process(ctx context.Context, transaction *protobuf.Transaction) (*protobuf.Empty, error) {
	srv.bus.Notify(notification.NewTransactionReceivedNotification(transaction))
	return &protobuf.Empty{}, nil
}


func (srv *DefaultDataMeshServer) ProcessQueue(empty *protobuf.Empty, server protobuf.DataMesh_ProcessQueueServer) error {
	log.Println("Processing client connected")
	consumer := protobuf.NewProcessQueueConsumer(server)
	srv.bus.Notify(notification.NewClientConnectedNotification(consumer))

	select {
		case <- consumer.Closed:
			log.Println("Processing client disconnected by server")
		case <- server.Context().Done():
			log.Println("Processing client disconnected (gone)")
	}

	srv.bus.Notify(notification.NewClientDisconnectedNotification(consumer))
	return nil
}

func (srv *DefaultDataMeshServer) Read(ctx context.Context, path *protobuf.Path) (*protobuf.Data, error) {
	version, data, err := srv.retriever.Get(path.Location)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &protobuf.Data{
		Path: &protobuf.Path{
			Version: version,
			Location: path.Location,
		},
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
