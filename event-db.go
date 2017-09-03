package main

import (
	"fmt"
	"github.com/nicolaferraro/event-db/pkg/core"
	"github.com/nicolaferraro/event-db/pkg/api"
	"context"
	"time"
	_"github.com/nicolaferraro/event-db/pkg/flux"
)

func main() {

	/*
	stop := make(chan interface{})
	a := make(chan int)

	go func() {
		defer close(a)
		defer fmt.Println("Closing channel a")
		for {
			fmt.Println("Write iteration")
			select {
			case <- stop:
				fmt.Println("Closing write")
				return
			case a <- rand.Int():
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	go func() {
		defer fmt.Println("Read goroutine closing")
		for {
			fmt.Println("Read iteration")
			select {
			case <- stop:
				fmt.Println("Closing read")
				return
			case v, ok := <- a:
				if (!ok) {
					return
				}
				fmt.Println("------ recv " + fmt.Sprint(v))
			}
		}
	}()

	time.Sleep(1 * time.Second)
	close(stop)

	*/

	startdb()
	time.Sleep(1 * time.Second) // To check deinitialization
}

func startdb() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database := core.NewDatabase(ctx, "Hello World")

	fmt.Println(database.Name())

	go func() {
		database.EventInputChannel() <- core.NewEvent(api.EventAdd)
		database.EventInputChannel() <- core.NewEvent(api.EventRemove)
		database.EventInputChannel() <- core.NewEvent(api.EventNoop)
	}()

	for _ = range []int{1, 2, 3} {
		e := <- database.EventOutputChannel()
		fmt.Println(e.Type())
	}

}
