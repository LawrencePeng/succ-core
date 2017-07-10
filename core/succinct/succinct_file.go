package succinct

import (
	. "./util"
)

type ExtractContext struct {
	Marker int64
}

type SuccinctFile interface {
	GetAlphabet() []int32
	GetSize() int32
	GetCompressedSize() int32
	CharAt(i int64) byte
	ExtractWith(offset int64, len int32, ctc *ExtractContext) string
	Extract(offset int64, len int32) string
	ExtractUntil(offset int64, delim int32) string
	ExtractUntilWith(offset int64, delim int32, ctx *ExtractContext) string
	ExtractBytes(offset int64, len int32, ctx ExtractContext) []byte
	ExtractBytesUntil(offset int64, delim int32) []byte
	ExtractBytesUntilWith(ctx ExtractContext, delim int32) []byte
	ExtractShortWith(offset int32, ctx ExtractContext) int16
	ExtractShort(offset int32) int16
	ExtractShortWithOffset(offset int32) int16
	ExtractShortWithCtx(ctx ExtractContext) int16
	ExtractIntWith(offset int32, ctx ExtractContext) int32
	ExtractInt(offset int32) int32
	ExtractIntWithOffset(offset int32) int32
	ExtractIntWithCtx(ctx ExtractContext) int32
	ExtractLongWith(offset int, ctx ExtractContext) int64
	ExtractLong(offset int32) int64
	ExtractLongWithOffset(offset int32) int64
	ExtractLongWithCtx(ctx ExtractContext) int64
	RangeSearch(buf1, buf2 []byte) Range
	RangeSearchWithSource(s1, s2 Source) Range
	BwdSearch(buf []byte) Range
	BwdSearchWithSource(s Source) Range
	ContinueBwdSearch(buf []byte, r Range) Range
	ContinueBwdSearchWithSource(source Source, r Range) Range
	Compare(buf []byte, i int32) int32
	CompareWithSource(s Source, i int32) int32
	CompareWithOffset(buf []byte, i int32, offset int32) int32
	CompareWithSourceAndOffSet(s Source, i int32, offset int32) int32
	FwdSearch(buf []byte) Range
	FwdSearchWithSource(s Source) Range
	ContinueFwdSearch(buf []byte, r Range) Range
	ContinueFwdSearchWithSource(source Source, r Range) Range
	Count(q []byte) int64
	CountWithSource(s Source) int64
	SuccinctIndexOffsets(r Range) []int64
	Search(query []byte) []int64
	SearchWithSource(s Source) []int64
	SameRecord(fir, sec int64) bool
}

type SuccinctIndexedFile interface {
	SuccinctFile
	OffsetToRecordId(pos int32) int32
	GetNumRecords() int32
	GetRecordOffset(recordId int32) int32
	GetRecordBytes(recordId int32) []byte
	ExtractRecordBytes(recordId int32, offset int32, len int32) []byte
	GetRecord(recordId int32) string
	ExtractRecord(recordId int32, offset int32, length int32) string
	RecordSearchIds(query Source) []int32
	RecordSearchIdsByByteArr(query []byte) []int32
}


