package mqcache

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type simpleMqCacheTestSuite struct {
	suite.Suite
	cache *SimpleMQCache
}

func TestSimpleMQCache(t *testing.T) {
	suite.Run(t, new(simpleMqCacheTestSuite))
}

func (s *simpleMqCacheTestSuite) SetupTest() {
	opts := NewRecommendedOptions(10, 1, time.Second)
	var err error
	s.cache, err = NewSimpleMQCache(opts)
	s.Nil(err)
	s.NotNil(s.cache)
}

func (s *simpleMqCacheTestSuite) TestSetAndGet() {
	s.cache.Set("1", "1")
	s.cache.Set("2", "2")
	var r interface{}
	var ok bool
	r, ok = s.cache.Get("1")
	s.True(ok)
	s.Equal("1", r.(string))
	r, ok = s.cache.Get("2")
	s.True(ok)
	s.Equal("2", r.(string))
	r, ok = s.cache.Get("3")
	s.False(ok)
	s.Nil(r)
}

func (s *simpleMqCacheTestSuite) TestLen() {
	s.cache.Set("1", "1")
	s.cache.Set("2", "2")
	total, queues := s.cache.Len()
	s.Equal(int64(2), total)
	s.Equal(int64(2), queues[0])
	s.Equal(int64(0), queues[1])
	s.Equal(8, len(queues))
}

func (s *simpleMqCacheTestSuite) TestDelete() {
	s.cache.Set("1", "1")
	var r interface{}
	var ok bool
	r, ok = s.cache.Get("1")
	s.True(ok)
	s.Equal("1", r.(string))
	s.cache.Delete("1")
	r, ok = s.cache.Get("1")
	s.False(ok)
	s.Nil(r)
}

func (s *simpleMqCacheTestSuite) TestContainer() {
	c := NewSimpleSizeContainer(123)
	r := c.GetValue()
	s.Equal(123, r.(int))
	c.SetValue("123")
	r2 := c.GetValue()
	s.Equal("123", r2.(string))
}
