package succinct_test

import "testing"
import (
	. "."
	"os"
	"./util"
)

var filePath string = "./resources/test_file"

func TestReadSuccinctBufferFromFile(t *testing.T) {
	f, err := os.Open(filePath)
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

	succ := BuildSuccinctBufferFromInput(&SuccinctSource{Bts:bts},
		&util.SuccinctConf{
			SaSamplingRate: int32(32),
			IsaSamplingRate: int32(32),
			NpaSamplingRate: int32(128),
	})
}


