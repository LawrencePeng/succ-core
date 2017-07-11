package util

import (
	"encoding/binary"
	"bytes"
)

func ReadArray(buf *bytes.Buffer) []int32 {
	length := ReadInt(buf)
	arr := make([]int32, length)
	for i := 0; i < int(length); i++ {
		arr[i] = ReadInt(buf)
	}
	return arr
}

func WriteArray(buf *bytes.Buffer, arr []int32)  {
	WriteInt(buf, int32(len(arr)))
	for _, v := range arr {
		WriteInt(buf, v)
	}
}

func ReadInt(buf *bytes.Buffer) int32 {
	return int32(binary.BigEndian.Uint32(buf.Next(INT_SIZE)))
}

func ReadLong(buf *bytes.Buffer) int64 {
	return int64(binary.BigEndian.Uint64(buf.Next(LONG_SIZE)))
}

func WriteInt(buf *bytes.Buffer, v int32) {
	binary.Write(buf, binary.BigEndian, v)
}

func WriteLong(buf *bytes.Buffer, v int64) {
	binary.Write(buf, binary.BigEndian, v)
}

func CheckBytes(input []byte) int {
	for i := 0; i < len(input); i++ {
		if input[i] < 0 {
			return i
		}
	}
	return -1
}

