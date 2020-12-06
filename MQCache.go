package mqcache

import (
	"strconv"
	"sync"
	"time"
)

//MQCache - implementation of
//@see https://www.usenix.org/legacy/events/usenix01/full_papers/zhou/zhou.pdf
type MQCache struct {
	q         []*queue
	qOut      *timeoutFifo
	queues    int
	hitFunc   HitFunc
	lock      sync.Mutex
	capacity  int64
	queuesLen []int64 //memory cache for Len function
	nilSize   int64
}

func NewMQCache(opts *Options) (*MQCache, error) {
	if err := opts.Init(); err != nil {
		return nil, err
	}
	mq := new(MQCache)
	mq.queues = opts.QueuesNum
	mq.hitFunc = opts.HitFunc
	mq.qOut = newTimeoutFifo(opts.QOutLen)
	mq.q = make([]*queue, mq.queues)
	mq.capacity = int64(opts.Capacity)
	mq.queuesLen = make([]int64, mq.queues)
	mq.nilSize = opts.NilSize
	for i := 0; i < mq.queues; i++ {
		mq.q[i] = newQueue(opts.Capacity, mq.nilSize, int64(opts.LifeTime))
	}
	return mq, nil
}

func (c *MQCache) len() int64 {
	sum := int64(0)
	for i := 0; i < c.queues; i++ {
		sum += c.q[i].len()
	}
	return sum
}

func (c *MQCache) find(key string) (int, *entry) {
	for i := c.queues - 1; i >= 0; i-- {
		if e, ok := c.q[i].peek(key); ok {
			return i, e
		}
	}
	return 0, nil
}

//@see https://www.usenix.org/legacy/events/usenix01/full_papers/zhou/zhou.pdf
func (c *MQCache) adjust(now int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for i := 1; i < c.queues-1; i++ {
		cacheEntry := c.q[i].peekFirst()
		if cacheEntry != nil && cacheEntry.expire < now {
			c.q[i].delete(cacheEntry.key)
			c.q[i-1].store(cacheEntry.key, cacheEntry.value, cacheEntry.hits, now)
		}
	}
}

func (c *MQCache) evict() *entry {
	for i := 0; i < c.queues; i++ {
		if e := c.q[i].evict(); e != nil {
			return e
		}
	}
	return nil
}

func (c *MQCache) evictMinimumSize(size int64) {
	sumSize := int64(0)
	for sumSize < size {
		victim := c.evict()
		if victim == nil {
			panic("can't allocate " + strconv.Itoa(int(size)) + " capacity in MQCache")
		}
		sumSize += c.GetSize(victim.value)
		c.qOut.push(victim.key, victim.hits)
	}
}

func (c *MQCache) store(key string, value SizeComputer, now int64) {
	if idx, cacheEntry := c.find(key); cacheEntry != nil { //update
		length := c.len() + c.GetSize(value) - cacheEntry.value.Size()
		if length > c.capacity {
			c.evictMinimumSize(length - c.capacity)
		}
		c.q[idx].store(cacheEntry.key, cacheEntry.value, cacheEntry.hits, now)
		return
	}

	length := c.len() + c.GetSize(value)
	if length > c.capacity {
		c.evictMinimumSize(length - c.capacity)
	}
	hits, _ := c.qOut.get(key)
	queueNum := c.hitFunc(int(hits), c.queues)
	c.q[queueNum].store(key, value, hits, now)
	c.qOut.delete(key)
}

func (c *MQCache) load(key string, now int64) *entry {
	idx, cacheEntry := c.find(key)
	if cacheEntry == nil {
		return nil
	}
	c.q[idx].delete(cacheEntry.key)
	cacheEntry.hit()
	queueNum := c.hitFunc(int(cacheEntry.hits), c.queues)
	c.q[queueNum].store(cacheEntry.key, cacheEntry.value, cacheEntry.hits, now)
	return cacheEntry
}

func (c *MQCache) delete(key string) {
	idx, cacheEntry := c.find(key)
	if cacheEntry != nil {
		c.q[idx].delete(cacheEntry.key)
		c.qOut.delete(cacheEntry.key)
	}
}

func (c *MQCache) Set(key string, value SizeComputer) {
	now := time.Now().UnixNano()
	c.lock.Lock()
	c.store(key, value, now)
	c.lock.Unlock()
	c.adjust(now)
}

func (c *MQCache) Get(key string) (SizeComputer, bool) {
	now := time.Now().UnixNano()
	c.lock.Lock()
	cacheEntry := c.load(key, now)
	var result SizeComputer
	result = nil
	ok := false
	if cacheEntry != nil {
		result = cacheEntry.value
		ok = true
	}
	c.lock.Unlock()
	c.adjust(now)
	return result, ok
}

func (c *MQCache) Delete(key string) {
	c.lock.Lock()
	c.delete(key)
	c.lock.Unlock()
	c.adjust(time.Now().UnixNano())
}

func (c *MQCache) Len() (totalLen int64, queuesLen []int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for i := 0; i < c.queues; i++ {
		c.queuesLen[i] = c.q[i].len()
		totalLen += c.queuesLen[i]
	}
	return totalLen, c.queuesLen
}

func (c *MQCache) GetSize(s SizeComputer) int64 {
	if nil == s {
		return c.nilSize
	}
	return s.Size()
}

