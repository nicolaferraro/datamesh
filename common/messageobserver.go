package common

type MessageObserver interface {

	Accept([]byte) 	error

}