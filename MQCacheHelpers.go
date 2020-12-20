package mqcache

import "errors"

func loadSimpleSizeContainer(cache *MQCache, key string) (*SimpleSizeContainer, bool, error) {
	result, ok := cache.Get(key)
	if !ok {
		return nil, false, nil
	}
	container, ok := result.(*SimpleSizeContainer)
	if !ok {
		return nil, false, errors.New("key " + key + "in mqcache don't contain valid item")
	}
	return container, true, nil
}

//SaveString - save string in the cache
func SaveString(cache *MQCache, key, value string) {
	if 1 == cache.nilSize {
		cache.Set(key, NewSimpleSizeContainer(value))
		return
	}
	cache.Set(key, StringWithSize(value))
}

//LoadString - load string from the cache
//If error occurred  - item not string, it returns error
//If item does not present - return second argument as false
func LoadString(cache *MQCache, key string) (string, bool, error) {
	if 1 == cache.nilSize {
		container, ok, err := loadSimpleSizeContainer(cache, key)
		if err != nil || nil == container || !ok {
			return "", ok, err
		}
		str, ok := container.value.(string)
		if !ok {
			return "", false, errors.New("key " + key + "in mqcache don't contain string")
		}
		return str, true, nil
	}
	result, ok := cache.Get(key)
	if !ok {
		return "", false, nil
	}
	strWS, ok := result.(StringWithSize)
	if !ok {
		return "", false, errors.New("key " + key + "in mqcache don't contain string")
	}
	return string(strWS), true, nil
}

//SaveBytes - save []bytes in the cache
func SaveBytes(cache *MQCache, key string, value []byte) {
	if 1 == cache.nilSize {
		cache.Set(key, NewSimpleSizeContainer(value))
		return
	}
	cache.Set(key, ByteSliceWithSize(value))
}

//LoadBytes - load []bytes from the cache
//If error occurred  - item not []byte, it returns error
//If item does not present - return second argument as false
func LoadBytes(cache *MQCache, key string) ([]byte, bool, error) {
	if 1 == cache.nilSize {
		container, ok, err := loadSimpleSizeContainer(cache, key)
		if err != nil || nil == container || !ok {
			return nil, ok, err
		}
		bytes, ok := container.value.([]byte)
		if !ok {
			return nil, false, errors.New("key " + key + "in mqcache don't contain []byte")
		}
		return bytes, true, nil
	}
	result, ok := cache.Get(key)
	if !ok {
		return nil, false, nil
	}
	bytesWS, ok := result.(ByteSliceWithSize)
	if !ok {
		return nil, false, errors.New("key " + key + "in mqcache don't contain []byte")
	}
	return bytesWS, true, nil
}
