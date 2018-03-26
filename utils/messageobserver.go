package utils

type MessageObserver interface {

	Accept([]byte) 	error

}