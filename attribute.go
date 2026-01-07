package ipp

// Attribute defines an ipp attribute
type Attribute struct {
	Tag   int8
	Name  string
	Value any
}

// Resolution defines the resolution attribute
type Resolution struct {
	Height int32
	Width  int32
	Depth  int8
}
