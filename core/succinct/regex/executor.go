package regex

import "github.com/emirpasic/gods/sets/treeset"
import "./parser"
import (
	"../util"
	"math"
)

const (
	FRONT_SORTED = iota
	END_SORTED
)

type Executor interface {
	Execute() *TypedTreeSet
}

type SuccinctBwdRegexExecutor struct {
	Regexable Regexable
	Regex  parser.Regex
	Greedy bool
	Alphabet []int32
}

type SuccinctFwdRegexExecutor struct {
	Regexable Regexable
	Regex  parser.Regex
	Greedy bool
	Alphabet []int32
}

type TypedTreeSet struct {
	Set *treeset.Set
	sortType int
}

func (fre *SuccinctFwdRegexExecutor) allocateSet(sortType int) *TypedTreeSet {
	if sortType == END_SORTED {
		return &TypedTreeSet{
			Set:      treeset.NewWith(END_COMP),
			sortType: END_SORTED,
		}
	}
	return &TypedTreeSet{
		Set:      treeset.NewWith(FRONT_COMP),
		sortType: FRONT_SORTED,
	}
}

func (fre *SuccinctFwdRegexExecutor) regexWildcard(left, right *TypedTreeSet, sortType int) *TypedTreeSet {
	wildcardRes := fre.allocateSet(sortType)

	lowerBoundEntry := RegexMatch{
		Offset: 0,
		Length: 0,
	}

	leftIter := left.Set.Iterator()
	for leftIter.Next() {
		leftEntry := leftIter.Value().(RegexMatch)
		lowerBoundEntry.Offset = leftEntry.End()
		rightEntry := ceiling(right, lowerBoundEntry)

		if rightEntry == nil {
			break
		}

		var lastMatch *RegexMatch
		for ; rightEntry != nil && fre.Regexable.SameRecord(int32(leftEntry.Offset), int32(rightEntry.Offset));  {
			lastMatch = rightEntry
			rightEntry = higher(right, *rightEntry)
		}

		if lastMatch != nil {
			dist := lastMatch.Offset - leftEntry.Offset
			wildcardRes.Set.Add(RegexMatch{
				Offset: leftEntry.Offset,
				Length: dist + lastMatch.Length,
			})

			for ; leftIter.Next() && leftEntry.Offset < lastMatch.Offset; {
				leftEntry = leftIter.Value().(RegexMatch)
			}
		}

	}
	return wildcardRes
}

func (fre *SuccinctFwdRegexExecutor) rangeResultsToRegexMatches(rangeRes map[SuccinctRegexMatch]bool, sortType int) *TypedTreeSet {
	regexMatches := fre.allocateSet(sortType)
	if fre.Greedy {
		succinctMatches := make(map[int64]*int)
		for match := range rangeRes {
			for i := match.From; i <= match.To; i++ {
				length := succinctMatches[i]
				if length == nil || int32(*length) < match.Length {
					length := int(match.Length)
					succinctMatches[i] = &length
				}
			}
		}

		for k, v := range succinctMatches {
			regexMatches.Set.Add(RegexMatch{
				int(fre.Regexable.SuccinctIndexToOffset(k)),
				*v,
			})
		}

		var prv *RegexMatch
		it := regexMatches.Set.Iterator()
		for ; it.Next(); {
			cur := it.Value()
			if prv != nil && prv.Contains(cur.(RegexMatch)) {
				regexMatches.Set.Remove(cur)
				it = regexMatches.Set.Iterator()
				for ; it.Value() != prv; it.Next() {}
				curPtr := cur.(RegexMatch)
				prv = &curPtr
			}
		}
	} else {
		for match := range rangeRes {
			if match.Empty() {
				offsets := fre.Regexable.RangeToOffsets(&match.Range)
				for _, offset := range offsets {
					regexMatches.Set.Add(RegexMatch{
						Offset: int(offset),
						Length: int(match.Length),
					})
				}
			}
		}
	}

	return regexMatches

}

