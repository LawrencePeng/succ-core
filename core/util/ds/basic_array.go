package ds

type BasicArray interface {
	Get(i int) int
	Set(i int, val int)
	Update(i int, val int) int
	Len() int
	Destroy()
}
