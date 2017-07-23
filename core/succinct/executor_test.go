package succinct_test

import "github.com/lexandro/go-assert"
import (
	"testing"
	"./regex/parser"
	"./regex"
	"./util"
	"."
)

var (
	input = "YoHoYoHoHoYoYoHoHoHo"
	succinctFile = succinct.BuildSuccinctFileBufferFromInput(input, &util.SuccinctConf{
		32,32,128,
	})
	re = "YO"
)

func TestGreedy(t *testing.T) {
	fre := regex.SuccinctFwdRegexExecutor{
		Regexable: succinctFile,
		Regex: parser.NewRegexParser(re).Parse(),
		Greedy: true,
		Alphabet: succinctFile.Alphabet(),
	}

	set := fre.Execute()
	assert.IsTrue(t, set.Set.Size() == 4)
}


