package service

import (
	"testing"
	"golang.org/x/net/context"
	"github.com/golang/protobuf/proto"
	"reflect"
	"github.com/stretchr/testify/assert"
)

const (
	testDefaultServerPort = 6543
)

type TestReceiver struct {
	messages	[][]byte
}

func (r *TestReceiver) Accept(m []byte) error  {
	r.messages = append(r.messages, m)
	return nil
}

func TestDataMeshClientServer(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	testReceiver := TestReceiver{}
	server := NewDefaultDataMeshServer(testDefaultServerPort, &testReceiver)
	go server.Start()

	client, err := NewDataMeshClientConnection("localhost", testDefaultServerPort);
	if err != nil {
		t.Fatal("Cannot create client", err)
	}

	const num = 5
	var evts []Event
	for i:=0; i<num; i++ {
		evt := Event{
			Name: "evt",
			Payload: []byte{byte(i) + 1,byte(i) + 2, byte(i) + 3},
		}
		evts = append(evts, evt)
		client.Push(ctx, &evt)
	}

	for i:=0; i<num; i++ {
		msg := testReceiver.messages[i]

		var evtr Event
		if err = proto.Unmarshal(msg, &evtr); err != nil {
			t.Fatal("Error while unmarshaling", err)
		}

		assert.True(t, reflect.DeepEqual(evts[i], evtr))
	}
}
