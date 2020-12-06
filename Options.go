package mqcache

import (
	"errors"
	"math"
	"time"
)

type HitFunc func(hits, maxNum int) int

//Cache Parameters
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

func RecommendedHitFunc(hits, maxNum int) int {
	if hits < 1 {
		return 0
	}
	r := int(math.Log(float64(hits)))
	if r < 0 {
		return 0
	} else if r > maxNum {
		return maxNum
	}
	return r
}
