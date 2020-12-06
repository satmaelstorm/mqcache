package mqcache

import "container/list"

type queue struct {
	items      map[string]*list.Element
	evictQueue *list.List
	lifeTime   int64
	length     int64
	nilSize    int64
}

func newQueue(maxItems int, nilSize, lifeTime int64) *queue {
	return &queue{
		items:      make(map[string]*list.Element, maxItems),
		evictQueue: list.New(),
		lifeTime:   lifeTime,
		nilSize:    nilSize,
		length:     0,
	}
}

func (c *queue) store(key string, value SizeComputer, hits int64, now int64) {
	if e, ok := c.items[key]; ok {
		c.evictQueue.MoveToBack(e)
		cacheEntry := e.Value.(*entry)
		oldSize := c.getSize(cacheEntry.value)
		cacheEntry.value = value
		cacheEntry.hits = hits
		cacheEntry.expire = now + c.lifeTime
		e.Value = cacheEntry
		c.length += c.getSize(value) - oldSize
		return
	}

	item := &entry{
		key:    key,
		value:  value,
		hits:   hits,
		expire: now + c.lifeTime,
	}
	c.length += c.getSize(value)
	c.items[key] = c.evictQueue.PushBack(item)
	return
}

func (c *queue) peekFirst() *entry {
	e := c.evictQueue.Front()
	if e != nil {
		return e.Value.(*entry)
	}
	return nil
}

func (c *queue) peek(key string) (*entry, bool) {
	if e, ok := c.items[key]; ok {
		if e.Value.(*entry) == nil {
			return nil, false
		}
		return e.Value.(*entry), true
	}
	return nil, false
}

func (c *queue) delete(key string) *entry {
	if e, ok := c.items[key]; ok {
		return c.remove(e)
	}
	return nil
}

func (c *queue) len() int64 {
	return c.length
}

func (c *queue) evict() *entry {
	e := c.evictQueue.Front()
	if e != nil {
		return c.remove(e)
	}
	return nil
}

func (c *queue) remove(e *list.Element) *entry {
	c.evictQueue.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
	c.length -= c.getSize(kv.value)
	return kv
}

func (c *queue) getSize(s SizeComputer) int64 {
	if nil == s {
		return c.nilSize
	}
	return s.Size()
}
