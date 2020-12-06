package mqcache

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type gobMQCacheTestSuite struct {
	suite.Suite
}

func TestGobMQCache(t *testing.T) {
	suite.Run(t, new(gobMQCacheTestSuite))
}

func (s *gobMQCacheTestSuite) TestBasic() {
	opts := NewRecommendedOptions(30000, 16, time.Millisecond)
	cache, err := NewGobMQCache(opts)
	s.Nil(err)
	s.NotNil(cache)
	l, err := cache.Set("1", []int{1,2,3})
	cacheLen, _ := cache.Len()
	s.Nil(err)
	s.Equal(int64(l), cacheLen)
	r, ok, err := cache.Get("2")
	s.Nil(err)
	s.Nil(r)
	s.False(ok)
	r, ok, err = cache.Get("1")
	s.Nil(err)
	s.NotNil(r)
	s.True(ok)
	slice, ok := r.([]int)
	s.True(ok)
	s.Equal(3, len(slice))
	s.Equal(1, slice[0])
	s.Equal(2, slice[1])
	s.Equal(3, slice[2])
}
