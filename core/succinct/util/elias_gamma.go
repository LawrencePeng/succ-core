package util

var PrefixSum [65536]int32

func init() {
	for i := 0; i < 65536; i++ {
		count := int32(0)
		offset := int32(0)
		sum := int32(0)

		for ;i != 0 && offset <= 16; {
			N := int32(0)

			for ;GetBit(int64(i), offset % 64) != 1; {
				N++
				offset++
			}
			offset++
			if offset + N <= 16 {
				n := int32(uint32(i) >> uint(offset)) &
					(int32(LOW_BITS_SET[N]))
				sum += n + 1 << uint(N)
				offset += N
				count++
			} else {
				offset -= N + 1
				break
			}
		}
		PrefixSum[i] = int32((offset << 24) | (count << 16) | sum)
	}
}

func EncodingSize(value int32) int32 {
	return 2 * (BitWidth(int64(value)) - 1) + 1
}

func PreOffset(block int32) int32 {
	return int32(uint32(PrefixSum[block]) >> 24) & 0xFF
}

func PreCount(block int32) int32 {
	return int32(uint32(PrefixSum[block]) >> 16) & 0xFF
}


func PreSum(block int32) int32 {
	return int32(uint32(PrefixSum[block])) & 0xFFFF
}

