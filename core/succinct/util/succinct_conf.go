package util

type SuccinctConf struct {
	SaSamplingRate  int
	IsaSamplingRate int
	NpaSamplingRate int
}

var DefaultSuccinctConf SuccinctConf = SuccinctConf{
	SaSamplingRate:  32,
	IsaSamplingRate: 32,
	NpaSamplingRate: 128,
}