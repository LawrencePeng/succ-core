package parser

type Regex interface {
	Type() int
}

const (
	MGRAM = iota
	CHAR_RANGE
	DOT

	BLANK_T
	CONCAT
	PRIMITIVE
	REPEAT
	UNION
	WILDCARD_T

	ZERO_OR_MORE
	ONE_OR_MORE
	MIN_TO_MAX
)

type RegexBlank struct {
}

func (b RegexBlank) Type() int {
	return BLANK_T
}

type RegexConcat struct {
	Left, Right Regex
}

func (c RegexConcat) Type() int {
	return CONCAT
}

type RegexPrimitive struct {
	PrimitiveStr string
	PrimitiveType int
}

func (p RegexPrimitive) Type() int {
	return PRIMITIVE
}

type RegexRepeat struct {
	RepeatType int
	Internal Regex
	Min int
	Max int
}

func (r RegexRepeat) Type() int {
	return REPEAT
}

type RegexUnion struct {
	First Regex
	Second Regex
}

func (u RegexUnion) Type() int {
	return UNION
}

type RegexWildcard struct {
	Left Regex
	Right Regex
}

func (w RegexWildcard) Type() int {
	return WILDCARD_T
}