func (bre *SuccinctBwdRegexExecutor) allocateSet(sortType int) *TypedTreeSet {
	if sortType == END_SORTED {
		return &TypedTreeSet{
			Set:      treeset.NewWith(END_COMP),
			sortType: END_SORTED,
		}
	}
	return &TypedTreeSet{
		Set:      treeset.NewWith(FRONT_COMP),
		sortType: FRONT_SORTED,
	}
}

func (bre *SuccinctBwdRegexExecutor) regexWildcard(left, right *TypedTreeSet, sortType int) *TypedTreeSet {
	wildcardRes := bre.allocateSet(sortType)

	lowerBoundEntry := RegexMatch{
		Offset: 0,
		Length: 0,
	}

	leftIter := left.Set.Iterator()
	for leftIter.Next() {
		leftEntry := leftIter.Value().(RegexMatch)
		lowerBoundEntry.Offset = leftEntry.End()
		rightEntry := ceiling(right, lowerBoundEntry)

		if rightEntry == nil {
			break
		}

		var lastMatch *RegexMatch
		for ; rightEntry != nil && bre.Regexable.SameRecord(int32(leftEntry.Offset), int32(rightEntry.Offset));  {
			lastMatch = rightEntry
			rightEntry = higher(right, *rightEntry)
		}

		if lastMatch != nil {
			dist := lastMatch.Offset - leftEntry.Offset
			wildcardRes.Set.Add(RegexMatch{
				Offset: leftEntry.Offset,
				Length: dist + lastMatch.Length,
			})

			for ; leftIter.Next() && leftEntry.Offset < lastMatch.Offset; {
				leftEntry = leftIter.Value().(RegexMatch)
			}
		}

	}
	return wildcardRes
}

func (fre *SuccinctFwdRegexExecutor) Execute() *TypedTreeSet {
	return fre.compute(fre.Regex, FRONT_SORTED)
}

func (fre *SuccinctFwdRegexExecutor) compute(r parser.Regex, sortType int) *TypedTreeSet {
	var results *TypedTreeSet
	switch r.Type() {
	case parser.WILDCARD_T:
		w := r.(parser.RegexWildcard)
		leftResults := fre.compute(w.Left, END_SORTED)
		rightResults := fre.compute(w.Right, FRONT_SORTED)
		results = fre.regexWildcard(leftResults, rightResults, sortType)
	default:
		succinctResults := fre.computeSuccinctly(r)
		results = fre.rangeResultsToRegexMatches(succinctResults, sortType)
	}
	return results
}

func (fre *SuccinctFwdRegexExecutor) computeSuccinctly(r parser.Regex) map[SuccinctRegexMatch]bool {
	results := make(map[SuccinctRegexMatch]bool)
	switch r.Type() {
	case parser.BLANK_T:
	case parser.PRIMITIVE:
		p := r.(parser.RegexPrimitive)
		switch p.PrimitiveType {
		case parser.MGRAM:
			mgram := p.PrimitiveStr
			ran := fre.Regexable.FwdSearchStr(mgram)
			if !ran.Empty() {
				results[SuccinctRegexMatch{
					Range: *ran,
					Length: int32(len(mgram)),
				}] = true
			}
		case parser.DOT:
			for _, b := range fre.Alphabet {
				if b == util.EOL || b == util.EOF {
					continue
				}
				ran := fre.Regexable.FwdSearchStr(string(b))
				if !ran.Empty() {
					results[SuccinctRegexMatch{
						Range: *ran,
						Length: 1,
					}] = true
				}
			}
		case parser.CHAR_RANGE:
			charRange := p.PrimitiveStr
			for _, b := range charRange {
				ran := fre.Regexable.FwdSearchStr(string(b))
				if !ran.Empty() {
					results[SuccinctRegexMatch{
						Range: *ran,
						Length: 1,
					}] = true
				}
			}
		}
	case parser.UNION:
		u := r.(parser.RegexUnion)
		fir := fre.computeSuccinctly(u.First)
		for k := range fir {
			results[k] = true
		}
		sec := fre.computeSuccinctly(u.Second)
		for k := range sec {
			results[k] = true
		}

	case parser.CONCAT:
		c := r.(parser.RegexConcat)
		leftResults := fre.computeSuccinctly(c.Left)
		for leftMatch := range leftResults {
			inter := fre.regexConcat(c.Right, leftMatch)
			for k := range inter {
				results[k] = true
			}
		}

	case parser.REPEAT:
		rep := r.(parser.RegexRepeat)
		switch rep.RepeatType {
		case parser.ZERO_OR_MORE:
			results = fre.regexRepeatOneOrMore(rep.Internal)
		case parser.ONE_OR_MORE:
			results = fre.regexRepeatOneOrMore(rep.Internal)
		case parser.MIN_TO_MAX:
			results = fre.regexRepeatMinToMax(rep.Internal, rep.Min, rep.Max)
		}

	default:
		panic("Invalid node in succinct regex parse tree.")
	}
	return results
}

