package util

func GETRANKL2(n int64) int64 {
	return int64(uint64(n) >> 32)
}

func GETRANKL1(n int64, i int) int64 {
	return int64(uint64(n&0x7fffffff)>>(32-uint64(i)*10)) & 0x3ff
}

func GETPOSL2(n int64) int64 {
	return int64(uint64(n) >> 31)
}

func GETPOSL1(n int64, i int) int64 {
	return int64(uint64(n&0x7fffffff)>>(31-uint64(i)*10)) & 0x3ff
}
