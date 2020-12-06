package mqcache

type entry struct {
	key    string
	value  SizeComputer
	hits   int64
	expire int64
}

func (e *entry) hit() {
	e.hits += 1
}
