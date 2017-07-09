package succinct

import "./util"

type SuccinctIndexedFileBuffer struct {
	sfb *SuccinctFileBuffer
	sa		[]int64
	isa		[]int64
	colOffset []int
	columns	  []byte
}


func NewSuccinctIndexedFileBuffer(input string, conf util.SuccinctConf) *SuccinctFileBuffer {

}