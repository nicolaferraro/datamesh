package flux

import (
	"fmt"
	"testing"
	"context"
)




func TestBaseConsumer_Consume(t *testing.T) {
	globalCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := ConcreteConsumer{}
	a.initConsumer(globalCtx, &a)

	c := make(chan interface{})
	a.Consume(globalCtx, c)

	c <- "Hello"
	c <- "World"
}


// Define common types

type ConcreteConsumer struct {
	BaseConsumer
}

func (c ConcreteConsumer) onReceive(element interface{}) {
	fmt.Println(c, element)
}