func (fre *SuccinctFwdRegexExecutor) regexRepeatZeroOrMore(r parser.Regex, leftMatch SuccinctRegexMatch) map[SuccinctRegexMatch]bool {
	repeatResults := make(map[SuccinctRegexMatch]bool)
	if leftMatch.Empty() {
		return repeatResults
	}

	repeatResults[leftMatch] = true
	inter := fre.regexRepeatOneOrMoreWithMatch(r, leftMatch)
	for k := range inter {
		repeatResults[k] = true
	}
	return repeatResults
}

func (fre *SuccinctFwdRegexExecutor) regexRepeatMinToMax(r parser.Regex, min, max int) map[SuccinctRegexMatch]bool {
	if min > 0 {
		min--
	} else {
		min = 0
	}

	if max > 0 {
		max--
	} else {
		max = 0
	}

	repeatResults := make(map[SuccinctRegexMatch]bool)
	internalResults := fre.computeSuccinctly(r)
	if len(internalResults) == 0 {
		return repeatResults
	}

	if min == 0 {
		for k := range internalResults {
			repeatResults[k] = true
		}
	}

	if max > 0 {
		for internalMatch := range internalResults {
			inter := fre.regexRepeatMinToMaxWithMatch(r, internalMatch, min, max)
			for k := range inter {
				repeatResults[k] = true
			}
		}
	}

	return repeatResults
}

func (fre *SuccinctFwdRegexExecutor) regexRepeatMinToMaxWithMatch(r parser.Regex, leftMatch SuccinctRegexMatch, min, max int) map[SuccinctRegexMatch]bool {
	if min > 0 {
		min--
	} else {
		min = 0
	}

	if max > 0 {
		max--
	} else {
		max = 0
	}

	repeatResults := make(map[SuccinctRegexMatch]bool)
	if leftMatch.Empty() {
		return repeatResults
	}

	concatResults := fre.regexConcat(r, leftMatch)
	if len(concatResults) == 0 {
		return repeatResults
	}

	if min == 0 {
		return repeatResults
	}

	if max > 0 {
		for concatMatch := range concatResults {
			inter := fre.regexRepeatMinToMaxWithMatch(r, concatMatch, min, max)
			for k := range inter {
				repeatResults[k] = true
			}
		}
	}

	return repeatResults

}

func (fre *SuccinctFwdRegexExecutor) regexRepeatOneOrMore(r parser.Regex) map[SuccinctRegexMatch]bool  {
	repeatResults := make(map[SuccinctRegexMatch]bool)
	internalResults := fre.computeSuccinctly(r)

	if len(internalResults) == 0 {
		return repeatResults
	}

	for k := range internalResults {
		repeatResults[k] = true
	}

	for internalMatch := range internalResults {
		inter := fre.regexRepeatOneOrMoreWithMatch(r, internalMatch)
		for k := range inter {
			repeatResults[k] = true
		}
	}

	return repeatResults
}

