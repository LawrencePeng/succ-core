package util

type Source interface {
	Len() int
	Get(i int) int
}

