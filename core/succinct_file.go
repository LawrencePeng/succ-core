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
	ExtractUntil(offset int64, delim int) string
	ExtractUntilWith(offset int64, delim int, ctx *ExtractContext) string
	ExtractBytes(offset int64, len int, ctx ExtractContext) []byte
	ExtractBytesUntil(offset int64, delim int) []byte
	ExtractBytesUntilWith(ctx ExtractContext, delim int) []byte
	ExtractShortWith(offset int, ctx ExtractContext) int16
	ExtractShort(offset int) int16
	ExtractShortWithOffset(offset int) int16
	ExtractShortWithCtx(ctx ExtractContext) int16
	ExtractIntWith(offset int, ctx ExtractContext) int
	ExtractInt(offset int) int
	ExtractIntWithOffset(offset int) int
	ExtractIntWithCtx(ctx ExtractContext) int
	ExtractLongWith(offset int, ctx ExtractContext) int64
	ExtractLong(offset int) int64
	ExtractLongWithOffset(offset int) int64
	ExtractLongWithCtx(ctx ExtractContext) int64
	RangeSearch(buf1, buf2 []byte) Range
	RangeSearchWithSource(s1, s2 Source) Range
	BwdSearch(buf []byte) Range
	BwdSearchWithSource(s Source) Range
	
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
