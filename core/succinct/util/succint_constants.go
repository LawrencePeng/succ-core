package util

import (
	"unsafe"
)

var BYTE_SIZE = 1
var SHORT_SIZE = 2
var INT_SIZE = 4
var LONG_SIZE = 8
var REF_SIZE_BYTE = int(unsafe.Sizeof(uintptr(0)))

var DEFAULT_SA_SAMPLING_SIZE = 32
var DEFAULT_ISA_SAMPLING_SIZE = 32
var DEFAULT_NSA_SAMPLING_SIZE = 128

var EOL = '\n'

