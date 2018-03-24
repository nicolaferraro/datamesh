package core

import (
	"github.com/nicolaferraro/event-db/tests/api"
	"context"
	"fmt"
)

// Default database struct
type DefaultDatabase struct {
	name	string
	eventInputChannel	chan api.Event
	eventOutputChannel	chan api.Event
}

// Constructor
func NewDatabase(ctx context.Context, name string) api.Database {
	db := DefaultDatabase{
		name: name,
		eventInputChannel: make(chan api.Event),
		eventOutputChannel: make(chan api.Event),
	}

	// Set database deinitialization
	go func() {
		<- ctx.Done()
		fmt.Println("Closing database channels")
		close(db.eventInputChannel)
		close(db.eventOutputChannel)
	}()

	go func() {
		for {
			select {
			case event := <- db.eventInputChannel:
				select {
				case db.eventOutputChannel <- event:
				case <- ctx.Done():
					return
				}
			case <- ctx.Done():
				return
			}
		}
	}()

	return db
}

// Interface methods

func (database DefaultDatabase) Name() string {
	return database.name
}


func (database DefaultDatabase) EventInputChannel()		chan<- api.Event {
	return database.eventInputChannel
}

func (database DefaultDatabase) EventOutputChannel()	<-chan api.Event {
	return database.eventOutputChannel
}