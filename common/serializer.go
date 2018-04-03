package common

import "log"

type ValueApplier interface {
	ApplyValue(interface{}) (bool, error)
}

type Serializer struct {
	queue	[]interface{}
	ingress	chan interface{}
	event	chan bool
	applier	ValueApplier
}

func NewSerializer(applier ValueApplier) *Serializer {
	serializer := Serializer{
		ingress:	make(chan interface{}, 10),
		event:		make(chan bool),
		applier: 	applier,
	}
	go serializer.runCycle()
	return &serializer
}

func (serializer *Serializer) Push(value interface{}) {
	serializer.ingress <- value
}

func (serializer *Serializer) OnNotification(Notification) error {
	go func() {
		serializer.event <- true
	}()
	return nil
}

func (serializer *Serializer) runCycle() {
	for {
		select {
		case <- serializer.event:
			serializer.applyPending()
		case value := <- serializer.ingress:
			serializer.queue = append(serializer.queue, value)
			serializer.applyPending()
		}
	}
}

func (serializer *Serializer) applyPending() {
	for idx, value := range serializer.queue {
		if ok, err := serializer.applier.ApplyValue(value); ok || err != nil {
			if err != nil {
				log.Println("Error while processing:", err)
			} else {
				if queueLen := len(serializer.queue); queueLen > 1 {
					if idx == 0 {
						serializer.queue = serializer.queue[1:]
					} else {
						newQueue := make([]interface{}, len(serializer.queue))
						newQueue = append(newQueue, serializer.queue[:idx]...)
						if idx + 1 < queueLen {
							newQueue = append(newQueue, serializer.queue[idx+1:]...)
						}
						serializer.queue = newQueue
					}

					serializer.applyPending()
				} else {
					serializer.queue = serializer.queue[0:0]
				}
				return
			}
		}
	}
}