package core

import (
	"sort"
	"unsafe"
)

type Succinct interface {
	FindCharacter(int) int
	GetCoreSize() int
	LookupSA(int64) int64
	LookupNPA(int64) int64
	LookupSPA(int64) int64
	LookupIPA(int64) int64
	LookupC(int64 int64) int
	BinSearchNPA(int64, int64, int64, bool) int64

}

type SuccinctCore struct {
	Alphabet         []int
	OriginalSize     int
	SamplingRateSA   int
	SamplingRateISA  int
	SamplingRateNSA  int
	SamplingBitWidth int
}

func (sc *SuccinctCore) FindCharacter(c int) int {
	return sort.SearchInts(sc.Alphabet[1:], c)
}

var INT_SIZE = int(unsafe.Sizeof(int(0)))

func (sc *SuccinctCore) BaseSize() int {
	return 6*INT_SIZE + (12 + len(sc.Alphabet)*INT_SIZE)
}