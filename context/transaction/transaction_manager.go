package transaction

import (
	"github.com/nicolaferraro/datamesh/context/projection"
	"github.com/nicolaferraro/datamesh/protobuf"
	"errors"
	"github.com/nicolaferraro/datamesh/common"
	logger "log"
	"github.com/nicolaferraro/datamesh/notification"
	"context"
)

type TransactionManager struct {
	projection		*projection.Projection
	bus				*notification.NotificationBus
	globalVersion	uint64
	serializer		*common.Serializer
	eventCache		*EventCache
}

func NewTransactionManager(ctx context.Context, projection *projection.Projection, bus *notification.NotificationBus) *TransactionManager {
	tx := TransactionManager{
		projection: projection,
		bus: bus,
	}
	tx.serializer = common.NewSerializer(ctx, &tx)
	tx.eventCache = NewEventCache()
	bus.Connect(&tx)
	return &tx
}

func (tx *TransactionManager) OnNotification(n notification.Notification) {
	if n.TransactionReceivedNotification != nil {
		logger.Printf("Received transaction for event %s\n", n.TransactionReceivedNotification.Transaction.Event.Name)
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
			return false, errors.New("Cannot apply empty or incomplete transaction")
		}

		eventVersion := transaction.Event.Version
		cachedEvent := tx.eventCache.Get(transaction.Event.ClientIdentifier)

		if eventVersion == 0 {
			if cachedEvent == nil {
				return false, nil
			}
			eventVersion = cachedEvent.Version
		}

		prjVersion := tx.projection.Version
		if eventVersion != prjVersion + 1 {
			return false, nil
		}

		for _, operation := range transaction.Operations {
			if operation.GetRead() != nil {
				path := operation.GetRead().Path

				currentVersion, _, err := tx.projection.Get(path.Location)
				if err != nil {
					return false, err
				}

				if path.Version < currentVersion {
					logger.Printf("Cannot apply transaction %d. Data read by transaction has changed from version %d to %d. Discarding.\n", eventVersion, path.Version, currentVersion)
					tx.bus.Notify(notification.NewTransactionFailedNotification(tx.projection, cachedEvent, errors.New("Transaction conflict")))
					return true, nil
				}
			}
		}

		logger.Printf("Applying transaction %d\n", eventVersion)
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
