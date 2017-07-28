package succinct

import (
	"./util"
	"bytes"
	"encoding/binary"
	"github.com/juju/errors"
	"os"
)

type SuccinctBuffer struct {
	Core          *SuccinctCore
	SA            []int64
	ISA           []int64
	ColumnOffsets []int32
	Columns       [][]byte
}

func ReadSuccinctBufferFromFile(file *os.File) (*SuccinctBuffer,
	*bytes.Buffer, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	size := stat.Size()
	bts := make([]byte, size)
	file.Read(bts)
	buf := bytes.NewBuffer(bts)

	sb, err := mapFromBuf(buf)
	if err != nil {
		return nil, buf, err
	}

	return sb, buf, nil
}

func BuildSuccinctBufferFromInput(input *SuccinctSource,
	conf *util.SuccinctConf) (*SuccinctBuffer, error) {
	buf := new(bytes.Buffer)

	var ISA []int32
	var columnOffsets []int32

	originalSize := (*input).Len() + 1

	samplingRateSA := conf.SaSamplingRate
	samplingRateISA := conf.IsaSamplingRate
	samplingRateNPA := conf.NpaSamplingRate

	sampleBitWidth := util.BitWidth(int64(originalSize))

	var alphabetSize int32

	{
		suffixSorter := &util.QSufSort{}
		suffixSorter.BuildSuffixArray(input)
		SA := suffixSorter.I
		ISA = suffixSorter.V

		alphabetSize = int32(len(suffixSorter.Alphabet))

		util.WriteInt(buf, originalSize)
		util.WriteInt(buf, samplingRateSA)
		util.WriteInt(buf, samplingRateISA)
		util.WriteInt(buf, samplingRateNPA)
		util.WriteInt(buf, sampleBitWidth)
		util.WriteInt(buf, alphabetSize)

		alphabetArr := suffixSorter.Alphabet

		for i := int32(0); i < alphabetSize; i++ {
			util.WriteInt(buf, alphabetArr[i])
		}

		pos := int32(0)
		prevSortedChar := int32(util.EOF)
		columnOffsets = make([]int32, alphabetSize)
		columnOffsets[pos] = 0
		pos++

		for i := int32(1); i < originalSize; i++ {
			if (*input).Get(SA[i]) != prevSortedChar {
				prevSortedChar = (*input).Get(SA[i])
				columnOffsets[pos] = i
				pos++
			}
		}
	}

	{
		var sampledSA, sampledISA *util.IntVector
		numSampledElementsSA := util.NumBlocks(originalSize, samplingRateSA)
		numSampledElementsISA := util.NumBlocks(originalSize, samplingRateISA)
		sampledSA = util.NewIntVector(numSampledElementsSA, sampleBitWidth)
		sampledISA = util.NewIntVector(numSampledElementsISA, sampleBitWidth)
		for val := int32(0); val < originalSize; val++ {
			idx := ISA[val]
			if idx%samplingRateSA == 0 {
				sampledSA.Add(idx/samplingRateSA, val)
			}
			if val%samplingRateISA == 0 {
				sampledISA.Add(val/samplingRateISA, idx)
			}
		}

		sampledSA.Bv.WriteToDataBuf(buf)
		sampledISA.Bv.WriteToDataBuf(buf)
	}

	{
		for i := int32(0); i < alphabetSize; i++ {
			util.WriteInt(buf, columnOffsets[i])
		}

		NPA := make([]int32, originalSize)
		for i := int32(1); i < originalSize; i++ {
			NPA[ISA[i-1]] = ISA[i]
		}
		NPA[ISA[originalSize-1]] = ISA[0]

		for i := int32(0); i < alphabetSize; i++ {
			startOffset := columnOffsets[i]
			var endOffset int32
			if i < alphabetSize-1 {
				endOffset = columnOffsets[i+1]
			} else {
				endOffset = originalSize
			}
			length := endOffset - startOffset
			columnVector := util.NewDeltaIntVectorFull(NPA, startOffset, length, samplingRateNPA)
			sz := columnVector.SerializedSize()
			util.WriteInt(buf, sz)
			columnVector.WriteToBuf(buf)
		}
	}

	return mapFromBuf(buf)
}

func mapFromBuf(buf *bytes.Buffer) (*SuccinctBuffer, error) {
	core := &SuccinctCore{}

	// setup core
	core.OriginalSize = util.ReadInt(buf)
	core.SamplingRateSA = util.ReadInt(buf)
	core.SamplingRateISA = util.ReadInt(buf)
	core.SamplingRateNPA = util.ReadInt(buf)
	core.SampleBitWidth = util.ReadInt(buf)

	// read alphabet
	core.Alphabet = util.ReadArray(buf)
	alphabetSize := int32(len(core.Alphabet))

	succBuf := &SuccinctBuffer{
		Core: core,
	}

	// read sa
	totalSampledBitsSA := util.NumBlocks(core.OriginalSize, core.SamplingRateSA) *
		core.SampleBitWidth
	saSize := util.BitsToBlock64(int64(totalSampledBitsSA)) * int32(util.LONG_SIZE)
	succBuf.SA = util.ToLongSlice(buf.Next(int(saSize)))

	// read isa
	totalSampledBitsISA := util.NumBlocks(core.OriginalSize, core.SamplingRateISA) *
		core.SampleBitWidth
	isaSize := util.BitsToBlock64(int64(totalSampledBitsISA)) * int32(util.LONG_SIZE)
	succBuf.ISA = util.ToLongSlice(buf.Next(int(isaSize)))

	// read coloffsets
	coloffsetsSize := alphabetSize * int32(util.INT_SIZE)
	succBuf.ColumnOffsets = util.ToIntSlice(buf.Next(int(coloffsetsSize)))

	// read columns
	succBuf.Columns = make([][]byte, alphabetSize)
	for i := int32(0); i < alphabetSize; i++ {
		columnSize := util.ReadInt(buf)
		succBuf.Columns[i] = buf.Next(int(columnSize))
	}

	return succBuf, nil
}

