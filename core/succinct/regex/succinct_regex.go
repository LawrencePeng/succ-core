package regex

import "./parser"
import (
	"../util"
)
func IsPrefixed(r parser.Regex) bool {
	switch r.Type() {
	case parser.BLANK_T:
		return false

	case parser.PRIMITIVE:
		return r.(parser.RegexPrimitive).PrimitiveType == parser.MGRAM

	case parser.REPEAT:
		return IsPrefixed(r.(parser.RegexRepeat).Internal)

	case parser.CONCAT:
		return IsPrefixed(r.(parser.RegexConcat).Left)

	case parser.WILDCARD_T:
		return IsPrefixed(r.(parser.RegexWildcard).Left)

	case parser.UNION:
		return IsPrefixed(r.(parser.RegexUnion).First) && IsPrefixed(r.(parser.RegexUnion).Second)

	default:
		return false
	}
}

func IsSuffixed(r parser.Regex) bool {
	switch r.Type() {
	case parser.BLANK_T:
		return false

	case parser.PRIMITIVE:
		return r.(parser.RegexPrimitive).PrimitiveType == parser.MGRAM

	case parser.REPEAT:
		return IsPrefixed(r.(parser.RegexRepeat).Internal)

	case parser.CONCAT:
		return IsPrefixed(r.(parser.RegexConcat).Right)

	case parser.WILDCARD_T:
		return IsPrefixed(r.(parser.RegexWildcard).Right)

	case parser.UNION:
		return IsPrefixed(r.(parser.RegexUnion).First) && IsPrefixed(r.(parser.RegexUnion).Second)

	default:
		return false
	}
}

type Regexable interface {
	SameRecord(int32, int32) bool
	SuccinctIndexToOffset(int64) int64
	RangeToOffsets(*util.Range) []int64
	Alphabet() []int32
	BwdSearchStr(string) *util.Range
	FwdSearchStr(string) *util.Range
	ContinueBwdSearchStr(string, *util.Range) *util.Range
	ContinueFwdSearchWithQuery(string, *util.Range, int32) *util.Range
}

func Regex(succ Regexable, q string, greedy bool) *TypedTreeSet {
	regex := parser.NewRegexParser(q).Parse()
	ex := SuccinctFwdRegexExecutor{succ, regex, greedy, succ.Alphabet() }
	return ex.Execute()
}