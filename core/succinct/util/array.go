package util

func GetRank164(arrayBuf []int64, startPos int32, size int32, i int64) int32 {
	sp := int32(0)
	ep := size - 1

	var m int32

	for ; sp <= ep; {
		m = (sp + ep) / 2
		if arrayBuf[startPos + m] == i {
			return m + 1
		} else if i < arrayBuf[startPos + m] {
			ep = m - 1
		} else {
			sp = m + 1
		}
	}

	return ep + 1
}

func GetRank132(arrayBuf []int32, startPos int32, size int32, i int32) int32 {
	sp := int32(0)
	ep := size - 1

	var m int32

	for ; sp <= ep; {
		m = (sp + ep) / 2
		if arrayBuf[startPos + m] == i {
			return m + 1
		} else if i < arrayBuf[startPos + m] {
			ep = m - 1
		} else {
			sp = m + 1
		}
	}

	return ep + 1
}
