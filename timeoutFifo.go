package mqcache

import "container/list"

type timeoutFifoNode struct {
	key   string
	value int64
}

type timeoutFifo struct {
	values   map[string]*list.Element
	queue    *list.List
	maxItems int
}

func newTimeoutFifo(maxItems int) *timeoutFifo {
	return &timeoutFifo{
		values:   make(map[string]*list.Element, maxItems),
		queue:    list.New(),
		maxItems: maxItems,
	}
}

func (f *timeoutFifo) push(key string, value int64) {
	if e, ok := f.values[key]; ok {
		e.Value.(*timeoutFifoNode).value = value
		return
	}
	node := &timeoutFifoNode{key: key, value: value}
	f.values[key] = f.queue.PushBack(node)

	if f.queue.Len() > f.maxItems {
		e := f.queue.Front()
		if e != nil {
			f.queue.Remove(e)
			delete(f.values, e.Value.(*timeoutFifoNode).key)
		}
	}
}

func (f *timeoutFifo) get(key string) (int64, bool) {
	if e, ok := f.values[key]; ok {
		return e.Value.(*timeoutFifoNode).value, true
	}
	return 0, false
}

func (f *timeoutFifo) delete(key string) {
	if e, ok := f.values[key]; ok {
		f.queue.Remove(e)
		delete(f.values, e.Value.(*timeoutFifoNode).key)
	}
}

func (f *timeoutFifo) len() int {
	return f.queue.Len()
}
