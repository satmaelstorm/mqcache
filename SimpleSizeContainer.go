package mqcache

type SimpleSizeContainer struct {
	value interface{}
}

func (s *SimpleSizeContainer) Size() int64 {
	return 1
}

func (s *SimpleSizeContainer) SetValue(value interface{}) {
	s.value = value
}

func (s *SimpleSizeContainer) GetValue() interface{} {
	return s.value
}

func NewSimpleSizeContainer(value interface{}) *SimpleSizeContainer {
	return &SimpleSizeContainer{value: value}
}
