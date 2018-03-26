package flux

import (
	"context"
)

// Utilities for consuming and producing events

type Consumer interface {

	Consume(ctx context.Context, sourceChannel <-chan interface{})

}

type Receiver interface {
	onReceive(element interface{})
}

type BaseConsumer struct {
	internalChannel	chan interface{}
	ctx	context.Context
}

func (b *BaseConsumer) initConsumer(ctx	context.Context, receiver Receiver) {
	b.ctx = ctx
	b.internalChannel = make(chan interface{})

	go func() {
		for {
			select {
			case <- ctx.Done():
				return
			case element, ok := <- b.internalChannel:
				if !ok {
					return
				}
				receiver.onReceive(element)
			}
		}
	}()
}

func (b *BaseConsumer) Consume(ctx context.Context, sourceChannel <-chan interface{}) {

	go func() {
		for {
			select {
				case <- ctx.Done():
					return
				case e, ok := <- sourceChannel:
					if !ok {
						return
					}

					select {
						case <- b.ctx.Done():
							return
						case b.internalChannel <- e:
					}
			}
		}
	}()

}




