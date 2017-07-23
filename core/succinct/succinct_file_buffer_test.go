package succinct_test

import (
	. "."
	"./util"
	"fmt"
	"os"
	"testing"
)

func TestCharAt(t *testing.T) {
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

	succ := BuildSuccinctFileBufferFromInput(string(bts), &util.SuccinctConf{
		SaSamplingRate:  int32(32),
		IsaSamplingRate: int32(32),
		NpaSamplingRate: int32(128),
	})
	fmt.Println(succ.CharAt(0))
	fmt.Println(succ.Extract(0, 100))
	fmt.Println(succ.ExtractUntil(0, int32('\n')))
	fmt.Println(succ.Count(&SuccinctSource{Bts: []byte("int")}))
	fmt.Println(len(succ.Search(&SuccinctSource{Bts: []byte("int")})))
}
