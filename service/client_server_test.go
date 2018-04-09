package service

import (
	"testing"
	"golang.org/x/net/context"
	"reflect"
	"github.com/stretchr/testify/assert"
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/nicolaferraro/datamesh/common"
	"github.com/nicolaferraro/datamesh/notification"
)

const (
	testDefaultServerPort = 6543
)

type TestStub struct {
	messages	[]*protobuf.Event
}

func (r *TestStub) Consume(evt *protobuf.Event) error  {
	r.messages = append(r.messages, evt)
	return nil
}

func (r *TestStub) ConnectEventConsumer(consumer common.EventConsumer) {
}

func (r *TestStub) Apply(transaction *protobuf.Transaction) error {
	return nil
}

func (r *TestStub) Get(key string) (uint64, interface{}, error) {
	return 0, nil, nil
}

func TestDataMeshClientServer(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	bus := notification.NewNotificationBus(ctx)
	testReceiver := TestStub{}
	server := NewDefaultDataMeshServer(testDefaultServerPort, bus, &testReceiver, &testReceiver)
	go server.Start()

	client, err := NewDataMeshClientConnection("localhost", testDefaultServerPort);
	if err != nil {
		t.Fatal("Cannot create client", err)
	}

	const num = 5
	var evts []protobuf.Event
	for i:=0; i<num; i++ {
		evt := protobuf.Event{
			Name: "evt",
			Payload: []byte{byte(i) + 1,byte(i) + 2, byte(i) + 3},
		}
		evts = append(evts, evt)
		client.Push(ctx, &evt)
	}

	for i:=0; i<num; i++ {
		assert.True(t, reflect.DeepEqual(evts[i], *testReceiver.messages[i]))
	}
}
