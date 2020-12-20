package mqcache

import (
	"github.com/stretchr/testify/suite"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type mqCacheTestSuite struct {
	suite.Suite
}

func TestMQCache(t *testing.T) {
	suite.Run(t, new(mqCacheTestSuite))
}

func (s *mqCacheTestSuite) TestPromote() {
	opts := NewRecommendedOptions(3, 1, time.Millisecond)
	s.Equal(8, opts.QueuesNum)
	cache, err := NewMQCache(opts)
	s.Nil(err)
	s.NotNil(cache)
	cache.Set("1", StringWithSize("1"))
	cache.Get("1")
	cache.Get("1")
	cache.Get("1")
	cache.Get("1")
	cache.Set("2", StringWithSize("2"))
	cache.Get("2")
	cache.Set("3", StringWithSize("3"))
	cache.Get("3")
	cache.Set("4", StringWithSize("4"))
	cache.Set("1", StringWithSize("1"))
	s.Nil(cache.q[0].peek("2"))
	cache.Set("2", nil)
	s.Equal(int64(3), cache.len())
	s.NotNil(cache.q[0].peek("2"))
	s.NotNil(cache.q[0].peek("4"))
	s.NotNil(cache.q[1].peek("1"))
	s.Nil(cache.q[0].peek("3"))
	r, ok := cache.Get("2")
	s.Nil(r)
	s.True(ok)
	r, ok = cache.Get("6")
	s.Nil(r)
	s.False(ok)
	cache.Get("2")
	cache.Get("2")
	cache.Get("2")
	cache.Get("2")
	time.Sleep(time.Millisecond)
	s.NotNil(cache.q[1].peek("1"))
	s.NotNil(cache.q[1].peek("2"))
	cache.Get("2")
	s.Nil(cache.q[1].peek("1"))
	s.NotNil(cache.q[1].peek("2"))
}

func (s *mqCacheTestSuite) TestSetAndGet() {
	opts := NewRecommendedOptions(3, 1, time.Millisecond)
	cache, err := NewMQCache(opts)
	s.Nil(err)
	s.NotNil(cache)
	cache.Set("1", StringWithSize("1"))
	cache.Set("2", StringWithSize("2"))
	var result SizeComputer
	var ok bool
	result, ok = cache.Get("1")
	s.True(ok)
	s.Equal("1", string(result.(StringWithSize)))
	result, ok = cache.Get("2")
	s.True(ok)
	s.Equal("2", string(result.(StringWithSize)))
	result, ok = cache.Get("3")
	s.False(ok)
	s.Nil(result)
}

func (s *mqCacheTestSuite) TestDelete() {
	opts := NewRecommendedOptions(3, 1, time.Millisecond)
	cache, err := NewMQCache(opts)
	s.Nil(err)
	s.NotNil(cache)
	cache.Set("1", StringWithSize("1"))
	var result SizeComputer
	var ok bool
	result, ok = cache.Get("1")
	s.True(ok)
	s.Equal("1", string(result.(StringWithSize)))
	cache.Delete("1")
	result, ok = cache.Get("1")
	s.False(ok)
	s.Nil(result)
}

func (s *mqCacheTestSuite) TestReplace() {
	opts := NewRecommendedOptions(3, 1, time.Millisecond)
	cache, err := NewMQCache(opts)
	s.Nil(err)
	s.NotNil(cache)
	cache.Set("1", StringWithSize("1"))
	var result SizeComputer
	var ok bool
	result, ok = cache.Get("1")
	s.True(ok)
	s.Equal("1", string(result.(StringWithSize)))
	cache.Set("1", StringWithSize("1111"))
	result, ok = cache.Get("1")
	s.True(ok)
	s.Equal("1111", string(result.(StringWithSize)))
}

func (s *mqCacheTestSuite) TestInvalidOpts() {
	opts := NewSimpleOptionsWithCapacityByItems(0, 0, 0, 0)
	s.NotNil(opts)
	cache, err := NewMQCache(opts)
	s.Nil(cache)
	s.NotNil(err)
}

func (s *mqCacheTestSuite) TestBytesSizeAndTooBig() {
	opts := NewSimpleOptionsWithCapacityByBytes(20, 2, 100, time.Second)
	s.NotNil(opts)
	cache, err := NewMQCache(opts)
	s.Nil(err)
	s.NotNil(cache)
	cache.Set("1", ByteSliceWithSize{1, 2, 3, 4, 5, 6, 7, 8, 9, 0})
	cache.Set("2", ByteSliceWithSize{1, 2, 3})
	cache.Set("3", ByteSliceWithSize{1, 2, 3, 4, 5, 6, 7})
	r, ok := cache.Get("1")
	s.True(ok)
	s.Equal(byte(1), r.(ByteSliceWithSize)[0])
	r, ok = cache.Get("2")
	s.True(ok)
	s.Equal(byte(1), r.(ByteSliceWithSize)[0])
	r, ok = cache.Get("3")
	s.True(ok)
	s.Equal(byte(1), r.(ByteSliceWithSize)[0])
	cache.Set("4", ByteSliceWithSize{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2})
	r, ok = cache.Get("1")
	s.False(ok)
	s.Nil(r)
	r, ok = cache.Get("2")
	s.False(ok)
	s.Nil(r)
	r, ok = cache.Get("3")
	s.True(ok)
	s.Equal(byte(1), r.(ByteSliceWithSize)[0])
	r, ok = cache.Get("4")
	s.True(ok)
	s.Equal(byte(2), r.(ByteSliceWithSize)[11])
	s.Panics(func() {
		cache.Set("5", ByteSliceWithSize{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1})
	})
}

func BenchmarkMQCache_Set(b *testing.B) {
	opts := NewRecommendedOptions(3000, 1, time.Microsecond)
	cache, _ := NewMQCache(opts)
	for i := 0; i < 3000; i++ {
		cache.Set(strconv.Itoa(i), IntWithSize(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(strconv.Itoa(rand.Intn(5000)), IntWithSize(rand.Intn(5000)))
	}
}

func BenchmarkMQCache_SetParallel(b *testing.B) {
	opts := NewRecommendedOptions(3000, 1, time.Millisecond)
	cache, _ := NewMQCache(opts)
	for i := 0; i < 3000; i++ {
		cache.Set(strconv.Itoa(i), IntWithSize(i))
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Set(strconv.Itoa(rand.Intn(5000)), IntWithSize(rand.Intn(5000)))
		}
	})
}

func BenchmarkMQCache_Get(b *testing.B) {
	opts := NewRecommendedOptions(3000, 1, time.Millisecond)
	cache, _ := NewMQCache(opts)
	for i := 0; i < 3000; i++ {
		cache.Set(strconv.Itoa(i), IntWithSize(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(strconv.Itoa(rand.Intn(5000)))
	}
}

func BenchmarkMQCache_GetParallel(b *testing.B) {
	opts := NewRecommendedOptions(3000, 1, time.Millisecond)
	cache, _ := NewMQCache(opts)
	for i := 0; i < 3000; i++ {
		cache.Set(strconv.Itoa(i), IntWithSize(i))
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Get(strconv.Itoa(rand.Intn(5000)))
		}
	})
}

func BenchmarkMQCache_SetGetParallel(b *testing.B) {
	opts := NewRecommendedOptions(3000, 1, time.Millisecond)
	cache, _ := NewMQCache(opts)
	for i := 0; i < 3000; i++ {
		cache.Set(strconv.Itoa(i), IntWithSize(i))
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if rand.Intn(10) < 5 {
				cache.Get(strconv.Itoa(rand.Intn(5000)))
			} else {
				cache.Set(strconv.Itoa(rand.Intn(5000)), IntWithSize(rand.Intn(5000)))
			}
		}
	})
}
