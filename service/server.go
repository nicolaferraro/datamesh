package service

import (
	"context"
	"google.golang.org/grpc"
	"fmt"
	"net"
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/nicolaferraro/datamesh/common"
	"encoding/json"
	"github.com/nicolaferraro/datamesh/notification"
	"github.com/golang/glog"
	meshcontext "github.com/nicolaferraro/datamesh/context"
	"time"
	"errors"
)

const (
	PingMaxDelay = 150 * time.Second
)

type DefaultDataMeshServer struct {
	port     			int
	consumer			common.EventConsumer
	mesh				*Mesh
	grpcServer 			*grpc.Server
}

func NewDefaultDataMeshServer(port int, consumer common.EventConsumer, msh *Mesh) *DefaultDataMeshServer {
	return &DefaultDataMeshServer{
		port: port,
		consumer: consumer,
		mesh: msh,
	}
}

func (srv *DefaultDataMeshServer) Push(ctx context.Context, evt *protobuf.Event) (*protobuf.Empty, error) {
	if err := srv.consumer.Consume(evt); err != nil {
		return nil, err
	}
	return &protobuf.Empty{}, nil
}

func (srv *DefaultDataMeshServer) Process(ctx context.Context, transaction *protobuf.Transaction) (*protobuf.Empty, error) {
	glog.V(10).Infof("Received transaction with version %d\n", transaction.Event.Version)
	meshContext := srv.mesh.GetContext(transaction.Context.Name, transaction.Context.Revision)
	meshContext.Notify(notification.NewTransactionReceivedNotification(transaction))
	return &protobuf.Empty{}, nil
}

func (srv *DefaultDataMeshServer) Connect(server protobuf.DataMesh_ConnectServer) error {
	glog.V(1).Info("Processing client connected")
	consumer := protobuf.NewProcessQueueConsumer(server)

	disconnect := make(chan bool, 1)
	ping := make(chan bool, 1)
	var meshContextRef *meshcontext.MeshContext
	go func() {
		contextReceived := false
		for {
			status, err := server.Recv()
			if err != nil || status.GetDisconnect() != nil {
				disconnect <- true
				return
			} else if status.GetConnect() != nil && !contextReceived {
				contextReceived = true
				glog.V(1).Infof("Processing client using context %s with revision %d", status.GetConnect().Name, status.GetConnect().Revision)
				meshContextRef = srv.mesh.GetContext(status.GetConnect().Name, status.GetConnect().Revision)
				meshContextRef.Notify(notification.NewClientConnectedNotification(consumer))
			} else if status.GetPing() != nil {
				select {
				case ping <- true:
				default:
				}
			}
		}
	}()

	WhileActive:
		for {
			select {
			case <-consumer.Closed:
				glog.V(1).Info("Processing client disconnected by server")
				break WhileActive;
			case <-server.Context().Done():
				glog.V(1).Info("Processing client disconnected (gone)")
				break WhileActive;
			case <-disconnect:
				glog.V(1).Info("Processing client sent a disconnect message")
				break WhileActive;
			case <-ping:
				glog.V(8).Info("Processing client sent a ping")
			case <-time.After(PingMaxDelay):
				glog.V(8).Infof("Ping not received within %d seconds: disconnecting client", PingMaxDelay / time.Second)
				break WhileActive;
			}
		}

	if meshContextRef != nil {
		meshContextRef.Notify(notification.NewClientDisconnectedNotification(consumer))
	}
	return errors.New("Disconnected")
}

func (srv *DefaultDataMeshServer) Read(ctx context.Context, req *protobuf.ReadRequest) (*protobuf.Data, error) {
	meshContext := srv.mesh.GetContext(req.Context.Name, req.Context.Revision)
	version, data, err := meshContext.GetDataRetriever().Get(req.Path.Location)
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
			Location: req.Path.Location,
		},
		Content: jsonData,
	}, nil
}

func (srv *DefaultDataMeshServer) Health(ctx context.Context, context *protobuf.Context) (*protobuf.Readiness, error) {
	meshContext := srv.mesh.GetContext(context.Name, context.Revision)

	readiness := &protobuf.Readiness{
		Context: context,
		Ready: meshContext.Status() == meshcontext.Ready,
	}

	return readiness, nil
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
