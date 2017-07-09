package util

import "bytes"

type DeltaIntVector struct {
	Samples			*IntVector
	DeltaOffsets	*IntVector
	Deltas			*BitVector
	SamplingRate	int32
}

func NewDeltaIntVector(elements []int32 , samplingRate int32) *DeltaIntVector {
	div := &DeltaIntVector{
		SamplingRate:samplingRate,
	}
	div.encode(elements, 0, int32(len(elements)))
	return div
}

func (div *DeltaIntVector) encode(elements []int32, startingOffset int32, length int32) {
	if length == 0 {
		return
	}

	samplesBuf := NewIntArrayList()
	deltasBuf := NewIntArrayList()
	deltaOffsetsBuf := NewIntArrayList()

	maxSample := int32(0)
	lastValue := int32(0)
	totalDeltaCount := int32(0)
	deltaCount := int32(0)
	cumulativeDeltaSize := int32(0)
	maxOffset := int32(0)

	for i := int32(0); i < length; i++ {
		if i % div.SamplingRate == 0 {
			samplesBuf.Add(elements[startingOffset + i])
			if elements[startingOffset + 1] > maxSample {
				maxSample = elements[startingOffset + 1]
			}
			if cumulativeDeltaSize > maxOffset {
				maxOffset = cumulativeDeltaSize
			}
			deltaOffsetsBuf.Add(cumulativeDeltaSize)
			if i != int32(0) {
				totalDeltaCount += deltaCount
				deltaCount = 0
			}
		} else {
			delta := elements[startingOffset + i] - lastValue
			deltasBuf.Add(delta)
			cumulativeDeltaSize += EncodingSize(delta)
			deltaCount++
		}
		lastValue = elements[startingOffset + i]
	}
	totalDeltaCount += deltaCount

	sampleBits := BitWidth(int64(maxSample))
	deltaOffsetBits := BitWidth(int64(maxOffset))

	if samplesBuf.Size() == 0 {
		div.Samples = nil
	} else {
		div.Samples = CopyFromIntArrayList(samplesBuf, sampleBits)
	}

	if cumulativeDeltaSize == 0 {
		div.Deltas = nil
	} else {
		div.Deltas = NewBitVectorWithSize(int64(cumulativeDeltaSize + 16))
		div.encodeDeltas(deltasBuf)
	}

	if deltaOffsetBits == 0 {
		div.DeltaOffsets = nil
	} else {
		div.DeltaOffsets = CopyFromIntArrayList(deltaOffsetsBuf, deltaOffsetBits)
	}

}
func (div *DeltaIntVector) encodeDeltas(deltasArr *IntArrayList) {
	pos := int64(0)
	for i := int32(0); i < deltasArr.Size(); i++ {
		deltaBits := int32(BitWidth(int64(deltasArr.Get(i))) - 1)
		pos += int64(deltaBits)

		div.Deltas.SetBit(pos)
		pos++
		div.Deltas.SetValue(pos, int64(deltasArr.Get(i)) - int64(1) << uint32(deltaBits), deltaBits)
		pos += int64(deltaBits)
	}
}

func (div *DeltaIntVector) prefixSum(deltaOffset, untilIdx int32) int32 {
	deltaSum := int32(0)
	deltaIdx := int64(0)
	currentDeltaOffset := deltaIdx
	for ;deltaIdx != int64(untilIdx); {
		block := int32(div.Deltas.GetValue(currentDeltaOffset, 16))
		cnt := PreCount(block)
		if cnt == 0 {
			deltaWidth := int32(0)
			for ;div.Deltas.GetBit(currentDeltaOffset) != 1; {
				deltaWidth++
				currentDeltaOffset++
			}
			currentDeltaOffset++
			deltaSum += int32(div.Deltas.GetValue(currentDeltaOffset, deltaWidth)) + int32(int64(1) << uint64(deltaWidth))
			currentDeltaOffset += int64(deltaWidth)
			deltaIdx += 1
		} else if deltaIdx + int64(cnt) <= int64(untilIdx) {
			deltaSum += PreSum(block)
			currentDeltaOffset += int64(PreOffset(block))
			deltaIdx += int64(cnt)
		} else {
			for ; deltaIdx != int64(untilIdx); {
				deltaWidth := int32(0)
				for ; div.Deltas.GetBit(currentDeltaOffset) != 1;  {
					deltaWidth ++
					currentDeltaOffset++
				}
				currentDeltaOffset++
				deltaSum += int32(div.Deltas.GetValue(currentDeltaOffset, deltaWidth) + int64(1) << uint64(deltaWidth))
				currentDeltaOffset += int64(deltaWidth)
				deltaIdx += 1
			}
		}
	}

	return deltaSum
}

