package common

type DataRetriever interface {
	Get(key string) (uint64, interface{}, error)
}