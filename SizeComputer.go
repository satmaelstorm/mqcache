package mqcache

import "unsafe"

type SizeComputer interface {
	Size() int64
}

type StringWithSize string

func (s StringWithSize) Size() int64 {
	return int64(len(s))
}

type ByteSliceWithSize []byte

func (b ByteSliceWithSize) Size() int64 {
	return int64(len(b))
}

type IntSliceWithSize []int

func (i IntSliceWithSize) Size() int64 {
	l := len(i)
	size := int64(unsafe.Sizeof(l))
	return int64(l) * size
}

type IntWithSize int

func (i IntWithSize) Size() int64 {
	return int64(unsafe.Sizeof(i))
}