func (fre *SuccinctFwdRegexExecutor) regexRepeatOneOrMoreWithMatch(r parser.Regex, leftMatch SuccinctRegexMatch) map[SuccinctRegexMatch]bool {
	repeatResults := make(map[SuccinctRegexMatch]bool)
	internalResults := fre.computeSuccinctly(r)

	if len(internalResults) == 0 {
		return repeatResults
	}

	for k := range internalResults {
		repeatResults[k] = true
	}

	for internalMatch := range internalResults {
		inter := fre.regexRepeatOneOrMoreWithMatch(r, internalMatch)
		for k := range inter {
			repeatResults[k] = true
		}
	}

	return repeatResults
}

func (fre *SuccinctFwdRegexExecutor) regexConcat(r parser.Regex, leftMatch SuccinctRegexMatch) map[SuccinctRegexMatch]bool  {
	concatResults := make(map[SuccinctRegexMatch]bool)

	if leftMatch.Empty() {
		return concatResults
	}

	switch r.Type() {
	case parser.BLANK_T:
	case parser.PRIMITIVE:
		p := r.(parser.RegexPrimitive)
		switch p.PrimitiveType {
		case parser.MGRAM:
			mgram := p.PrimitiveStr
			ran := fre.Regexable.ContinueFwdSearchWithQuery(mgram, &leftMatch.Range, int32(leftMatch.Size()))
			if !ran.Empty() {
				m := SuccinctRegexMatch{
					Range: *ran,
					Length: leftMatch.Length + int32(len(mgram)),
				}
				concatResults[m] = true
			}
		case parser.DOT:
			for _, b := range fre.Alphabet {
				if b == util.EOL || b == util.EOF {
					continue
				}

				ran := fre.Regexable.ContinueFwdSearchWithQuery(string(b), &leftMatch.Range, leftMatch.Length + 1)
				if !ran.Empty() {
					m := SuccinctRegexMatch{
						Range: *ran,
						Length: leftMatch.Length + 1,
					}
					concatResults[m] = true
				}
			}
		case parser.CHAR_RANGE:
			charRange := p.PrimitiveStr
			for _, b := range charRange {
				ran := fre.Regexable.ContinueFwdSearchWithQuery(string(b), &leftMatch.Range, leftMatch.Length)
				if !ran.Empty() {
					m := SuccinctRegexMatch{
						Range: *ran,
						Length: leftMatch.Length + 1,
					}
					concatResults[m] = true
				}
			}
		}
	case parser.UNION:
		u := r.(parser.RegexUnion)
		fir := fre.regexConcat(u.First, leftMatch)
		for k := range fir {
			concatResults[k] = true
		}
		sec := fre.regexConcat(u.Second, leftMatch)
		for k := range sec {
			concatResults[k] = true
		}

	case parser.CONCAT:
		c := r.(parser.RegexConcat)
		leftOfRightResults := fre.regexConcat(c.Left, leftMatch)
		for leftOfRightMatch := range leftOfRightResults {
			inter := fre.regexConcat(c.Right, leftOfRightMatch)
			for k := range inter {
				concatResults[k] = true
			}
		}

	case parser.REPEAT:
		rep := r.(parser.RegexRepeat)
		switch rep.RepeatType {
		case parser.ZERO_OR_MORE:
			concatResults = fre.regexRepeatZeroOrMore(rep.Internal, leftMatch)
		case parser.ONE_OR_MORE:
			concatResults = fre.regexRepeatOneOrMoreWithMatch(rep.Internal, leftMatch)
		case parser.MIN_TO_MAX:
			concatResults = fre.regexRepeatMinToMaxWithMatch(rep.Internal, leftMatch, rep.Min, rep.Max)
		}

	default:
		panic("Invalid node in succinct regex parse tree.")
	}

	return concatResults
}

