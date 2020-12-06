package mqcache

type SimpleMQCache struct {
	cache *MQCache
}

func NewSimpleMQCache(opts *Options) (*SimpleMQCache, error) {
	cache, err := NewMQCache(opts)
	if err != nil {
		return nil, err
	}
	cache.nilSize = 1
	return &SimpleMQCache{cache: cache}, nil
}

func (c *SimpleMQCache) Set(key string, value interface{}) {
	c.cache.Set(key, NewSimpleSizeContainer(value))
}

func (c *SimpleMQCache) Get(key string) (interface{}, bool) {
	if result, ok := c.cache.Get(key); ok {
		if result != nil {
			return result.(*SimpleSizeContainer).value, true
		}
		return nil, true
	}
	return nil, false
}

func (c *SimpleMQCache) Delete(key string) {
	c.cache.Delete(key)
}

func (c *SimpleMQCache) Len() (totalLen int64, queuesLen []int64) {
	return c.cache.Len()
}
