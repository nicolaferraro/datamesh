package processor

import (
	"context"
	"github.com/nicolaferraro/datamesh/notification"
	"github.com/nicolaferraro/datamesh/common"
	"time"
	"github.com/nicolaferraro/datamesh/protobuf"
)

const (
	DefaultRetransmissionTimeout 	= 3 * time.Second
	DefaultRetransmissionGraceTime 	= 1 * time.Second
)

type EventProcessor struct {
	communicator		*Communicator
	serializer			*common.Serializer
	projectionVersion	uint64
	logVersion			uint64
	states				map[uint64]eventInfo
}

type eventInfo struct {
	Event				*protobuf.Event
	State				eventState
	NumTransmissions	int
}

type eventState int
const (
	eventStateRequested		= eventState(iota)
	eventStateQueued
	eventStateGracePeriod
)

type timeoutInfo struct {
	Event				*protobuf.Event
	Grace				bool
	RefNumTransmissions	int
}

type nextVersionInfo struct {
	Version	uint64
}

func NewEventProcessor(ctx context.Context, bus *notification.NotificationBus) *EventProcessor {
	proc := EventProcessor{}
	proc.communicator = NewCommunicator(ctx, bus)
	proc.serializer = common.NewSerializer(ctx, &proc)
	proc.states = make(map[uint64]eventInfo)

	bus.Connect(&proc)
	return &proc
}

func (proc *EventProcessor) OnNotification(n notification.Notification) {
	proc.serializer.Push(n)
}

func (proc *EventProcessor) ExecuteSerially(value interface{}) (bool, error) {
	if n, ok := value.(notification.Notification); ok {
		if n.EventAppendedNotification != nil {
			// Event is recorded into event log or replayed from it
			ver := n.EventAppendedNotification.Event.Version
			_, hasState := proc.states[ver]
			if ver <= proc.logVersion || hasState {
				return true, nil // discard, should not happen
			} else if ver == proc.logVersion + 1 {
				proc.logVersion = ver
			}

			if n.EventAppendedNotification.Replay {
				proc.sendWithTimeout(n.EventAppendedNotification.Event, 1)
			}

			proc.states[ver] = eventInfo{
				Event: n.EventAppendedNotification.Event,
				State: eventStateRequested,
				NumTransmissions: 1,
			}

		} else if n.TransactionProcessedNotification != nil {
			// Transaction is applied into the projection correctly
			lastVer := proc.projectionVersion
			if n.TransactionProcessedNotification.Version > proc.projectionVersion {
				proc.projectionVersion = n.TransactionProcessedNotification.Version
				for v := lastVer; v<proc.projectionVersion; v++ {
					delete(proc.states, v)
				}
			}

			if _, hasNextState := proc.states[proc.projectionVersion + 1]; hasNextState {
				proc.serializer.Push(nextVersionInfo{
					Version: proc.projectionVersion + 1,
				})
			}
		} else if n.TransactionFailedNotification != nil {
			// Transaction manager has tried to do the transaction and failed
			// The turn is correct but it should be recomputed (immediately)
			ver := n.TransactionFailedNotification.Event.Version
			state, hasState := proc.states[ver]
			if ver > proc.projectionVersion && hasState {
				proc.sendWithTimeout(n.TransactionFailedNotification.Event, state.NumTransmissions + 1)
				proc.states[ver] = eventInfo{
					Event: n.TransactionFailedNotification.Event,
					State: eventStateRequested,
					NumTransmissions: state.NumTransmissions + 1,
				}
			}
		}

		return true, nil
	} else if timeout, ok := value.(timeoutInfo); ok {
		ver := timeout.Event.Version
		state, hasState := proc.states[ver]
		if ver > proc.projectionVersion && hasState && timeout.RefNumTransmissions == state.NumTransmissions {
			// valid timeout

			if !timeout.Grace && state.State == eventStateRequested {
				// First timeout
				next := (ver == proc.projectionVersion + 1)
				if next {
					proc.states[ver] = eventInfo{
						Event: timeout.Event,
						State: eventStateGracePeriod,
						NumTransmissions: state.NumTransmissions,
					}
					timer := time.NewTimer(DefaultRetransmissionGraceTime)
					go func() {
						<- timer.C
						proc.serializer.Push(timeoutInfo{
							Event: timeout.Event,
							Grace: true,
							RefNumTransmissions: state.NumTransmissions,
						})
					}()
				} else {
					proc.states[ver] = eventInfo{
						Event: timeout.Event,
						State: eventStateQueued,
						NumTransmissions: state.NumTransmissions,
					}
				}

			} else if state.State == eventStateGracePeriod {
				// Still not processed
				proc.sendWithTimeout(state.Event, state.NumTransmissions + 1)
				proc.states[ver] = eventInfo{
					Event: timeout.Event,
					State: eventStateRequested,
					NumTransmissions: state.NumTransmissions + 1,
				}
			}
		}
	} else if next, ok := value.(nextVersionInfo); ok {
		state, hasState := proc.states[next.Version]
		if next.Version > proc.projectionVersion && hasState && state.State == eventStateQueued {
			// valid next
			proc.states[next.Version] = eventInfo{
				Event: state.Event,
				State: eventStateGracePeriod,
				NumTransmissions: state.NumTransmissions,
			}
			timer := time.NewTimer(DefaultRetransmissionGraceTime)
			go func() {
				<- timer.C
				proc.serializer.Push(timeoutInfo{
					Event: state.Event,
					Grace: true,
					RefNumTransmissions: state.NumTransmissions,
				})
			}()
		}
	}
	return false, nil
}


func (proc *EventProcessor) sendWithTimeout(evt *protobuf.Event, numTransmissions int) {
	proc.communicator.Send(evt)
	timer := time.NewTimer(DefaultRetransmissionTimeout)
	go func() {
		<- timer.C
		proc.serializer.Push(timeoutInfo{
			Event: evt,
			Grace: false,
			RefNumTransmissions: numTransmissions,
		})
	}()
}
