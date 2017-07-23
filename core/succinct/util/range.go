package util

type Range struct {
	From int64
	To   int64
}

func (r *Range) Contains(point int64) bool {
	return point >= r.From && point <= r.To
}


func (r *Range) ContainsRange(rr *Range) bool {
	return rr.From >= r.From && rr.To <= r.To
}

func (r *Range) AdvanceFrom() {
	if !r.Empty() {
		r.From++
	}
}

func (r *Range) Empty() bool {
	return r.From > r.To
}

func (r *Range) Size() int64 {
	return r.To - r.From + 1
}



