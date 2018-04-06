package controller

import (
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/nicolaferraro/datamesh/common"
	"log"
)

type Notifier struct {
	serializer	*common.Serializer
	consumers	[]common.CloseableEventConsumer
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
		log.Printf("Requesting processing of event %d\n", evt.Version)
		if len(n.consumers) > 0 {
			next := (n.last + 1) % len(n.consumers)

			consumer := n.consumers[next]
			if consumer.IsClosed() {
				log.Println("Removing closed consumer")
				newConsumers := append(n.consumers[0:next], n.consumers[next+1:]...)
				n.consumers = newConsumers
				return n.ApplyValue(value)
			} else if err := consumer.Consume(evt); err != nil {
				log.Println("Error while push event to the client: ", err)
				log.Println("Removing client from the list of consumers")

				newConsumers := append(n.consumers[0:next], n.consumers[next+1:]...)
				n.consumers = newConsumers
				consumer.Close()
				return n.ApplyValue(value)
			} else {
				n.last = next
				log.Println("Event pushed back to the client")
			}
		} else {
			log.Println("No consumers available for pushing the event")
		}
		return true, nil
	} else if consumer, ok := value.(common.CloseableEventConsumer); ok {
		n.consumers = append(n.consumers, consumer)
		return true, nil
	}
	return false, nil
}

func (n *Notifier) ConnectEventConsumer(consumer common.CloseableEventConsumer) {
	n.serializer.Push(consumer)
}



