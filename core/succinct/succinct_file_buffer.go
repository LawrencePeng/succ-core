package succinct

import (
	"./util"
	"os"
	"bytes"
)

type SuccinctFileBuffer struct {
	SuccBuf *SuccinctBuffer
}

type SuccinctSource struct {
	Bts []byte
}

func (succ *SuccinctSource) Len() int32 {
	return int32(len(succ.Bts))
}

func (succ *SuccinctSource) Get(i int32) int32 {
	return int32(succ.Bts[i])
}

func BuildSuccinctFileBufferFromInput(input string, conf *util.SuccinctConf) *SuccinctFileBuffer {
	bts := []byte(input)

	succBuf := BuildSuccinctBufferFromInput(&SuccinctSource{
		Bts: bts,
	}, conf)
	return &SuccinctFileBuffer{
		SuccBuf: succBuf,
	}
}

func ReadSuccinctFileBufferFromFile(file *os.File) *SuccinctFileBuffer {
	succBuf := ReadSuccinctBufferFromFile(file)
	return &SuccinctFileBuffer{
		SuccBuf: succBuf,
	}
}

func (succFBuf *SuccinctFileBuffer) Size() int32 {
	return succFBuf.SuccBuf.Core.OriginalSize
}

func (succFBuf *SuccinctFileBuffer) CompressedSize() int32 {
	return succFBuf.SuccBuf.CoreSize()
}

func (succFBuf *SuccinctFileBuffer) CharAt(i int64) byte {
	return byte(succFBuf.SuccBuf.LookUpC(succFBuf.SuccBuf.LookUpSA(i)))
}

func (succFBuf *SuccinctFileBuffer) ExtractWith(offset int64, len int32, ctx *ExtractContext) string {
	buf := new(bytes.Buffer)
	s := succFBuf.SuccBuf.LookUpISA(offset)
	for k := int32(0); k < len && int32(offset) + k < succFBuf.Size(); k++ {
		nextCh := succFBuf.SuccBuf.LookUpC(s)
		buf.WriteByte(byte(nextCh))
		s = succFBuf.SuccBuf.LookUpNPA(s)
	}

	if ctx != nil {
		ctx.Marker = s
	}

	return string(buf.Bytes())
}

