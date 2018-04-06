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

func (serializer *Serializer) Size() int {
	queue := serializer.queue
	return len(queue)
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
			//log.Printf("APPENDED. LEN: %d\n", len(serializer.queue))
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
						//pl := len(serializer.queue)
						serializer.queue = serializer.queue[1:]
						//log.Printf("REMOVED 0. LEN: %d -> %d\n", pl, len(serializer.queue))
					} else {
						//pl := len(serializer.queue)
						newQueue := append(serializer.queue[:idx], serializer.queue[idx+1:]...)
						serializer.queue = newQueue
						//log.Printf("REMOVED 1. LEN: %d -> %d\n", pl, len(serializer.queue))
					}
					serializer.applyPending()
				} else {
					//pl := len(serializer.queue)
					serializer.queue = serializer.queue[0:0]
					//log.Printf("REMOVED 2. LEN: %d -> %d\n", pl, len(serializer.queue))
				}
				return
			}
		}
	}
}