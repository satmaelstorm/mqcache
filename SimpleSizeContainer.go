package mqcache

//SimpleSizeContainer - container for SimpleMQCache
type SimpleSizeContainer struct {
	value interface{}
}

//Size - implements SizeComputer. Always size equal to 1
func (s *SimpleSizeContainer) Size() int64 {
	return 1
}

//SetValue - set value in the container
func (s *SimpleSizeContainer) SetValue(value interface{}) {
	s.value = value
}

//GetValue - get value from the container
func (s *SimpleSizeContainer) GetValue() interface{} {
	return s.value
}

//NewSimpleSizeContainer - creates new SimpleSizeContainer
func NewSimpleSizeContainer(value interface{}) *SimpleSizeContainer {
	return &SimpleSizeContainer{value: value}
}
