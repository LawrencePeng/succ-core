package succinct_test

import "testing"
import . "."
import (
	"./util"
	"fmt"
	"github.com/lexandro/go-assert"
	"os"
)

func TestSuccinctIndexedFileBuffer_Record(t *testing.T) {
	f, err := os.Open("./resources/test_file")
	if err != nil {
		panic("")
	}
	stat, err := f.Stat()
	if err != nil {
		panic("")
	}
	size := stat.Size()
	bts := make([]byte, size)
	f.Read(bts)

	pos := make([]int32, 0)
	pos = append(pos, int32(0))
	for i := int64(0); i < size; i++ {
		if bts[i] == '\n' {
			pos = append(pos, int32(i+1))
		}
	}

	input := string(bts)
	succ := BuildSuccinctIndexedFileBufferFromInput(&input, pos, &util.SuccinctConf{
		SaSamplingRate:  int32(32),
		IsaSamplingRate: int32(32),
		NpaSamplingRate: int32(128),
	})

	for i := 0; i < len(pos); i++ {
		assert.Equals(t, pos[i], succ.RecordOffset(int32(i)))
	}

	for i := 0; i < len(pos); i++ {
		fmt.Println(succ.Record(int32(i)))
	}

	ids := succ.RecordSearchIds(&SuccinctSource{Bts: []byte("int")})
	for i := 0; i < len(ids); i++ {
		fmt.Println(ids[i])
	}
}

func TestReadSuccinctFileBufferFromFile(t *testing.T) {
	f, err := os.Open("./resources/test_file")
	if err != nil {
		panic("")
	}
	stat, err := f.Stat()
	if err != nil {
		panic("")
	}
	size := stat.Size()
	bts := make([]byte, size)
	f.Read(bts)

	pos := make([]int32, 0)
	pos = append(pos, int32(0))
	for i := int64(0); i < size; i++ {
		if bts[i] == '\n' {
			pos = append(pos, int32(i+1))
		}
	}

	input := string(bts)
	succ := BuildSuccinctIndexedFileBufferFromInput(&input, pos, &util.SuccinctConf{
		SaSamplingRate:  int32(32),
		IsaSamplingRate: int32(32),
		NpaSamplingRate: int32(128),
	})
	outputPath := "./output/code.succinct"
	of, _ := os.Create(outputPath)

	succ.WriteToFile(of)
	of.Close()
	of, _ = os.Open(outputPath)
	succ = ReadSuccinctIndexFileBufferFromFile(of)
	fmt.Println(succ.Record(0))
}
