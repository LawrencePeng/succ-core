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
	Alphabet        []int32
	OriginalSize    int32
	SamplingRateSA  int32
	SamplingRateISA int32
	SamplingRateNPA int32
	SampleBitWidth  int32
}

func (sc *SuccinctCore) FindCharacter(c int32) int32 {
	bg := int32(1)
	end := int32(len(sc.Alphabet)) - 1

	if c < bg || c > end {
		return -1
	}

	for bg <= end {
		mid := (bg + end) / 2
		if sc.Alphabet[mid] == c {
			return mid
		} else if sc.Alphabet[mid] > c {
			end = mid - 1
		} else {
			bg = mid + 1
		}
	}

	return -1
}

func (sc *SuccinctCore) BaseSize() int32 {
	return 6*int32(util.INT_SIZE) + (12 + int32(len(sc.Alphabet)*util.INT_SIZE))
}
