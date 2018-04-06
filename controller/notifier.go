package controller

import (
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/nicolaferraro/datamesh/common"
	"log"
)

type Notifier struct {
	serializer	*common.Serializer
	consumers	[]common.EventConsumer
	last		int
}

func NewNotifier() *Notifier {
	notifier := Notifier{
		last: -1,
	}
	notifier.serializer = common.NewSerializer(&notifier)
	return &notifier
}

func (n *Notifier) Notify(evt *protobuf.Event) {
	n.serializer.Push(evt)
}

func (n *Notifier) ApplyValue(value interface{}) (bool, error) {
	if evt, ok := value.(*protobuf.Event); ok {
		if len(n.consumers) > 0 {
			next := (n.last + 1) % len(n.consumers)

			consumer := n.consumers[next]
			if err := consumer.Consume(evt); err != nil {
				log.Println("Error while push event to the client: ", err)
				log.Println("Removing client from the list of consumers")

				newConsumers := make([]common.EventConsumer, 0)
				if next > 0 {
					newConsumers = append(newConsumers, n.consumers[0:next]...)
				}
				if next < len(n.consumers) - 1 {
					newConsumers = append(newConsumers, n.consumers[next+1:]...)
				}
				n.consumers = newConsumers
			} else {
				n.last = next
			}
		} else {
			log.Println("No consumers available for pushing the event")
		}
		return true, nil
	} else if consumer, ok := value.(common.EventConsumer); ok {
		n.consumers = append(n.consumers, consumer)
		return true, nil
	}
	return false, nil
}

func (n *Notifier) ConnectEventConsumer(consumer common.EventConsumer) {
	n.serializer.Push(consumer)
}



