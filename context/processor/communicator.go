package processor

import (
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/nicolaferraro/datamesh/common"
	"context"
	"github.com/nicolaferraro/datamesh/notification"
	"github.com/golang/glog"
)

type Communicator struct {
	serializer	*common.Serializer
	consumers	[]common.EventConsumer
	last		int
}

func NewCommunicator(ctx context.Context, bus *notification.NotificationBus) *Communicator {
	comm := Communicator{
		last: -1,
	}
	comm.serializer = common.NewSerializer(ctx, &comm)
	bus.Connect(&comm)
	return &comm;
}

func (c *Communicator) OnNotification(n notification.Notification) {
	if n.ClientConnectedNotification != nil {
		c.serializer.Push(n.ClientConnectedNotification)
	} else if n.ClientDisconnectedNotification != nil {
		c.serializer.Push(n.ClientDisconnectedNotification)
	}
}

func (c *Communicator) Send(evt *protobuf.Event) {
	c.serializer.Push(evt)
}

func (c *Communicator) ExecuteSerially(value interface{}) (bool, error) {
	if evt, ok := value.(*protobuf.Event); ok {
		glog.V(4).Infof("Requesting processing of event %d\n", evt.Version)
		if len(c.consumers) > 0 {
			next := (c.last + 1) % len(c.consumers)

			consumer := c.consumers[next]
			if err := consumer.Consume(evt); err != nil {
				glog.V(4).Info("Error while pushing event to the client: ", err)
			} else {
				glog.V(4).Info("Event pushed back to the client")
			}
			c.last = next
		} else {
			glog.V(4).Info("No consumers available for pushing the event")
		}
		return true, nil
	} else if connected, ok := value.(*notification.ClientConnectedNotification); ok {
		glog.V(1).Info("Remote client added to the list of available clients")
		c.consumers = append(c.consumers, connected.Client)
		return true, nil
	} else if disconnected, ok := value.(*notification.ClientDisconnectedNotification); ok {
		for idx, cons := range c.consumers {
			if cons == disconnected.Client {
				glog.V(1).Info("Remote client removed from the list of available clients")
				c.consumers = append(c.consumers[:idx], c.consumers[idx+1:]...)
				return true, nil
			}
		}

		glog.Error("Cannot find remote client in the list of available clients")
		return true, nil
	}
	return false, nil
}