func higher(TSet *TypedTreeSet, match RegexMatch) *RegexMatch {
	vals := TSet.Set.Values()
	length := len(vals)
	if length == 0 {
		return nil
	}

	var comp func(interface{}, interface{}) int
	if TSet.sortType == END_SORTED {
		comp = END_COMP
	} else {
		comp = FRONT_COMP
	}

	low := 0
	high := length - 1
	for ; high > low; {
		mid := (low + high) / 2
		if comp(vals[mid], mid) >= 0 {
			high = mid
		} else {
			low = mid - 1
		}
	}

	if low > high {
		return nil
	}

	if comp(vals[low], match) == 0 && low < length - 1{
		low++
	}
	ret := vals[low].(RegexMatch)
	return &ret

}

func ceiling(TSet *TypedTreeSet, match RegexMatch) *RegexMatch {
	vals := TSet.Set.Values()
	length := len(vals)
	if length == 0 {
		return nil
	}
	
	var comp func(interface{}, interface{}) int
	if TSet.sortType == END_SORTED {
		comp = END_COMP
	} else {
		comp = FRONT_COMP
	}
	
	low := 0
	high := length - 1
	for ; high > low; {
		mid := (low + high) / 2
		if comp(vals[mid], match) >= 0 {
			high = mid
		} else {
			low = mid - 1
		}
	}

	if low > high {
		return nil
	}
	ret := vals[low].(RegexMatch)
	return &ret
}

func (bre *SuccinctBwdRegexExecutor) rangeResultsToRegexMatches(rangeRes map[SuccinctRegexMatch]bool, sortType int) *TypedTreeSet {
	regexMatches := bre.allocateSet(sortType)
	if bre.Greedy {
		succinctMatches := make(map[int64]*int)
		for match := range rangeRes {
			for i := match.From; i <= match.To; i++ {
				length := succinctMatches[i]
				if length == nil || int32(*length) < match.Length {
					length := int(match.Length)
					succinctMatches[i] = &length
				}
			}
		}

		for k, v := range succinctMatches {
			regexMatches.Set.Add(RegexMatch{
				int(bre.Regexable.SuccinctIndexToOffset(k)),
				*v,
			})
		}

		var prv *RegexMatch
		it := regexMatches.Set.Iterator()
		for ; it.Next(); {
			cur := it.Value()
			if prv != nil && prv.Contains(cur.(RegexMatch)) {
				regexMatches.Set.Remove(cur)
				it = regexMatches.Set.Iterator()
				for ; it.Value() != prv; it.Next() {}
				prvPtr := cur.(RegexMatch)
				prv = &prvPtr
			}
		}
	} else {
		for match := range rangeRes {
			if match.Empty() {
				offsets := bre.Regexable.RangeToOffsets(&match.Range)
				for _, offset := range offsets {
					regexMatches.Set.Add(RegexMatch{
						Offset: int(offset),
						Length: int(match.Length),
					})
				}
			}
		}
	}

	return regexMatches

}

func (bre *SuccinctBwdRegexExecutor) Execute() *TypedTreeSet {
	return bre.compute(bre.Regex, FRONT_SORTED)
}

func (bre *SuccinctBwdRegexExecutor) compute(r parser.Regex, sortType int) *TypedTreeSet {
	var results *TypedTreeSet
	switch r.Type() {
	case parser.WILDCARD_T:
		w := r.(parser.RegexWildcard)
		leftResults := bre.compute(w.Left, END_SORTED)
		rightResults := bre.compute(w.Right, FRONT_SORTED)
		results = bre.regexWildcard(leftResults, rightResults, sortType)
	default:
		succinctResults := bre.computeSuccinctly(r)
		results = bre.rangeResultsToRegexMatches(succinctResults, sortType)
	}
	return results
}

