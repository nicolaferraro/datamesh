package common

import (
	"log"
	"context"
)

type SerialExecutor interface {
	ExecuteSerially(interface{}) (bool, error)
}

type Serializer struct {
	queue		[]interface{}
	ingress		chan interface{}
	event		chan bool
	executor	SerialExecutor
}

func NewSerializer(ctx context.Context, executor SerialExecutor) *Serializer {
	serializer := Serializer{
		ingress:	make(chan interface{}, 10),
		event:		make(chan bool),
		executor: 	executor,
	}
	go serializer.runCycle(ctx)
	return &serializer
}

func (serializer *Serializer) Size() int {
	queue := serializer.queue
	return len(queue)
}

func (serializer *Serializer) Push(value interface{}) {
	serializer.ingress <- value
}

func (serializer *Serializer) OnNotification(n interface{}) {
	go func() {
		serializer.event <- true
	}()
}

func (serializer *Serializer) runCycle(ctx context.Context) {
	for {
		select {
		case <- ctx.Done():
			return
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
		if ok, err := serializer.executor.ExecuteSerially(value); ok || err != nil {
			if err != nil {
				log.Println("Error while processing:", err)
			} else {
				if queueLen := len(serializer.queue); queueLen > 1 {
					if idx == 0 {
						serializer.queue = serializer.queue[1:]
					} else {
						var newQueue []interface{}
						if idx == 0 {
							newQueue = serializer.queue[1:]
						} else {
							newQueue = append(serializer.queue[:idx], serializer.queue[idx+1:]...)
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