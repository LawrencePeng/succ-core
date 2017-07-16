package succinct

import (
	"./util"
	"os"
	"bytes"
	"math"
)

type SuccinctIndexedFileBuffer struct {
	SuccFBuf *SuccinctFileBuffer
	Offsets  []int32
}

func BuildSuccinctIndexedFileBufferFromInput(source *string, offset []int32,
	conf *util.SuccinctConf) *SuccinctIndexedFileBuffer {
	succFBuf := BuildSuccinctFileBufferFromInput(*source, conf)
	return &SuccinctIndexedFileBuffer{
		SuccFBuf: succFBuf,
		Offsets:  offset,
	}
}

func (succIFB *SuccinctIndexedFileBuffer) WriteToFile(file *os.File) {
	buf := new(bytes.Buffer)
	succIFB.SuccFBuf.SuccBuf.WriteToBuf(buf)
	util.WriteArray(buf, succIFB.Offsets)
	buf.WriteTo(file)
}

func ReadSuccinctIndexFileBufferFromFile(file *os.File) *SuccinctIndexedFileBuffer {
	succFBuf, buf := ReadSuccinctFileBufferFromFile(file)
	if succFBuf == nil || buf == nil {
		return nil
	}
	offsets := util.ReadArray(buf)
	return &SuccinctIndexedFileBuffer{
		SuccFBuf: succFBuf,
		Offsets:  offsets,
	}
}

func (succIFB *SuccinctIndexedFileBuffer) CompressedSize() int32 {
	return succIFB.SuccFBuf.CompressedSize() + int32(len(succIFB.Offsets) * util.INT_SIZE)
}

func (succIFB *SuccinctIndexedFileBuffer) RecordOffset(recordId int32) int32 {
	return succIFB.Offsets[recordId]
}

func (succIFB *SuccinctIndexedFileBuffer) RecordBytes(recordId int32) []byte {
	if int(recordId) > len(succIFB.Offsets) || recordId < 0 {
		panic("wrong recordid in RecordBytes")
	}

	begOffset := succIFB.Offsets[recordId]
	var endOffset int32
	if int(recordId) == len(succIFB.Offsets) - 1 {
		endOffset = succIFB.SuccFBuf.SuccBuf.Core.OriginalSize - 1
	} else {
		endOffset = succIFB.Offsets[recordId + 1]
	}

	length := endOffset - begOffset - 1
	return succIFB.SuccFBuf.ExtractBytes(int64(begOffset), length)
}

func (succIFB *SuccinctIndexedFileBuffer) Record(recordId int32) string {
	if int(recordId) > len(succIFB.Offsets) || recordId < 0 {
		panic("wrong recordid in Record")
	}

	begOffset := succIFB.Offsets[recordId]
	var endOffset int32
	if int(recordId) == len(succIFB.Offsets) - 1 {
		endOffset = succIFB.SuccFBuf.SuccBuf.Core.OriginalSize - 1
	} else {
		endOffset = succIFB.Offsets[recordId + 1]
	}

	length := endOffset - begOffset - 1
	return succIFB.SuccFBuf.Extract(int64(begOffset), length)
}

func (succIFB *SuccinctIndexedFileBuffer) ExtractRecord(recordId int32, offset, length int32) string {
	if int(recordId) > len(succIFB.Offsets) || recordId < 0 {
		panic("wrong recordid in ExtractRecord")
	}

	if length == 0 {
		return ""
	}

	begOffset := succIFB.Offsets[recordId] + offset
	var nextRecordOffset int32
	if int(recordId) == len(succIFB.Offsets) - 1 {
		nextRecordOffset = succIFB.SuccFBuf.SuccBuf.Core.OriginalSize - 1
	} else {
		nextRecordOffset = succIFB.Offsets[recordId + 1]
	}

	length = int32(math.Min(float64(nextRecordOffset - begOffset - 1), float64(length)))
	return succIFB.SuccFBuf.Extract(int64(begOffset), length)
}

func (succIFB *SuccinctIndexedFileBuffer) OffsetToRecordId(pos int32) int32 {
	sp := int32(0)
	ep := int32(len(succIFB.Offsets) - 1)

	var m int32

	for ; sp <= ep; {
		m = (sp + ep) / 2
		if succIFB.Offsets[m] == pos {
			return m
		} else if pos < succIFB.Offsets[m] {
			ep = m - 1
		} else {
			sp = m + 1
		}
	}

	return ep
}

func (succIFB *SuccinctIndexedFileBuffer) RecordSearchIds(q *SuccinctSource) []int32 {
	results := &util.HashSet{M: make(map[int32]bool)}
	r := succIFB.SuccFBuf.BwdSearch(q)

	sp := r.From
	ep := r.To

	if ep - sp + 1 <= 0 {
		return []int32{}
	}

	for i := int64(0); i < ep - sp + 1; i++ {
		results.Add(succIFB.OffsetToRecordId(int32(succIFB.SuccFBuf.SuccBuf.LookUpSA(sp + i))))
	}

	ret := make([]int32, 0)
	for k := range results.M {
		ret = append(ret, k)
	}
	return ret
}

func (succIFB *SuccinctIndexedFileBuffer) SameRecord(fir, sec int64) bool {
	return succIFB.OffsetToRecordId(int32(fir)) == succIFB.OffsetToRecordId(int32(sec))
}

