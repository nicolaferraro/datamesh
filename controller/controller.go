package controller

import (
	"github.com/nicolaferraro/datamesh/projection"
	"github.com/nicolaferraro/datamesh/log"
	"github.com/nicolaferraro/datamesh/protobuf"
	"errors"
	"github.com/nicolaferraro/datamesh/common"
)

type Controller struct {
	projection		*projection.Projection
	log				*log.Log
	globalVersion	uint64
	serializer		*common.Serializer
}

func NewController(projection *projection.Projection, log *log.Log) *Controller {
	ctrl := Controller{
		projection: projection,
		log: log,
	}
	ctrl.serializer = common.NewSerializer(&ctrl)
	log.Listen(ctrl.serializer)
	return &ctrl
}

func (ctrl *Controller) Apply(transaction *protobuf.Transaction) error {
	ctrl.serializer.Push(transaction)
	return nil
}

func (ctrl *Controller) ApplyValue(value interface{}) (bool, error) {
	if transaction, ok := value.(*protobuf.Transaction); ok {
		if transaction == nil || transaction.Event == nil {
			return false, errors.New("Cannot apply empty or incomplete transaction")
		}

		eventVersion := transaction.Event.Version
		if eventVersion == 0 {
			event := ctrl.log.Cache.Get(transaction.Event.ClientIdentifier)
			if event == nil {
				return false, nil
			}
			eventVersion = event.Version
		}

		prjVersion := ctrl.projection.Version
		if eventVersion != prjVersion + 1 {
			return false, nil
		}

		for _, operation := range transaction.Operations {
			if err := ctrl.applyOperation(operation); err != nil {
				ctrl.projection.Rollback()
				return false, err
			}
		}

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
