package mqcache

import "unsafe"

//SizeComputer - interface, for MQCache items.
type SizeComputer interface {
	//Size - function must return size of item in the cache
	Size() int64
}

//StringWithSize - implementation of SizeComputer for string
type StringWithSize string

//Size - implementation of SizeComputer for string
func (s StringWithSize) Size() int64 {
	return int64(len(s))
}

//ByteSliceWithSize - implementation of SizeComputer for []byte
type ByteSliceWithSize []byte

//Size - implementation of SizeComputer for []byte
func (b ByteSliceWithSize) Size() int64 {
	return int64(len(b))
}

//IntSliceWithSize - implementation of SizeComputer for int
type IntSliceWithSize []int

//Size - implementation of SizeComputer for int
func (i IntSliceWithSize) Size() int64 {
	l := len(i)
	size := int64(unsafe.Sizeof(l))
	return int64(l) * size
}

//IntWithSize - implementation of SizeComputer for []int
type IntWithSize int

//Size - implementation of SizeComputer for []int
func (i IntWithSize) Size() int64 {
	return int64(unsafe.Sizeof(i))
}
