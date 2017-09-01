package api

// The database interface
type Database interface {
	Named
	EventReceiver
	EventProducer

}