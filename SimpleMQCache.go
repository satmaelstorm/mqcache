package mqcache

//SimpleMQCache - wrapper of MQCache for simple use with interface{}
type SimpleMQCache struct {
	cache *MQCache
}

//NewSimpleMQCache - creates new SimpleMQCache
func NewSimpleMQCache(opts *Options) (*SimpleMQCache, error) {
	cache, err := NewMQCache(opts)
	if err != nil {
		return nil, err
	}
	cache.nilSize = 1
	return &SimpleMQCache{cache: cache}, nil
}

//Set - stores item in the cache
func (c *SimpleMQCache) Set(key string, value interface{}) {
	c.cache.Set(key, NewSimpleSizeContainer(value))
}

//Get - load item from the cache
func (c *SimpleMQCache) Get(key string) (interface{}, bool) {
	if result, ok := c.cache.Get(key); ok {
		if result != nil {
			return result.(*SimpleSizeContainer).value, true
		}
		return nil, true
	}
	return nil, false
}

//Delete - deletes item in the cache
func (c *SimpleMQCache) Delete(key string) {
	c.cache.Delete(key)
}

//Len - return cache size
func (c *SimpleMQCache) Len() (totalLen int64, queuesLen []int64) {
	return c.cache.Len()
}
