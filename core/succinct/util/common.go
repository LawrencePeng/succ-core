package util

const Two32 = int64(1 << 32)


func IntLog2(n int64) int {
	if n < 0 {
		return -1
	}

	var l int
	if (n != 0) && ((n & (n - 1)) == 0) {
		l = 0
	} else {
		l = 1
	}

	n = n >> 1
	for ; n > 0; {
		l ++
		n = n >> 1
	}

	return l
}


func Mod(a int64, n int64) int64 {
	for ;a < 0; {
		a += n
	}
	return a % n
}

func PopCount(x uint64) int {
	x = (x & 0x5555555555555555) + ((x & 0xAAAAAAAAAAAAAAAA) >> 1)
	x = (x & 0x3333333333333333) + ((x & 0xCCCCCCCCCCCCCCCC) >> 2)
	x = (x & 0x0F0F0F0F0F0F0F0F) + ((x & 0xF0F0F0F0F0F0F0F0) >> 4)
	x *= 0x0101010101010101
	return int((x >> 56) & 0xFF)
}

func NumBlocks(n int, blockSize int) int {
	if n % blockSize == 0 {
		return n / blockSize
	} else {
		return n / blockSize + 1
	}
}