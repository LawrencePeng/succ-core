package util

const (
	FBS = 16
	FBS_HIBIT = 4
	NUM_BUCKETS = 32
)

type IntArrayList struct {
	Buckets		[][]int32
	CurrentIdx	int32
}


func NewIntArrayList() *IntArrayList {
	buckets := make([][]int32, NUM_BUCKETS)
	buckets[0] = make([]int32, FBS)
	for i := 1; i < NUM_BUCKETS; i++ {
		buckets[i] = nil
	}
	cur := int32(0)
	return &IntArrayList{
		Buckets:buckets,
		CurrentIdx: cur,
	}
}

func NumOfLeadingZero(i int32) int32 {
	if i == int32(0) {
		return 32
	}

	n := int32(1)
	if uint32(i) >> 16 == 0 {
		n += 16
		i = i << 16
	}
	if uint32(i) >> 24 == 0 {
		n += 8
		i = i << 8
	}
	if uint32(i) >> 28 == 0 {
		n += 4
		i = i << 4
	}
	if uint32(i) >> 30 == 0 {
		n += 2
		i = i << 2
	}

	n += int32(i >> 31)
	return n
}

func (intArrList *IntArrayList) Add(val int32) {
	pos := intArrList.CurrentIdx + FBS
	intArrList.CurrentIdx++
	hibit := 31 - NumOfLeadingZero(pos)
	bucketOff := pos ^ (1 << uint32(hibit))
	bucketIdx := hibit - FBS_HIBIT
	if intArrList.Buckets[bucketIdx] == nil {
		size := 1 << uint32(bucketIdx + FBS_HIBIT)
		intArrList.Buckets[bucketIdx] = make([]int32, size)
	}
	intArrList.Buckets[bucketIdx][bucketOff] = val
}

func (intArrayList *IntArrayList) Get(idx int32) int32 {
	pos := idx + FBS
	hibit := 31 - NumOfLeadingZero(pos)
	bucketOff := pos ^ (1 << uint32(hibit))
	bucketIdx := hibit - FBS_HIBIT
	return intArrayList.Buckets[bucketIdx][bucketOff]
}

func (intArrayList *IntArrayList) Size() int32 {
	return intArrayList.CurrentIdx
}