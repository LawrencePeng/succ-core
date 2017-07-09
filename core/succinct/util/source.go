package util

type Source interface {
	Len() int32
	Get(i int32) int32
}

