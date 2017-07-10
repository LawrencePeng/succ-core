package util

type SuccinctConf struct {
	SaSamplingRate  int32
	IsaSamplingRate int32
	NpaSamplingRate int32
}

var DefaultSuccinctConf SuccinctConf = SuccinctConf{
	SaSamplingRate:  32,
	IsaSamplingRate: 32,
	NpaSamplingRate: 128,
}