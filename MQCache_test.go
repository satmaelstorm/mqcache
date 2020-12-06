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

func (s *mqCacheTestSuite) TestBasic() {
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
