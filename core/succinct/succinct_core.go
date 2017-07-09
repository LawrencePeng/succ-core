package succinct

import "./util"

type Succinct interface {
	FindCharacter(int) int
	GetCoreSize() int32
	LookupSA(int64) int64
	LookupNPA(int64) int64
	LookupSPA(int64) int64
	LookupIPA(int64) int64
	LookupC(int64 int64) int32
	BinSearchNPA(int64, int64, int64, bool) int64

}

type SuccinctCore struct {
	Alphabet         []int32
	OriginalSize     int32
	SamplingRateSA   int32
	SamplingRateISA  int32
	SamplingRateNSA  int32
	SamplingBitWidth int32
}

func (sc *SuccinctCore) BaseSize() int {
	return 6 * util.INT_SIZE + (12 + len(sc.Alphabet)*util.INT_SIZE)
}