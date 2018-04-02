package common

type DataRetriever interface {
	Get(key string) (interface{}, error)
}