func (div *DeltaIntVector) Get(i int32) int32 {
	sampleIdx := i / div.SamplingRate
	deltaOffsetIdx := i % div.SamplingRate
	val := div.Samples.Get(sampleIdx)

	if deltaOffsetIdx == 0 {
		return val
	}

	deltaOffset := div.DeltaOffsets.Get(sampleIdx)
	val += div.prefixSum(deltaOffset, deltaOffsetIdx)
	return val
}

func (div *DeltaIntVector) SerializedSize() int32 {
	var samplesSize, deltaOffsetSize, deltaSize int32
	if div.Samples == nil {
		samplesSize = int32(INT_SIZE)
	} else {
		samplesSize = div.Samples.SerializedSize()
	}

	if div.DeltaOffsets == nil {
		deltaOffsetSize = int32(INT_SIZE)
	} else {
		deltaOffsetSize = div.DeltaOffsets.SerializedSize()
	}

	if div.Deltas == nil {
		deltaSize = int32(INT_SIZE)
	} else {
		deltaSize = div.Deltas.SerializedSize()
	}

	return int32(INT_SIZE) + samplesSize + deltaOffsetSize + deltaSize
}

func (div *DeltaIntVector) WriteToBuf(buf bytes.Buffer)  {
	WriteInt(buf, div.SamplingRate)
	if div.Samples != nil {
		div.Samples.WriteToBuf(buf)
	} else {
		WriteInt(buf, 0)
	}

	if div.DeltaOffsets != nil {
		div.DeltaOffsets.WriteToBuf(buf)
	} else {
		WriteInt(buf, 0)
	}

	if div.Deltas != nil {
		div.Deltas.WriteToBuf(buf)
	} else {
		WriteInt(buf, 0)
	}
}

func DIVPrefixSum(deltas []int64, deltaOffset int32, untilIdx int32) int32 {
	deltaSum := int32(0)
	deltaIdx := int64(0)

	currentDeltaOffset := deltaIdx
	for ; deltaIdx != int64(untilIdx); {
		block := int32(BitVecGetValue(deltas, currentDeltaOffset, 16))
		cnt := PreCount(block)
		if cnt == 0 {
			deltaWidth := int32(0)
			for ; BitVecGetBit(deltas, currentDeltaOffset) != -1; {
				deltaWidth++
				currentDeltaOffset++
			}

			currentDeltaOffset++
			deltaSum +=
				int32(BitVecGetValue(deltas, currentDeltaOffset, deltaWidth) + int64(1) << uint(deltaWidth))
			currentDeltaOffset += int64(deltaWidth)
			deltaIdx ++
		} else if deltaIdx + int64(cnt) <= int64(untilIdx) {
			deltaSum += PreSum(block)
			currentDeltaOffset += int64(PreOffset(block))
			deltaIdx += int64(cnt)
		} else {
			for ;deltaIdx != int64(untilIdx);  {
				deltaWidth := int32(0)
				for ; BitVecGetBit(deltas, currentDeltaOffset) != 1; {
					deltaWidth++
					currentDeltaOffset++
				}
				currentDeltaOffset++
				deltaSum +=
					int32(BitVecGetValue(deltas, currentDeltaOffset, deltaWidth) + int64(1) << uint(deltaWidth))
				currentDeltaOffset += int64(deltaWidth)
				deltaIdx += 1
			}
		}
	}
	return deltaSum
}

func DIVGet(buf bytes.Buffer, i int32) int32 {
	samplingRate := ReadInt(buf)
	sampleBits := ReadInt(buf)
	sampleBlocks := ReadInt(buf)

	samples :=

}