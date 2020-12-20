package mqcache

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type helpersMqCacheTestSuite struct {
	suite.Suite
	cacheByItems *MQCache
	cacheBySize  *MQCache
}

func TestMQCacheHelpers(t *testing.T) {
	suite.Run(t, new(helpersMqCacheTestSuite))
}

func (s *helpersMqCacheTestSuite) SetupSuite() {
	opts := NewRecommendedOptions(100, 1, time.Second)
	s.cacheByItems, _ = NewMQCache(opts)
	opts2 := NewRecommendedOptions(100, 16, time.Second)
	s.cacheBySize, _ = NewMQCache(opts2)
	s.cacheByItems.Set("notValid", &SimpleSizeContainer{value: 123})
	s.cacheByItems.Set("notValid2", IntSliceWithSize{1, 2, 3})
	s.cacheBySize.Set("notValid", IntWithSize(64))
}

func (s *helpersMqCacheTestSuite) TestStringSimple() {
	SaveString(s.cacheByItems, "1", "1")
	result, ok, err := LoadString(s.cacheByItems, "1")
	s.Nil(err)
	s.True(ok)
	s.Equal("1", result)
	result, ok, err = LoadString(s.cacheByItems, "2")
	s.Nil(err)
	s.False(ok)
	s.Equal("", result)
	result, ok, err = LoadString(s.cacheByItems, "notValid")
	s.NotNil(err)
	s.False(ok)
	s.Equal("", result)
	result, ok, err = LoadString(s.cacheByItems, "notValid2")
	s.NotNil(err)
	s.False(ok)
	s.Equal("", result)
}

func (s *helpersMqCacheTestSuite) TestStringBySize() {
	SaveString(s.cacheBySize, "1", "1")
	result, ok, err := LoadString(s.cacheBySize, "1")
	s.Nil(err)
	s.True(ok)
	s.Equal("1", result)
	result, ok, err = LoadString(s.cacheBySize, "2")
	s.Nil(err)
	s.False(ok)
	s.Equal("", result)
	result, ok, err = LoadString(s.cacheBySize, "notValid")
	s.NotNil(err)
	s.False(ok)
	s.Equal("", result)
}

func (s *helpersMqCacheTestSuite) TestBytesSimple() {
	SaveBytes(s.cacheByItems, "1", []byte{1, 2, 3})
	result, ok, err := LoadBytes(s.cacheByItems, "1")
	s.Nil(err)
	s.True(ok)
	s.Equal(byte(1), result[0])
	s.Equal(byte(2), result[1])
	s.Equal(byte(3), result[2])
	result, ok, err = LoadBytes(s.cacheByItems, "2")
	s.Nil(err)
	s.False(ok)
	s.Nil(result)
	result, ok, err = LoadBytes(s.cacheByItems, "notValid")
	s.NotNil(err)
	s.False(ok)
	s.Nil(result)
	result, ok, err = LoadBytes(s.cacheByItems, "notValid2")
	s.NotNil(err)
	s.False(ok)
	s.Nil(result)
}

func (s *helpersMqCacheTestSuite) TestBytesBySize() {
	SaveBytes(s.cacheBySize, "1", []byte{1, 2, 3})
	result, ok, err := LoadBytes(s.cacheBySize, "1")
	s.Nil(err)
	s.True(ok)
	s.Equal(byte(1), result[0])
	s.Equal(byte(2), result[1])
	s.Equal(byte(3), result[2])
	result, ok, err = LoadBytes(s.cacheBySize, "2")
	s.Nil(err)
	s.False(ok)
	s.Nil(result)
	result, ok, err = LoadBytes(s.cacheBySize, "notValid")
	s.NotNil(err)
	s.False(ok)
	s.Nil(result)
}
