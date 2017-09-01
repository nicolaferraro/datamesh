package main

import (
	"fmt"
	"github.com/nicolaferraro/event-db/pkg/core"
	"github.com/nicolaferraro/event-db/pkg/api/eventtype"
	"context"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database := core.NewDatabase(ctx, "Hello World")
	defer close(database.EventInputChannel())

	fmt.Println(database.Name())

	go func() {
		database.EventInputChannel() <- core.NewEvent(eventtype.Add)
		database.EventInputChannel() <- core.NewEvent(eventtype.Remove)
		database.EventInputChannel() <- core.NewEvent(eventtype.Noop)
	}()

	for _ = range []int{1, 2, 3} {
		e := <- database.EventOutputChannel()
		fmt.Println(e.Type())
	}

}
