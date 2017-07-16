package util_test

import (
	"testing"
	"fmt"
)
import . "."
func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func TestIntLog2(t *testing.T) {
	assertEqual(t, int(IntLog2(0)), 1, "")
	assertEqual(t, int(IntLog2(1)), 0, "")
	assertEqual(t, int(IntLog2(2)), 1, "")
	assertEqual(t, int(IntLog2(3)), 2, "")
	assertEqual(t, int(IntLog2(4)), 2, "")
	assertEqual(t, int(IntLog2(5)), 3, "")
	assertEqual(t, int(IntLog2(6)), 3, "")
	assertEqual(t, int(IntLog2(99)), 7, "")
	assertEqual(t, int(IntLog2(-5)), -1, "")
}

func TestMod(t *testing.T) {
	assertEqual(t, Mod(-2, 3), int64(1), "")
	assertEqual(t, Mod(5, 2), int64(1), "")
	assertEqual(t, Mod(13, 13), int64(0), "")
	assertEqual(t, Mod(15, 17), int64(15), "")
}

func TestPopCount(t *testing.T) {
	assertEqual(t, PopCount(uint64(0)), int32(0), "")
	assertEqual(t, PopCount(0xFFFFFFFFFFFFFFFF), int32(64), "")
	assertEqual(t, PopCount(0xFFFF0000), int32(16), "")
}

func TestNumBlocks(t *testing.T) {
	assertEqual(t, NumBlocks(0, 5), int32(0), "a")
	for i := 1; i <= 5; i++ {
		assertEqual(t, NumBlocks(int32(i), 5), int32(1), "b")
	}
	assertEqual(t, NumBlocks(6, 5), int32(2), "c")
	assertEqual(t, NumBlocks(256, 5), int32(52), "d")
	assertEqual(t, NumBlocks(255, 5), int32(51), "e")
}
