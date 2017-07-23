package succinct

import (
	"bytes"
	"os"

	"./util"
	"math"
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

func ReadSuccinctFileBufferFromFile(file *os.File) (*SuccinctFileBuffer, *bytes.Buffer) {
	succBuf, buf := ReadSuccinctBufferFromFile(file)
	return &SuccinctFileBuffer{
		SuccBuf: succBuf,
	}, buf
}

func (succFBuf *SuccinctFileBuffer) Size() int32 {
	return succFBuf.SuccBuf.Core.OriginalSize
}

func (succFBuf *SuccinctFileBuffer) CompressedSize() int32 {
	return succFBuf.SuccBuf.CoreSize()
}

func (succFBuf *SuccinctFileBuffer) CharAt(i int64) byte {
	return byte(succFBuf.SuccBuf.LookUpC(succFBuf.SuccBuf.LookUpISA(i)))
}

func (succFBuf *SuccinctFileBuffer) ExtractWith(offset int64, len int32, ctx *ExtractContext) string {
	buf := new(bytes.Buffer)
	s := succFBuf.SuccBuf.LookUpISA(offset)
	for k := int32(0); k <= len && int32(offset)+k <= succFBuf.Size(); k++ {
		nextCh := succFBuf.SuccBuf.LookUpC(s)
		if nextCh < 0 || nextCh > 0xFFFF {
			break
		}

		buf.WriteByte(byte(nextCh))
		s = succFBuf.SuccBuf.LookUpNPA(s)
	}

	if ctx != nil {
		ctx.Marker = s
	}

	return string(buf.Bytes())
}

func (succFBuf *SuccinctFileBuffer) Extract(offset int64, len int32) string {
	return succFBuf.ExtractWith(offset, len, nil)
}

func (succFBuf *SuccinctFileBuffer) ExtractUntilWith(offset int64, delim int32, ctx *ExtractContext) string {
	buf := new(bytes.Buffer)

	s := succFBuf.SuccBuf.LookUpISA(offset)

	var nextCh int32
	nextCh = succFBuf.SuccBuf.LookUpC(s)
	if nextCh == delim || nextCh == int32(util.EOF) {
		if ctx != nil {
			ctx.Marker = s
		}
		return string(buf.Bytes())
	}
	buf.WriteByte(byte(nextCh))
	s = succFBuf.SuccBuf.LookUpNPA(s)

	for {
		nextCh = succFBuf.SuccBuf.LookUpC(s)
		if nextCh == delim || nextCh == int32(util.EOF) {
			break
		}
		buf.WriteByte(byte(nextCh))
		s = succFBuf.SuccBuf.LookUpNPA(s)
	}
	if ctx != nil {
		ctx.Marker = s
	}
	return string(buf.Bytes())
}

func (succFBuf *SuccinctFileBuffer) ExtractUntil(offset int64, delim int32) string {
	return succFBuf.ExtractUntilWith(offset, delim, nil)
}

func (succFBuf *SuccinctFileBuffer) ExtractBytesWith(offset int64, len int32, ctx *ExtractContext) []byte {
	buf := new(bytes.Buffer)
	s := succFBuf.SuccBuf.LookUpISA(offset)
	for k := int32(0); k < len && int32(offset)+k < succFBuf.Size(); k++ {
		nextCh := succFBuf.SuccBuf.LookUpC(s)
		buf.WriteByte(byte(nextCh))
		s = succFBuf.SuccBuf.LookUpNPA(s)
	}

	if ctx != nil {
		ctx.Marker = s
	}

	return buf.Bytes()
}

func (succFBuf *SuccinctFileBuffer) ExtractBytes(offset int64, len int32) []byte {
	return succFBuf.ExtractBytesWith(offset, len, nil)
}

func (succFBuf *SuccinctFileBuffer) BwdSearch(source *SuccinctSource) *util.Range {
	ran := &util.Range{
		From: 0,
		To:   -1,
	}

	m := int32((*source).Len())
	var c1, c2 int64

	alphaSize := int32(len(succFBuf.SuccBuf.Core.Alphabet))

	pos := succFBuf.SuccBuf.Core.FindCharacter((*source).Get(m - 1))
	if pos >= 0 {
		ran.From = int64(succFBuf.SuccBuf.ColumnOffsets[pos])
		if pos+1 == alphaSize {
			ran.To = int64(succFBuf.SuccBuf.Core.OriginalSize)
		} else {
			ran.To = int64(succFBuf.SuccBuf.ColumnOffsets[pos+1] - 1)
		}
	} else {
		return &util.Range{From: 0, To: -1}
	}

	for i := m - 2; i >= 0; i-- {
		pos = succFBuf.SuccBuf.Core.FindCharacter(source.Get(i))
		if pos >= 0 {
			c1 = int64(succFBuf.SuccBuf.ColumnOffsets[pos])
			if pos+1 == alphaSize {
				c2 = int64(succFBuf.SuccBuf.Core.OriginalSize)
			} else {
				c2 = int64(succFBuf.SuccBuf.ColumnOffsets[pos+1] - 1)
			}
		} else {
			return &util.Range{From: 0, To: -1}
		}

		if c1 > c2 {
			return &util.Range{From: 0, To: -1}
		}

		ran.From = succFBuf.SuccBuf.BinSearchNPA(ran.From, c1, c2, false)
		ran.To = succFBuf.SuccBuf.BinSearchNPA(ran.To, c1, c2, true)

		if ran.From > ran.To {
			return &util.Range{From: 0, To: -1}
		}
	}

	return ran
}

func (succFBuf *SuccinctFileBuffer) BwdSearchStr(str string) *util.Range {
	return succFBuf.BwdSearch(&SuccinctSource{
		Bts: []byte(str),
	})
}

func (succFBuf *SuccinctFileBuffer) FwdSearchStr(str string) *util.Range {
	return succFBuf.FwdSearchWithSource(&SuccinctSource{
		Bts: []byte(str),
	})
}

func (succFBuf *SuccinctFileBuffer) ContinueBwdSearchStr(q string, ran *util.Range) *util.Range  {
	return succFBuf.ContinueBwdSearch(q, ran)
}

func (succFBuf *SuccinctFileBuffer) ContinueBwdSearchWithSource(source *SuccinctSource, r *util.Range) *util.Range {
	if r.Empty() {
		return r
	}

	newRange := &util.Range{From: r.From, To: r.To}
	m := (*source).Len()
	var c1, c2 int64

	alphaSize := int32(len(succFBuf.SuccBuf.Core.Alphabet))
	for i := m - 1; i >= 0; i-- {
		pos := succFBuf.SuccBuf.Core.FindCharacter((*source).Get(i))
		if pos >= 0 {
			c1 = int64(succFBuf.SuccBuf.ColumnOffsets[pos])
			if pos+1 == alphaSize {
				c2 = int64(succFBuf.SuccBuf.Core.OriginalSize)
			} else {
				c2 = int64(succFBuf.SuccBuf.ColumnOffsets[pos+1] - 1)
			}
		} else {
			return &util.Range{
				From: 0,
				To:   -1,
			}
		}

		if c1 > c2 {
			return &util.Range{
				From: 0,
				To:   -1,
			}
		}

		newRange.From = succFBuf.SuccBuf.BinSearchNPA(newRange.From, c1, c2, false)
		newRange.To = succFBuf.SuccBuf.BinSearchNPA(newRange.To, c1, c2, true)

		if newRange.From > newRange.To {
			return &util.Range{
				From: 0,
				To:   -1,
			}
		}
	}

	return newRange

}

func (succFBuf *SuccinctFileBuffer) ContinueBwdSearch(q string, ran *util.Range) *util.Range {
	return succFBuf.ContinueBwdSearchWithSource(&SuccinctSource{
		Bts: []byte(q),
	}, ran)
}

func (succFBuf *SuccinctFileBuffer) Compare(source *SuccinctSource, i int64) int32 {
	j := int32(0)

	var c, b int32

	c = succFBuf.SuccBuf.LookUpC(i)
	b = (*source).Get(j)

	if b < c {
		return -1
	} else if b > c {
		return 1
	}
	i = succFBuf.SuccBuf.LookUpNPA(i)

	for j < (*source).Len() {
		c = succFBuf.SuccBuf.LookUpC(i)
		b = (*source).Get(j)

		if b < c {
			return -1
		} else if b > c {
			return 1
		}
		i = succFBuf.SuccBuf.LookUpNPA(i)
	}
	return 0
}

func (succFBuf *SuccinctFileBuffer) CompareWithSourceAndOffSet(s *SuccinctSource, i int32, offset int32) int32 {
	j := int32(0)

	for ; offset != 0; offset-- {
		i = int32(succFBuf.SuccBuf.LookUpNPA(int64(i)))
	}

	var c, b int32

	c = succFBuf.SuccBuf.LookUpC(int64(i))
	b = (*s).Get(j)
	if b < c {
		return -1
	} else if b > c {
		return 1
	}

	i = int32(succFBuf.SuccBuf.LookUpNPA(int64(i)))
	j++

	for j < (*s).Len() {
		c = succFBuf.SuccBuf.LookUpC(int64(i))
		b = (*s).Get(j)
		if b < c {
			return -1
		} else if b > c {
			return 1
		}

		i = int32(succFBuf.SuccBuf.LookUpNPA(int64(i)))
		j++
	}
	return 0
}

func (succFBuf *SuccinctFileBuffer) FwdSearchWithSource(source *SuccinctSource) *util.Range {
	st := succFBuf.SuccBuf.Core.OriginalSize - 1

	sp := int32(0)
	var s int32
	for sp < st {
		s = (sp + st) / 2
		if succFBuf.Compare(source, int64(s)) > 0 {
			sp = s + 1
		} else {
			st = s
		}
	}

	et := succFBuf.SuccBuf.Core.OriginalSize - 1

	ep := sp - 1
	var e int32
	for ep < et {
		e = int32(math.Ceil(float64((ep + et) / 2)))
		if succFBuf.Compare(source, int64(e)) == 0 {
			ep = e
		} else {
			et = e - 1
		}
	}

	return &util.Range{
		From: int64(sp),
		To:   int64(ep),
	}
}


func (succFBuf *SuccinctFileBuffer) ContinueFwdSearchWithQuery(q string, r *util.Range, offset int32) *util.Range {
	return succFBuf.ContinueFwdSearchWithSource(&SuccinctSource{Bts:[]byte(q)}, r, offset)
}


func (succFBuf *SuccinctFileBuffer) ContinueFwdSearchWithSource(source *SuccinctSource, r *util.Range, offset int32) *util.Range {
	if source.Len() == 0 || r.Empty() {
		return r
	}

	st := int32(r.To)
	sp := int32(r.From)
	var s int32

	for sp < st {
		s = (sp + st) / 2
		if succFBuf.CompareWithSourceAndOffSet(source, s, offset) > 0 {
			sp = sp + 1
		} else {
			st = s
		}
	}

	et := int32(r.To)
	ep := sp - 1
	var e int32

	for ep < et {
		e = int32(math.Ceil(float64((ep + et) / 2)))
		if succFBuf.CompareWithSourceAndOffSet(source, e, offset) == 0 {
			ep = e
		} else {
			et = e - 1
		}
	}

	return &util.Range{
		From: int64(sp),
		To:   int64(ep),
	}
}

func (succFBuf *SuccinctFileBuffer) Count(query *SuccinctSource) int64 {
	r := succFBuf.BwdSearch(query)
	return r.To - r.From + 1
}

func (succFBuf *SuccinctFileBuffer) SuccinctIndexToOffset(i int64) int64 {
	return succFBuf.SuccBuf.LookUpSA(i)
}

func (succFBuf *SuccinctFileBuffer) RangeToOffsets(r *util.Range) []int64 {
	if r.Empty() {
		return []int64{}
	}

	offsets := make([]int64, r.Size())
	for i := int64(0); i < r.Size(); i++ {
		offsets[i] = succFBuf.SuccBuf.LookUpSA(r.From + i)
	}

	return offsets
}

func (succFBuf *SuccinctFileBuffer) Alphabet() []int32 {
	return succFBuf.SuccBuf.Core.Alphabet
}

func (succFBuf *SuccinctFileBuffer) Search(query *SuccinctSource) []int64 {
	return succFBuf.RangeToOffsets(succFBuf.BwdSearch(query))
}

func (succFBuf *SuccinctFileBuffer) SameRecord(fir, sec int32) bool {
	return true
}
