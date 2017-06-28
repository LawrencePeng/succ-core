package util

import (
	"unsafe"
)

const (
	BYTE_SIZE = 1
	SHORT_SIZE = 2
	INT_SIZE = unsafe.Sizeof(int(0))
	LONG_SIZE = 4
	REF_SIZE_BYTE = unsafe.Sizeof(uintptr(0))

	DEFAULT_SA_SAMPLING_SIZE = 32
	DEFAULT_ISA_SAMPLING_SIZE = 32
	DEFAULT_NSA_SAMPLING_SIZE = 128

	EOL = '\n'
)
