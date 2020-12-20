package mqcache

import (
	"bytes"
	"encoding/gob"
)

type item struct {
	Val interface{}
}

//GobMQCache is the implementation of MQCache encode data to gob
//and stores it like []byte. Easy way to build cache with limitation by memory
type GobMQCache struct {
	cache *MQCache
}

//NewGobMQCache create new GobMQCache
func NewGobMQCache(opts *Options) (*GobMQCache, error) {
	cache, err := NewMQCache(opts)
	if err != nil {
		return nil, err
	}
	cache.nilSize = 16
	return &GobMQCache{cache: cache}, nil
}

//Set store variable in the cache. Return error, if can't encode variable to gob
//Return number of stored bytes
func (c *GobMQCache) Set(key string, value interface{}) (int, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(item{Val: value})
	if err != nil {
		return 0, err
	}
	b := buf.Bytes()
	c.cache.Set(key, ByteSliceWithSize(b))
	return len(b), nil
}

//Get load variable from cache. Return error, if can't decode gob
//Return value and sign of success
func (c *GobMQCache) Get(key string) (interface{}, bool, error) {
	if result, ok := c.cache.Get(key); ok {
		if result != nil {
			var ret item
			buf := bytes.NewBuffer(result.(ByteSliceWithSize))
			decoder := gob.NewDecoder(buf)
			err := decoder.Decode(&ret)
			if err != nil {
				return nil, false, err
			}
			return ret.Val, true, nil
		}
		return nil, true, nil
	}
	return nil, false, nil
}

//Delete entry from cache
func (c *GobMQCache) Delete(key string) {
	c.cache.Delete(key)
}

//Len returns length of cache in bytes
func (c *GobMQCache) Len() (totalLen int64, queuesLen []int64) {
	return c.cache.Len()
}
