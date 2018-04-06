package controller

import (
	"github.com/nicolaferraro/datamesh/projection"
	"github.com/nicolaferraro/datamesh/log"
	"github.com/nicolaferraro/datamesh/protobuf"
	"errors"
	"github.com/nicolaferraro/datamesh/common"
	logger "log"
)

type Controller struct {
	projection		*projection.Projection
	log				*log.Log
	notifier		*Notifier
	globalVersion	uint64
	serializer		*common.Serializer
}

func NewController(projection *projection.Projection, log *log.Log, notifier *Notifier) *Controller {
	ctrl := Controller{
		projection: projection,
		log: log,
		notifier: notifier,
	}
	ctrl.serializer = common.NewSerializer(&ctrl)
	log.Listen(ctrl.serializer)
	return &ctrl
}

func (ctrl *Controller) Apply(transaction *protobuf.Transaction) error {
	logger.Printf("Received transaction for event %s\n", transaction.Event.Name)
	ctrl.serializer.Push(transaction)
	return nil
}

func (ctrl *Controller) ApplyValue(value interface{}) (bool, error) {
	if transaction, ok := value.(*protobuf.Transaction); ok {
		if transaction == nil || transaction.Event == nil {
			return false, errors.New("Cannot apply empty or incomplete transaction")
		}

		eventVersion := transaction.Event.Version
		cachedEvent := ctrl.log.Cache.Get(transaction.Event.ClientIdentifier)
		//if cachedEvent != nil {
		//	logger.Printf("Mapping client identifier %s to version %d\n", transaction.Event.ClientIdentifier, cachedEvent.Version)
		//}
		if eventVersion == 0 {
			if cachedEvent == nil {
				return false, nil
			}
			eventVersion = cachedEvent.Version
		}

		prjVersion := ctrl.projection.Version
		if eventVersion != prjVersion + 1 {
			return false, nil
		}

		for _, operation := range transaction.Operations {
			if operation.GetRead() != nil {
				path := operation.GetRead().Path

				currentVersion, _, err := ctrl.projection.Get(path.Location)
				if err != nil {
					return false, err
				}

				if path.Version < currentVersion {
					logger.Printf("Cannot apply transaction %d. Data read by transaction has changed from version %d to %d. Discarding.\n", eventVersion, path.Version, currentVersion)
					if cachedEvent != nil {
						ctrl.notifier.Notify(cachedEvent)
					}
					return true, nil
				}
			}
		}

		logger.Printf("Applying transaction %d\n", eventVersion)
		for _, operation := range transaction.Operations {
			if err := ctrl.applyOperation(operation); err != nil {
				ctrl.projection.Rollback()
				return false, err
			}
		}

		//logger.Printf("Serializer size is now %d\n", ctrl.serializer.Size())

		if err := ctrl.projection.Commit(); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (ctrl *Controller) applyOperation(operation *protobuf.Operation) error {
	if operation == nil || operation.GetRead() != nil {
		return nil
	} else if operation.GetUpsert() != nil {
		return ctrl.applyUpsert(operation.GetUpsert())
	} else if operation.GetDelete() != nil {
		return ctrl.applyDelete(operation.GetDelete())
	} else {
		return errors.New("Unsupported operation type")
	}
}

func (ctrl *Controller) applyUpsert(operation *protobuf.UpsertOperation) error {
	unmarshalled, err := operation.Data.Unmarshal()
	if err != nil {
		return err
	}
	expanded, err := unmarshalled.Expand()
	if err != nil {
		return err
	}

	if len(expanded) > 1 {
		if err := ctrl.projection.Delete(operation.Data.Path.Location); err != nil {
			return err
		}
	}
	for _, data := range expanded {
		if err := ctrl.projection.Upsert(data.Path.Location, data.Content); err != nil {
			return err
		}
	}
	return nil
}

func (ctrl *Controller) applyDelete(operation *protobuf.DeleteOperation) error {
	return ctrl.projection.Delete(operation.Path.Location)
}