func (bre *SuccinctBwdRegexExecutor) computeSuccinctly(r parser.Regex) map[SuccinctRegexMatch]bool {
	results := make(map[SuccinctRegexMatch]bool)
	switch r.Type() {
	case parser.BLANK_T:
	case parser.PRIMITIVE:
		p := r.(parser.RegexPrimitive)
		switch p.PrimitiveType {
		case parser.MGRAM:
			mgram := p.PrimitiveStr
			ran := bre.Regexable.BwdSearchStr(mgram)
			if !ran.Empty() {
				results[SuccinctRegexMatch{
					Range:  *ran,
					Length: int32(len(mgram)),
				}] = true
			}
		case parser.DOT:
			for _, b := range bre.Alphabet {
				if b == util.EOL || b == util.EOF {
					continue
				}
				ran := bre.Regexable.BwdSearchStr(string(b))
				if !ran.Empty() {
					results[SuccinctRegexMatch{*ran, int32(1)}] = true
				}
			}
		case parser.CHAR_RANGE:
			str := p.PrimitiveStr
			for _, c := range str {
				ran := *bre.Regexable.BwdSearchStr(string(c))
				if !ran.Empty() {
					results[SuccinctRegexMatch{
						Range:  ran,
						Length: int32(1),
					}] = true
				}
			}
		}

	case parser.UNION:
		u := r.(parser.RegexUnion)
		fir := bre.computeSuccinctly(u.First)
		for k := range fir {
			results[k] = true
		}
		sec := bre.computeSuccinctly(u.Second)
		for k := range sec {
			results[k] = true
		}

	case parser.CONCAT:
		c := r.(parser.RegexConcat)
		rightResults := bre.computeSuccinctly(c.Right)
		for rightMatch := range rightResults {
			inter := bre.regexConcat(c.Left, rightMatch)
			for k := range inter {
				inter[k] = true
			}
		}

	case parser.REPEAT:
		rep := r.(parser.RegexRepeat)

		switch rep.RepeatType {
		case parser.ZERO_OR_MORE:
			results = bre.regexRepeatOneOrMore(rep.Internal)
		case parser.ONE_OR_MORE:
			results = bre.regexRepeatOneOrMore(rep.Internal)
		case parser.MIN_TO_MAX:
			results = bre.regexRepeatMinToMax(rep.Internal, rep.Min, rep.Max)
		}
	}

	return results
}

func (bre *SuccinctBwdRegexExecutor) regexRepeatMinToMax(r parser.Regex, min, max int) map[SuccinctRegexMatch]bool {
	if min > 0 {
		min--
	} else {
		min = 0
	}

	if max > 0 {
		max--
	} else {
		max = 0
	}

	repeatResults := make(map[SuccinctRegexMatch]bool)

	internalResults := bre.computeSeedToRepeat(r)
	if len(internalResults) == 0 {
		return repeatResults
	}

	if min == 0 {
		for k := range internalResults {
			repeatResults[k] = true
		}
	}

	if max > 0 {
		for internalMatch := range internalResults {
			inter := bre.regexRepeatMinToMaxWithMatch(r, internalMatch, min, max)
			for k := range inter {
				repeatResults[k] = true
			}
		}
	}

	return repeatResults
}

func (bre *SuccinctBwdRegexExecutor) regexRepeatOneOrMore(r parser.Regex) map[SuccinctRegexMatch]bool {
	repeatResults := make(map[SuccinctRegexMatch]bool)
	internalResults := bre.computeSeedToRepeat(r)
	if len(internalResults) == 0 {
		return repeatResults
	}

	for k := range internalResults {
		repeatResults[k] = true
	}

	for internalMatch := range internalResults {
		inter := bre.regexRepeatOneOrMoreWithMatch(r, internalMatch)
		for k := range inter {
			repeatResults[k] = true
		}
	}

	return repeatResults
}