func (succBuf *SuccinctBuffer) WriteToBuf(buf *bytes.Buffer) {
	util.WriteInt(buf, succBuf.Core.OriginalSize)
	util.WriteInt(buf, succBuf.Core.SamplingRateSA)
	util.WriteInt(buf, succBuf.Core.SamplingRateISA)
	util.WriteInt(buf, succBuf.Core.SamplingRateNPA)
	util.WriteInt(buf, succBuf.Core.SampleBitWidth)
	alphabetSize := int32(len(succBuf.Core.Alphabet))
	util.WriteInt(buf, alphabetSize)

	for i := int32(0); i < alphabetSize; i++ {
		util.WriteInt(buf, succBuf.Core.Alphabet[i])
	}

	for i := int32(0); i < int32(len(succBuf.SA)); i++ {
		util.WriteLong(buf, succBuf.SA[i])
	}

	for i := int32(0); i < int32(len(succBuf.ISA)); i++ {
		util.WriteLong(buf, succBuf.ISA[i])
	}

	for i := int32(0); i < int32(len(succBuf.ColumnOffsets)); i++ {
		util.WriteInt(buf, succBuf.ColumnOffsets[i])
	}

	for i := int32(0); i < int32(len(succBuf.Columns)); i++ {
		util.WriteInt(buf, int32(len(succBuf.Columns[i])))
		binary.Write(buf, binary.BigEndian, succBuf.Columns[i])
	}
}

func (succBuf *SuccinctBuffer) CoreSize() int32 {
	coreSize := succBuf.Core.BaseSize()                                  // core size
	coreSize += int32(len(succBuf.SA)) * int32(util.LONG_SIZE)           // sa size
	coreSize += int32(len(succBuf.ISA)) * int32(util.LONG_SIZE)          // isa size
	coreSize += int32(len(succBuf.ColumnOffsets)) * int32(util.INT_SIZE) // coloff size
	for i := int32(0); i < int32(len(succBuf.Columns)); i++ {            // col size
		coreSize += int32(len(succBuf.Columns[i])) * int32(util.BYTE_SIZE)
	}
	return coreSize
}

func (succBuf *SuccinctBuffer) LookUpNPA(i int64) (int64, error) {
	if i > int64(succBuf.Core.OriginalSize-1) || i < 0 {
		return -1, errors.New("wrong range of i in LookUpNPA")
	}

	alphabetSize := int32(len(succBuf.Core.Alphabet))

	colId := util.GetRank132(succBuf.ColumnOffsets, 0, alphabetSize, int32(i)) - 1

	if colId >= alphabetSize || int64(succBuf.ColumnOffsets[colId]) > i {
		return -1, errors.New("LookUpNPA Wrong colId")
	}

	return int64(util.DIVGet(&succBuf.Columns[colId], int32(i)-succBuf.ColumnOffsets[colId])), nil
}

func (succBuf *SuccinctBuffer) LookUpSA(i int64) (int64, error) {
	if i > int64(succBuf.Core.OriginalSize-1) || i < 0 {
		return -1, errors.New("wrong range of i in LookUpSA")
	}

	var err error = nil
	j := int32(0)
	for int32(i)%succBuf.Core.SamplingRateSA != 0 {
		i, err = succBuf.LookUpNPA(i)
		if err != nil {
			return -1, err
		}
		j++
	}

	saVal := util.IntVecGet(succBuf.SA, int32(i)/succBuf.Core.SamplingRateSA, succBuf.Core.SampleBitWidth)
	if saVal < j {
		return int64(succBuf.Core.OriginalSize - (j - saVal)), nil
	}
	return int64(saVal - j), nil
}

func (succBuf *SuccinctBuffer) LookUpISA(i int64) (int64, error) {
	if i > int64(succBuf.Core.OriginalSize-1) || i < 0 {
		return -1, errors.New("wrong range of i in LookUpSA")
	}

	var err error = nil
	var neoPos int64

	sampleIdx := int32(i) / succBuf.Core.SamplingRateISA
	pos := util.IntVecGet(succBuf.ISA, sampleIdx, succBuf.Core.SampleBitWidth)
	i -= int64(sampleIdx * succBuf.Core.SamplingRateISA)
	for i != 0 {
		i--
		neoPos, err = succBuf.LookUpNPA(int64(pos))
		if err != nil {
			return -1, err
		}
		pos = int32(neoPos)
	}

	return int64(pos), nil
}

func (succBuf *SuccinctBuffer) LookUpC(i int64) int32 {
	alphaSize := int32(len(succBuf.Core.Alphabet))
	idx := util.GetRank132(succBuf.ColumnOffsets, 0, alphaSize, int32(i)) - 1
	return succBuf.Core.Alphabet[idx]
}

func (succBuf *SuccinctBuffer) BinSearchNPA(val, startIdx, endIdx int64, flag bool) int64 {
	if endIdx < startIdx {
		return endIdx
	}

	alphaSize := int32(len(succBuf.Core.Alphabet))
	colId := util.GetRank132(succBuf.ColumnOffsets, 0, alphaSize, int32(startIdx)) - 1
	colValue := succBuf.ColumnOffsets[colId]

	sp := int32(startIdx) - colValue
	ep := int32(endIdx) - colValue

	res := util.BinarySearch(&succBuf.Columns[colId], int32(val), sp, ep, flag)
	return int64(colValue + res)
}
