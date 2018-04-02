package common

import "github.com/nicolaferraro/datamesh/protobuf"

type TransactionExecutor interface {
	Apply(transaction *protobuf.Transaction) error
}