func (bre *SuccinctBwdRegexExecutor) computeSeedToRepeat(r parser.Regex) map[SuccinctRegexMatch]bool {
	results := bre.computeSuccinctly(r)

	if !bre.Greedy {
		return results
	}

	initRepeats := &TypedTreeSet{
		Set: treeset.NewWith(func(a interface{}, b interface{}) int {
			o1 := a.(SuccinctRegexMatch)
			o2 := b.(SuccinctRegexMatch)

			return int(o1.From - o2.From)
		}),
		sortType: 0,
	}

	for result := range results {
		inter := bre.regexConcat(r, result)
		for k := range inter {
			initRepeats.Set.Add(k)
		}
	}

	it := initRepeats.Set.Iterator()
	if it.Value() == nil {
		return results
	}
	it.Next()
	first := it.Value()
	start := first.(SuccinctRegexMatch).From
	end := first.(SuccinctRegexMatch).To

	for ; it.Next(); {
		current := it.Value().(SuccinctRegexMatch)
		if current.From <= end {
			end = int64(math.Max(float64(current.To), float64(end)))
		} else {
			// remove subrange
			newSubRange := make(map[SuccinctRegexMatch]bool)
			for match := range results {
				if match.Range.ContainsRange(&util.Range{
					From:start,
					To:end,
				}) {
					delete(results, match)
				}
				if match.From == start && match.To != end {
					newSubRange[SuccinctRegexMatch{
						Range: util.Range{
							From: end + 1,
							To: match.To,
						},
						Length: match.Length,
					}] = true
				} else if match.To == end && match.From != start {
					newSubRange[SuccinctRegexMatch{
						Range: util.Range{
							From: match.From,
							To: start - 1,
						},
						Length: match.Length,
					}] = true
				} else if match.From != start && match.To != end {
					newSubRange[SuccinctRegexMatch{
						Range: util.Range{
							From: match.From,
							To: start - 1,
						},
						Length: match.Length,
					}] = true
					newSubRange[SuccinctRegexMatch{
						Range: util.Range{
							From: end + 1,
							To: match.To,
						},
						Length: match.Length,
					}] = true
				}
			}

			for k := range newSubRange {
				results[k] = true
			}

			start = current.From
			end = current.To
		}
	}

	// rm subrange
	newSubRanges := make(map[SuccinctRegexMatch]bool)
	for match := range results {
		if match.Range.ContainsRange(&util.Range{
			From:start,
			To:end,
		}) {
			delete(results, match)
		}
		if match.From == start && match.To != end {
			newSubRanges[SuccinctRegexMatch{
				Range: util.Range{
					From: end + 1,
					To: match.To,
				},
				Length: match.Length,
			}] = true
		} else if match.To == end && match.From != start {
			newSubRanges[SuccinctRegexMatch{
				Range: util.Range{
					From: match.From,
					To: start - 1,
				},
				Length: match.Length,
			}] = true
		} else if match.From != start && match.To != end {
			newSubRanges[SuccinctRegexMatch{
				Range: util.Range{
					From: match.From,
					To: start - 1,
				},
				Length: match.Length,
			}] = true
			newSubRanges[SuccinctRegexMatch{
				Range: util.Range{
					From: end + 1,
					To: match.To,
				},
				Length: match.Length,
			}] = true
		}
	}

	for k := range newSubRanges {
		results[k] = true
	}

	return results
}

func (bre *SuccinctBwdRegexExecutor) regexRepeatOneOrMoreWithMatch(r parser.Regex, rightMatch SuccinctRegexMatch) map[SuccinctRegexMatch]bool {
	repeatResults := make(map[SuccinctRegexMatch]bool)
	if rightMatch.Empty() {
		return repeatResults
	}

	concatResults := bre.regexConcat(r, rightMatch)
	if len(concatResults) == 0 {
		return repeatResults
	}

	for k := range concatResults {
		repeatResults[k] = true
	}

	for concatMatch := range concatResults {
		inter := bre.regexRepeatOneOrMoreWithMatch(r, concatMatch)
		for k := range inter {
			repeatResults[k]  = true
		}
	}
	return repeatResults
}

func (bre *SuccinctBwdRegexExecutor) regexRepeatZeroOrMoreWithMatch(r parser.Regex, rightMatch SuccinctRegexMatch) map[SuccinctRegexMatch]bool {
	repeatResults := make(map[SuccinctRegexMatch]bool)
	if rightMatch.Empty() {
		return repeatResults
	}

	repeatResults[rightMatch] = true
	inter := bre.regexRepeatOneOrMoreWithMatch(r, rightMatch)
	for k := range inter {
		repeatResults[k] = true
	}

	return repeatResults
}

