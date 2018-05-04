package transaction

import (
	"github.com/nicolaferraro/datamesh/context/projection"
	"github.com/nicolaferraro/datamesh/protobuf"
	"errors"
	"github.com/nicolaferraro/datamesh/common"
	"github.com/nicolaferraro/datamesh/notification"
	"github.com/golang/glog"
	"context"
	"time"
)

const (
	MaxTransactionShift = 50
	MaxPreflightBufferSize = 25
	MaxTimeInPreflightBuffer = 10 * time.Second
)

type TransactionManager struct {
	contextId		string
	projection		*projection.Projection
	bus				*notification.NotificationBus
	globalVersion	uint64
	serializer		*common.Serializer
	eventCache		*EventCache
	preflightBuffer	map[string]int64
}

func NewTransactionManager(ctx context.Context, contextId string, projection *projection.Projection, bus *notification.NotificationBus) *TransactionManager {
	tx := TransactionManager{
		contextId: contextId,
		projection: projection,
		bus: bus,
		preflightBuffer: make(map[string]int64),
	}
	tx.serializer = common.NewSerializer(ctx, &tx)
	tx.eventCache = NewEventCache()
	bus.Connect(&tx)
	return &tx
}

func (tx *TransactionManager) OnNotification(n notification.Notification) {
	if n.TransactionReceivedNotification != nil {
		glog.V(4).Infof("Received transaction for event %s in context %s", n.TransactionReceivedNotification.Transaction.Event.Name, tx.contextId)
		tx.serializer.Push(n.TransactionReceivedNotification.Transaction)
	} else if n.EventAppendedNotification != nil {
		tx.eventCache.Register(n.EventAppendedNotification.Event)
		tx.serializer.OnNotification(n) // Forward to unlock serializer
	} else if n.TransactionProcessedNotification != nil {
		if n.TransactionProcessedNotification.Error == nil {
			tx.eventCache.Prune(n.TransactionProcessedNotification.Version)
		}
	}
}

func (tx *TransactionManager) ExecuteSerially(value interface{}) (bool, error) {
	if transaction, ok := value.(*protobuf.Transaction); ok {
		if transaction == nil || transaction.Event == nil {
			return false, errors.New("Cannot apply empty or incomplete transaction in context " + tx.contextId)
		}

		// Delete old snapshots
		tx.cleanupPreflightBuffer()

		eventVersion := transaction.Event.Version
		cachedEvent := tx.eventCache.Get(transaction.Event.ClientIdentifier)

		if eventVersion == 0 {
			if cachedEvent == nil {
				if len(tx.preflightBuffer) > MaxPreflightBufferSize {
					return false, errors.New("Cannot find event in cache and buffer is full in context " + tx.contextId)
				} else {
					tx.preflightBuffer[transaction.Event.ClientIdentifier] = time.Now().Unix()
					glog.V(10).Info("Buffering transaction ", transaction.Event.ClientIdentifier, " in context", tx.contextId)
					return false, nil
				}
			}
			eventVersion = cachedEvent.Version
		}
		delete(tx.preflightBuffer, transaction.Event.ClientIdentifier)

		prjVersion := tx.projection.Version
		if eventVersion <= prjVersion {
			glog.V(1).Infof("Discarding old transaction for version %d in context %s", eventVersion, tx.contextId)
			return true, nil
		} else if eventVersion > prjVersion + MaxTransactionShift {
			glog.V(1).Infof("Discarding new transaction %d to keep the transaction buffer size low in context %s", eventVersion, tx.contextId)
			tx.bus.Notify(notification.NewTransactionFailedNotification(tx.projection, cachedEvent, errors.New("Too much traffic in context " + tx.contextId)))
			return true, nil
		} else if eventVersion != prjVersion + 1 {
			return false, nil // enqueue if not next
		}

		for _, operation := range transaction.Operations {
			if operation.GetRead() != nil {
				path := operation.GetRead().Path

				currentVersion, _, err := tx.projection.Get(path.Location)
				if err != nil {
					return false, err
				}

				if path.Version < currentVersion {
					glog.V(1).Infof("Cannot apply transaction %d in context %s. Data read by transaction has changed from version %d to %d. Discarding.\n", eventVersion, tx.contextId, path.Version, currentVersion)
					tx.bus.Notify(notification.NewTransactionFailedNotification(tx.projection, cachedEvent, errors.New("Transaction conflict in context " + tx.contextId)))
					return true, nil
				}
			}
		}

		glog.V(1).Infof("Applying transaction %d in context %s", eventVersion, tx.contextId)
		for _, operation := range transaction.Operations {
			if err := tx.applyOperation(operation); err != nil {
				tx.projection.Rollback()
				return false, err
			}
		}

		if err := tx.projection.Commit(); err != nil {
			return false, err
		}

		tx.bus.Notify(notification.NewTransactionProcessedNotification(tx.projection, eventVersion))
		return true, nil
	}
	return false, nil
}

func (tx *TransactionManager) cleanupPreflightBuffer() {
	now := time.Now().Unix()
	var oldEvents []string
	for k, v := range tx.preflightBuffer {
		if v < now - int64(MaxTimeInPreflightBuffer) {
			if oldEvents == nil {
				oldEvents = make([]string, 0)
			}
			oldEvents = append(oldEvents, k)
		}
	}
	if oldEvents != nil {
		for _, id := range oldEvents {
			glog.V(4).Infof("Deleting transaction %s from preflight buffer: expired", id)
			delete(tx.preflightBuffer, id)
		}
	}
}

func (tx *TransactionManager) applyOperation(operation *protobuf.Operation) error {
	if operation == nil || operation.GetRead() != nil {
		return nil
	} else if operation.GetUpsert() != nil {
		return tx.applyUpsert(operation.GetUpsert())
	} else if operation.GetDelete() != nil {
		return tx.applyDelete(operation.GetDelete())
	} else {
		return errors.New("Unsupported operation type")
	}
}

func (tx *TransactionManager) applyUpsert(operation *protobuf.UpsertOperation) error {
	unmarshalled, err := operation.Data.Unmarshal()
	if err != nil {
		return err
	}
	expanded, err := unmarshalled.Expand()
	if err != nil {
		return err
	}

	if len(expanded) > 1 {
		if err := tx.projection.Delete(operation.Data.Path.Location); err != nil {
			return err
		}
	}
	for _, data := range expanded {
		if err := tx.projection.Upsert(data.Path.Location, data.Content); err != nil {
			return err
		}
	}
	return nil
}

func (tx *TransactionManager) applyDelete(operation *protobuf.DeleteOperation) error {
	return tx.projection.Delete(operation.Path.Location)
}
