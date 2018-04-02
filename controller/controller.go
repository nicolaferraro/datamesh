package controller

import (
	"github.com/nicolaferraro/datamesh/projection"
	"github.com/nicolaferraro/datamesh/log"
	"github.com/nicolaferraro/datamesh/protobuf"
	"errors"
)

type Controller struct {
	projection		*projection.Projection
	log				*log.Log
	globalVersion	uint64
}

func NewController(projection *projection.Projection, log *log.Log) *Controller {
	return &Controller{
		projection: projection,
		log: log,
	}
}

func (ctrl *Controller) Apply(transaction *protobuf.Transaction) error {
	if transaction == nil || transaction.Event == nil {
		return errors.New("Cannot apply empty or incomplete transaction")
	}

	event := ctrl.log.Cache.Get(transaction.Event.ClientIdentifier)
	if event == nil {
		return errors.New("Cannot find the event that triggered transaction in cache")
	}

	//tvers := event.Version
	for _, operation := range transaction.Operations {
		if err := ctrl.applyOperation(operation); err != nil {
			ctrl.projection.Rollback()
			return err
		}
	}

	return ctrl.projection.Commit()
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
	return ctrl.projection.Upsert(operation.Data.Path.Path, operation.Data.Content.Fields)
}

func (ctrl *Controller) applyDelete(operation *protobuf.DeleteOperation) error {
	return ctrl.projection.Delete(operation.Path.Path)
}