func (bre *SuccinctBwdRegexExecutor) regexRepeatMinToMaxWithMatch(r parser.Regex, rightMatch SuccinctRegexMatch, min, max int) map[SuccinctRegexMatch]bool {
	if min > 0 {
		min--
	} else {
		min = 0
	}

	if max > 0 {
		max--
	} else {
		max = 0
	}

	repeatResults := make(map[SuccinctRegexMatch]bool)
	if rightMatch.Empty() {
		return repeatResults
	}

	concatResults := bre.regexConcat(r, rightMatch)
	if len(concatResults) == 0 {
		return repeatResults
	}

	if min == 0 {
		for k := range concatResults {
			repeatResults[k] = true
		}
	}

	if max > 0 {
		for concatMatch := range concatResults {
			inter := bre.regexRepeatMinToMaxWithMatch(r, concatMatch, min, max)
			for k := range inter {
				repeatResults[k] = true
			}
		}
	}

	return repeatResults
}

func (bre *SuccinctBwdRegexExecutor) regexConcat(r parser.Regex, rightMatch SuccinctRegexMatch) map[SuccinctRegexMatch]bool {
	concatResults := make(map[SuccinctRegexMatch]bool)

	if rightMatch.Empty() {
		return concatResults
	}

	switch r.Type() {
	case parser.BLANK_T:
	case parser.PRIMITIVE:
		p := r.(parser.RegexPrimitive)
		switch p.PrimitiveType {
		case parser.MGRAM:
			mgram := p.PrimitiveStr
			ran := bre.Regexable.ContinueBwdSearchStr(mgram, &rightMatch.Range)
			if !ran.Empty() {
				concatResults[SuccinctRegexMatch{
					Range:	*ran,
					Length: rightMatch.Length + int32(len(mgram)),
				}] = true
			}
		case parser.DOT:
			for _, b := range bre.Alphabet {
				if b == util.EOL || b == util.EOF {
					continue
				}
				ran := bre.Regexable.ContinueBwdSearchStr(string(b), &rightMatch.Range)
				if !ran.Empty() {
					concatResults[SuccinctRegexMatch{
						Range:	*ran,
						Length: rightMatch.Length + int32(1),
					}] = true
				}
			}
		case parser.CHAR_RANGE:
			charRange := p.PrimitiveStr
			for _, c := range charRange {
				ran := bre.Regexable.ContinueBwdSearchStr(string(c), &rightMatch.Range)
				if !ran.Empty() {
					concatResults[SuccinctRegexMatch{
						Range:	*ran,
						Length: rightMatch.Length + int32(1),
					}] = true
				}
			}
		}
	case parser.UNION:
		u := r.(parser.RegexUnion)
		left := bre.regexConcat(u.First, rightMatch)
		for k := range left {
			concatResults[k] = true
		}
		right := bre.regexConcat(u.Second, rightMatch)
		for k := range right {
			concatResults[k] = true
		}

	case parser.CONCAT:
		c := r.(parser.RegexConcat)
		rightOfLeftResults := bre.regexConcat(c.Right, rightMatch)
		for rightOfLeftMatch := range rightOfLeftResults {
			inter := bre.regexConcat(c.Left, rightOfLeftMatch)
			for k := range inter {
				concatResults[k] = true
			}
		}

	case parser.REPEAT:
		rep := r.(parser.RegexRepeat)
		switch rep.RepeatType {
		case parser.ZERO_OR_MORE:
			concatResults = bre.regexRepeatZeroOrMoreWithMatch(rep.Internal, rightMatch)
		case parser.ONE_OR_MORE:
			concatResults = bre.regexRepeatOneOrMoreWithMatch(rep.Internal, rightMatch)
		case parser.MIN_TO_MAX:
			concatResults = bre.regexRepeatMinToMaxWithMatch(rep.Internal, rightMatch, rep.Min, rep.Max)
		}

	default:
		panic("Invalid node in succinct regex parse tree.")
	}

	return concatResults
}



