package core

import (
	. "./util"
)

type ExtractContext struct {
	Marker int64
}

type SuccinctFile interface {
	GetAlphabet() []int
	GetSize() int
	GetCompressedSize() int
	CharAt(i int64) byte
	ExtractWith(offset int64, len int, ctc *ExtractContext) string
	Extract(offset int64, len int) string
	ExtractUtil(offset int64, delim int) string
	ExtractUntilWith(offset int64, delim int, ctx *ExtractContext) string

}

type SuccinctIndexedFile interface {
	SuccinctFile
	OffsetToRecordId(pos int) int
	GetNumRecords() int
	GetRecordOffset(recordId int) int
	GetRecordByBytes(recordId int) []byte
	ExtractRecordBytes(recordId int, offset int, len int) []byte
	GetRecord(recordId int) string
	ExtractRecord(recordId int, offset int, length int) string
	RecordSearchIds(query Source) []int
	RecordSearchIdsByByteArr(query []byte) []int
}
