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

	database := core.NewDatabase(ctx, "Hello")
	defer close(database.EventInputChannel())

	fmt.Println(database.Name())

	go func() {
		select {
		case database.EventInputChannel() <- core.NewEvent(eventtype.Add):
		case <-ctx.Done():
			return
		}

		select {
		case database.EventInputChannel() <- core.NewEvent(eventtype.Remove):
		case <-ctx.Done():
			return
		}

		select {
		case database.EventInputChannel() <- core.NewEvent(eventtype.Noop):
		case <-ctx.Done():
			return
		}
	}()

	for _ = range []int{1, 2, 3} {
		e := <- database.EventOutputChannel()
		fmt.Println(e.Type())
	}

}
