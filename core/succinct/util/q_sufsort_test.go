package util_test

import (
	. "."
	"os"
	"testing"
	"sort"
	"bytes"
	. ".."
)

const testFileRaw = "../resources/test_file"
const testFileSA = "../resources/test_file.sa"
const testFileIPA = "../resources/test_file.isa"

var suf_q *QSufSort =&QSufSort{}
var bts []byte


func init() {
	rawFile, err := os.Open(testFileRaw)
	if err != nil {
		panic("fail to open rawfile")
	}

	stat, _ := rawFile.Stat()
	sz := stat.Size()
	bts = make([]byte, sz)
	rawFile.Read(bts)

	source := &SuccinctSource{
		Bts:bts,
	}

	suf_q.BuildSuffixArray(source)
}

func TestGetAlphabet(t *testing.T) {
	alphabet := suf_q.Alphabet

	set := HashSet{M: make(map[int32]bool)}

	for _, bt := range bts {
		set.Add(int32(bt))
	}
	set.Add(EOF)
	exp := make([]int32, set.Len())
	i := int32(0)
	for k := range set.M {
		exp[i] = k
		i++
	}

	tr := make([]int, set.Len())
	for i := 0; i < set.Len(); i++ {
		tr[i] = int(exp[i])
	}
	sort.Ints(tr)
	for i := 0; i < set.Len(); i++ {
		exp[i] = int32(tr[i])
	}

	for i := 0; i < len(alphabet); i++ {
		if alphabet[i] != exp[i] {
			panic(alphabet[i])
		}
	}
}

func TestSA(t *testing.T) {
	sa := suf_q.SA()

	saFile, _ := os.Open(testFileSA)

	stat, _ := saFile.Stat()
	bts := make([]byte, stat.Size())

	saFile.Read(bts)
	testSa := ReadArray(bytes.NewBuffer(bts))

	for i := int64(0); i < stat.Size(); i++ {
		if testSa[i] != sa[i] {
			panic(i)
		}
	}

}