package regex

type RegexMatch struct {
	Offset, Length int
}

func (rm RegexMatch) Begin() int {
	return rm.Offset
}

func (rm RegexMatch) End() int {
	return rm.Offset + rm.Length
}

func (rm RegexMatch) Contains(ano RegexMatch) bool {
	return ano.Begin() >= rm.Begin() && ano.End() <= rm.End()
}

func END_COMP(a, b interface{}) int {
	o1 := a.(RegexMatch)
	o2 := b.(RegexMatch)

	o1End := o1.Offset + o1.Length
	o2End := o2.Offset + o2.Length
	if o1End == o2End {
		if o1.Length == o2.Length {
			return 0
		} else if o1.Length < o2.Length {
			return -1
		} else {
			return 1
		}
	}

	diff := o1End - o2End
	if diff < 0 {
		return -1
	} else {
		return 1
	}
}

func FRONT_COMP(a, b interface{}) int {
	o1 := a.(RegexMatch)
	o2 := b.(RegexMatch)

	if o1.Offset == o2.Offset {
		if o1.Length == o2.Length {
			return 0
		} else if o1.Length < o2.Length {
			return -1
		} else {
			return 1
		}
	}

	diff := o1.Offset - o2.Offset
	if diff < 0 {
		return -1
	} else {
		return 1
	}
}



