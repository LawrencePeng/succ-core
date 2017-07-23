package parser_test

import (
	"testing"
	"github.com/lexandro/go-assert"
	. "."
)

func TestParse(t *testing.T) {
	mgramParser1 := NewRegexParser("abcd")
	mgramRegEx1 := mgramParser1.Parse()
	assert.Equals(t, mgramRegEx1.Type(), PRIMITIVE)
	assert.Equals(t,mgramRegEx1.(RegexPrimitive).PrimitiveStr, "abcd")

	mgramParser2 := NewRegexParser("\\|a\\(b\\)c\\*d\\+\\{\\}")
	mframRegEx2 := mgramParser2.Parse()
	assert.Equals(t, mframRegEx2.Type(), PRIMITIVE)
	assert.Equals(t, mframRegEx2.(RegexPrimitive).PrimitiveStr, "|a(b)c*d+{}")

	unionParser := NewRegexParser("a|b")
	unionRegex := unionParser.Parse()
	assert.Equals(t, unionRegex.Type(), UNION)
	uRE := unionRegex.(RegexUnion)
	assert.Equals(t, uRE.First.Type(), PRIMITIVE)
	assert.Equals(t, uRE.First.(RegexPrimitive).PrimitiveStr, "a")
	assert.Equals(t, uRE.Second.Type(), PRIMITIVE)
	assert.Equals(t, uRE.Second.(RegexPrimitive).PrimitiveStr, "b")

	concatParser1 := NewRegexParser("a(b|c)")
	concatRegex1 := concatParser1.Parse()
	assert.Equals(t, concatRegex1.Type(), CONCAT)
	cRE1 := concatRegex1.(RegexConcat)
	assert.Equals(t, cRE1.Left.Type(), PRIMITIVE)
	assert.Equals(t, cRE1.Left.(RegexPrimitive).PrimitiveStr, "a")
	assert.Equals(t, cRE1.Right.Type(), UNION)

	concatParser2 := NewRegexParser("a(b)")
	concatRegex2 := concatParser2.Parse()
	assert.Equals(t, concatRegex2.Type(), PRIMITIVE)
	assert.Equals(t, concatRegex2.(RegexPrimitive).PrimitiveStr, "ab")

	repeatParser1 := NewRegexParser("a*")
	repeatRegex1 := repeatParser1.Parse()
	assert.Equals(t, repeatRegex1.Type(), REPEAT)
	rRE1 := repeatRegex1.(RegexRepeat)
	assert.Equals(t, rRE1.RepeatType, ZERO_OR_MORE)
	assert.Equals(t, rRE1.Internal.Type(), PRIMITIVE)
	assert.Equals(t, rRE1.Internal.Type(), PRIMITIVE)
	assert.Equals(t, rRE1.Internal.(RegexPrimitive).PrimitiveStr, "a")

	repeatParser2 := NewRegexParser("a+")
	repeatRegex2 := repeatParser2.Parse()
	assert.Equals(t, repeatRegex2.Type(), REPEAT)
	rRE2 := repeatRegex2.(RegexRepeat)
	assert.Equals(t, rRE2.RepeatType, ONE_OR_MORE)
	assert.Equals(t, rRE2.Internal.Type(), PRIMITIVE)
	assert.Equals(t, rRE2.Internal.Type(), PRIMITIVE)
	assert.Equals(t, rRE2.Internal.(RegexPrimitive).PrimitiveStr, "a")

	repeatParser3 := NewRegexParser("a{1,2}")
	repeatRegex3 := repeatParser3.Parse()
	assert.Equals(t, repeatRegex3.Type(), REPEAT)
	rRE3 := repeatRegex3.(RegexRepeat)
	assert.Equals(t, rRE3.RepeatType, MIN_TO_MAX)
	assert.Equals(t, rRE3.Internal.Type(), PRIMITIVE)
	assert.Equals(t, rRE3.Internal.Type(), PRIMITIVE)
	assert.Equals(t, rRE3.Internal.(RegexPrimitive).PrimitiveStr, "a")
	assert.Equals(t, rRE3.Min, 1)
	assert.Equals(t, rRE3.Max, 2)

	wildcardParser1 := NewRegexParser(".*abc.*")
	wildcardRegex1 := wildcardParser1.Parse()
	assert.Equals(t, PRIMITIVE, wildcardRegex1.Type())
	pRE1 := wildcardRegex1.(RegexPrimitive)
	assert.Equals(t, "abc", pRE1.PrimitiveStr)
	assert.Equals(t, MGRAM, pRE1.PrimitiveType)

	wildcardParser2 := NewRegexParser("a.*b")
	wildcardRegex2 := wildcardParser2.Parse()
	assert.Equals(t, WILDCARD_T, wildcardRegex2.Type())
	wRE1 := wildcardRegex2.(RegexWildcard)
	assert.Equals(t, PRIMITIVE, wRE1.Left.Type())
	assert.Equals(t, "a", wRE1.Left.(RegexPrimitive).PrimitiveStr)
	assert.Equals(t, PRIMITIVE, wRE1.Right.Type())
	assert.Equals(t, "b", wRE1.Right.(RegexPrimitive).PrimitiveStr)

}
