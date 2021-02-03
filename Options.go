package mqcache

import (
	"errors"
	"math"
	"time"
)

//HitFunc - function for promotions item in cache lines by hits
type HitFunc func(hits, maxNum int) int

//Options - Cache Parameters
type Options struct {
	//Queues number
	QueuesNum int
	//Capacity of cache
	Capacity int
	//Element lifetime in queues
	LifeTime time.Duration
	//Qout maximum length
	QOutLen int
	//Distribution function of the number of hits in line
	HitFunc HitFunc
	NilSize int64
}

//NewRecommendedOptions - creates recommended by algorithm authors parameters
func NewRecommendedOptions(capacity int, nilSize int64, lifeTimeInQueue time.Duration) *Options {
	opts := &Options{
		QueuesNum: 8,
		Capacity:  capacity,
		LifeTime:  lifeTimeInQueue,
		QOutLen:   4 * capacity,
		NilSize:   nilSize,
	}
	return opts
}

//NewSimpleOptionsWithCapacityByItems - creates options for cache with size by items count
func NewSimpleOptionsWithCapacityByItems(
	capacity, queuesNum, qOutLen int,
	lifeTimeInQueue time.Duration) *Options {
	opts := &Options{
		QueuesNum: queuesNum,
		Capacity:  capacity,
		LifeTime:  lifeTimeInQueue,
		QOutLen:   qOutLen,
		HitFunc:   nil,
		NilSize:   1,
	}
	return opts
}

//NewSimpleOptionsWithCapacityByBytes - creates options for cache with size by bytes
func NewSimpleOptionsWithCapacityByBytes(
	capacity, queuesNum, qOutLen int,
	lifeTimeInQueue time.Duration) *Options {
	opts := &Options{
		QueuesNum: queuesNum,
		Capacity:  capacity,
		LifeTime:  lifeTimeInQueue,
		QOutLen:   qOutLen,
		HitFunc:   nil,
		NilSize:   16,
	}
	return opts
}

//NewOptions - creates simple options
func NewOptions(
	capacity,
	queuesNum,
	qOutLen int,
	nilSize int64,
	lifeTimeInQueue time.Duration,
	hitFunc HitFunc,
) *Options {
	opts := &Options{
		QueuesNum: queuesNum,
		Capacity:  capacity,
		LifeTime:  lifeTimeInQueue,
		QOutLen:   qOutLen,
		HitFunc:   hitFunc,
		NilSize:   nilSize,
	}
	return opts
}

//Init - validate options
func (o *Options) Init() error {
	if o.QOutLen < 1 {
		return errors.New("length of out queue must be greater than 0")
	}
	if o.QueuesNum < 1 {
		return errors.New("number of queues must be greater than 0")
	}
	if o.Capacity < 1 {
		return errors.New("max items must be greater then 0")
	}
	if nil == o.HitFunc {
		o.HitFunc = RecommendedHitFunc
	}
	return nil
}

//RecommendedHitFunc - returns recommended by algorithm authors HitFunc
func RecommendedHitFunc(hits, queueNum int) int {
	if hits < 1 {
		return 0
	}
	r := int(math.Log(float64(hits)))
	if r < 0 {
		return 0
	} else if r >= queueNum {
		return queueNum - 1
	}
	return r
}
