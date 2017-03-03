package balancer

import "container/heap"

// Pool implements heap.Interface and holds Workers.
// https://golang.org/pkg/container/heap/
type Pool []*connector

func newPool() Pool {
	p := Pool{}
	// initiate heap
	heap.Init(&p)
	return p
}

func (p Pool) Len() int           { return len(p) }
func (p Pool) Less(i, j int) bool { return p[i].pending < p[j].pending }
func (p Pool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
}

// Push ..
func (p *Pool) Push(item interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*p = append(*p, item.(*connector))
}

// Pop ...
func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	// item.index = -1 // for safety
	*p = old[0 : n-1]
	return item
}