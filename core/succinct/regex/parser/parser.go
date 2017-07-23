package parser

import (
	"strings"
	"math"
)

const (
	WILDCARD = '@'
	REVERSED = "(){}[]|+*@."
)
var BLANK = RegexBlank{}

type RegexParser struct {
	Exp string
}

func NewRegexParser(exp string) *RegexParser {
	exp = strings.Replace(exp, ".*", string(WILDCARD), math.MaxInt16)

	for ;exp != "" && exp[0] == WILDCARD; {
		exp = exp[1:]
	}

	for exp != "" && exp[len(exp) - 1] == WILDCARD {
		exp = exp[:len(exp) - 1]
	}

	return &RegexParser{
		Exp: exp,
	}
}

func (parser *RegexParser) Parse() Regex {
	r := parser.regex()
	parser.checkAST(r, nil)
	return r
}
func (parser *RegexParser) checkAST(node Regex, parent Regex) {
	switch node.Type() {
	case BLANK_T:
		break

	case PRIMITIVE:
		break

	case WILDCARD_T:
		w := node.(RegexWildcard)
		if parent != nil && parent.Type() != WILDCARD_T {
			panic("Wildcard node has non-wildcard parent.")
		}
		parser.checkAST(w.Left, w)
		parser.checkAST(w.Right, w)
		break

	case REPEAT:
		r := node.(RegexRepeat)
		parser.checkAST(r.Internal, r)
		break

	case CONCAT:
		c := node.(RegexConcat)
		parser.checkAST(c.Left, c)
		parser.checkAST(c.Right, c)
		break

	case UNION:
		u := node.(RegexUnion)
		parser.checkAST(u.First, u)
		parser.checkAST(u.Second, u)
		break
	}
}
func (parser *RegexParser) regex() Regex {
	t := parser.term()
	if parser.more() && parser.peek() == '|' {
		parser.eat('|')
		r := parser.regex()
		return RegexUnion{t, r}
	}
	return t
}

func (parser *RegexParser) term() Regex {
	var f Regex
	f = BLANK

	for ;parser.more() && parser.peek() != ')' && parser.peek() != '|'; {
		if parser.peek() == WILDCARD {
			parser.eat(WILDCARD)
			nextF := parser.factor()
			if f.Type() == BLANK_T || nextF.Type() == BLANK_T {
				panic("evil input: empty children for wildcard op")
			}
			f = RegexWildcard{f, nextF}
		} else {
			nextF := parser.factor()
			f = parser.concat(f, nextF)
		}
	}

	return f
}

func (parser *RegexParser) concat(a Regex, b Regex) Regex {
	if a.Type() == BLANK_T {
		return b
	} else if a.Type() == PRIMITIVE && b.Type() == PRIMITIVE {
		primitiveA := a.(RegexPrimitive)
		primitiveB := b.(RegexPrimitive)

		if primitiveA.PrimitiveType == MGRAM &&
			primitiveB.PrimitiveType == MGRAM {
			return RegexPrimitive{primitiveA.PrimitiveStr + primitiveB.PrimitiveStr, MGRAM}
		}
	}
	return RegexConcat{a, b}
}

func (parser *RegexParser) factor() Regex  {
	b := parser.base()

	if parser.more() && parser.peek() == '*' {
		parser.eat('*')
		b = RegexRepeat{ZERO_OR_MORE,b, 	0, math.MaxInt16}
	} else if parser.more() && parser.peek() == '+' {
		parser.eat('+')
		b = RegexRepeat{ONE_OR_MORE,b, 1, math.MaxInt16}
	} else if parser.more() && parser.peek() == '{' {
		parser.eat('{')
		min := parser.nextInt()
		parser.eat(',')
		max := parser.nextInt()
		parser.eat('}')
		b = RegexRepeat{MIN_TO_MAX,b, min, max}
	}

	return b
}

func (parser *RegexParser) nextInt() int {
	num := 0
	for ; parser.peek() >= 48 && parser.peek() <= 57;  {
		num = num * 10 + int(parser.next()[0]) - 48
	}
	return num
}

func (parser *RegexParser) base() Regex  {
	if parser.more() && parser.peek() == '(' {
		parser.eat('(')
		r := parser.regex()
		parser.eat(')')
		return r
	}
	return parser.primitive()
}
func (parser *RegexParser) primitive() Regex  {
	if parser.more() && parser.peek() == '[' {
		parser.eat('[')
		charRange := ""
		for ;parser.peek() != ']'; {
			charRange += parser.next()
		}
		charRange = parser.expandCharRange(charRange)
		parser.eat(']')
		return RegexPrimitive{".", CHAR_RANGE}
	} else if parser.more() && parser.peek() == '.' {
		parser.eat('.')
		return RegexPrimitive{".", DOT}
	}

	m := ""
	for ; parser.more() &&
		!strings.Contains(REVERSED, string(parser.peek())); {
		m += parser.nextChar()
	}

	if m == "" {
		return BLANK
	}

	return RegexPrimitive{m, MGRAM}
}
func (parser *RegexParser) nextChar() string {
	if parser.peek() == '\\' {
		parser.eat('\\')
	}
	return parser.next()
}
func (parser *RegexParser) expandCharRange(charRange string) string {
	expandedCharRange := ""
	for i := 0; i < len(charRange); i++ {
		if charRange[i] == '-' {
			beg := charRange[i - 1]
			end := charRange[i + 1]

			for c := beg + 1; c < end; c++ {
				expandedCharRange += string(c)
			}
			i++
		}
		if charRange[i] == '\\' {
			i++
		}

		expandedCharRange += string(charRange[i])
	}

	return expandedCharRange
}
func (parser *RegexParser) next() string {
	c := parser.peek()
	parser.eat(c)
	return string(c)
}
func (parser *RegexParser) eat(c byte) {
	if parser.peek() == c {
		parser.Exp = parser.Exp[1:]
	} else {
		panic("Could not parse regex")
	}
}

func (parser *RegexParser) peek() byte {
	return parser.Exp[0]
}
func (parser *RegexParser) more() bool {
	return len(parser.Exp) > 0
}
