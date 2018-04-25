package processor

import (
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/nicolaferraro/datamesh/common"
	"context"
	"github.com/nicolaferraro/datamesh/notification"
	"github.com/golang/glog"
)

type Communicator struct {
	contextId	string
	serializer	*common.Serializer
	consumers	[]common.EventConsumer
	last		int
}

func NewCommunicator(ctx context.Context, contextId string, bus *notification.NotificationBus) *Communicator {
	comm := Communicator{
		contextId: contextId,
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
		glog.V(4).Infof("Requesting processing of event %d to a client in context %s", evt.Version, c.contextId)
		if len(c.consumers) > 0 {
			glog.V(4).Infof("%d clients available to process the event in context %s", len(c.consumers), c.contextId)
			next := (c.last + 1) % len(c.consumers)

			consumer := c.consumers[next]
			if err := consumer.Consume(evt); err != nil {
				glog.V(1).Info("Error while pushing event to the client in context " + c.contextId + ": ", err)
			} else {
				glog.V(4).Infof("Event pushed back to the client in context %s", c.contextId)
			}
			c.last = next
		} else {
			glog.V(1).Infof("No clients available for pushing the event in context %s", c.contextId)
		}
		return true, nil
	} else if connected, ok := value.(*notification.ClientConnectedNotification); ok {
		glog.V(1).Infof("Remote client added to the list of available clients in context %s", c.contextId)
		c.consumers = append(c.consumers, connected.Client)
		return true, nil
	} else if disconnected, ok := value.(*notification.ClientDisconnectedNotification); ok {
		for idx, cons := range c.consumers {
			if cons == disconnected.Client {
				glog.V(1).Infof("Remote client removed from the list of available clients in context %s", c.contextId)
				c.consumers = append(c.consumers[:idx], c.consumers[idx+1:]...)
				return true, nil
			}
		}

		glog.Errorf("Cannot find remote client in the list of available clients in context %s", c.contextId)
		return true, nil
	}
	return false, nil
}
