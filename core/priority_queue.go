package core

import (
	"errors"
	"math"
	"sync"
	"time"
)

const proxyNotFound = -1

// Queue is the priority queue implementation
type Queue struct {
	BlockWeight int64

	mut     *sync.Mutex
	nodes   []*node
	proxies []Proxy
}

// NewQueue returns a newly initialized priority queue
func NewQueue() *Queue {
	return &Queue{
		BlockWeight: 4,

		mut:     &sync.Mutex{},
		nodes:   []*node{},
		proxies: []Proxy{},
	}
}

// Rebuild rebuilds the priority with a clean slate from the given proxies
func (q *Queue) Rebuild(proxies []Proxy) {
	q.mut.Lock()
	defer q.mut.Unlock()

	q.nodes = make([]*node, 0, len(proxies))
	for i := range proxies {
		node := &node{
			i:           i,
			lastUsed:    time.Now(),
			blockWeight: q.BlockWeight,
		}
		q.nodes = append(q.nodes, node)
	}
	q.proxies = proxies
}

// Pop returns the top proxy
func (q *Queue) Pop() (Proxy, error) {
	q.mut.Lock()
	defer q.mut.Unlock()

	if len(q.nodes) == 0 {
		return Proxy{}, errors.New("Invalid queue state")
	}

	i := q.nodes[0].i
	q.nodes[0].setLastUsed(time.Now())
	q.heapifyIndex(0)
	return q.proxies[i], nil
}

func (q *Queue) resolveIndexFromAddress(addr string) int {
	for i, proxy := range q.proxies {
		if proxy.Addr == addr {
			return i
		}
	}
	return proxyNotFound
}

// Update blocks the proxy at index `proxyIdx`. The assumption is that
// the index corresponds to the original index of the proxy
// in the proxies array when `Rebuild` was called and that
// the order has not changed
func (q *Queue) Update(addr string) {
	proxyIdx := q.resolveIndexFromAddress(addr)
	if proxyIdx == proxyNotFound {
		return
	}
	q.mut.Lock()
	q.mut.Unlock()

	// TODO: make this more efficient?
	for i, node := range q.nodes {
		if node.i != proxyIdx {
			continue
		}
		node.incrBlock()
		q.heapifyIndex(i)
		break
	}
}

func (q *Queue) heapify() {
	for i := len(q.nodes)/2 - 1; i >= 0; i-- {
		q.heapifyIndex(i)
	}
}

func (q *Queue) heapifyIndex(i int) {
	min := i
	minScore, _ := q.nodeScore(min)
	left := left(i)
	right := right(i)
	if score, ok := q.nodeScore(left); ok && score < minScore {
		min = left
		minScore = score
	}
	if score, ok := q.nodeScore(right); ok && score < minScore {
		min = right
		minScore = score
	}
	if i != min {
		q.nodes[min], q.nodes[i] = q.nodes[i], q.nodes[min]
		q.heapifyIndex(min)
	}
}

func (q *Queue) nodeScore(i int) (int64, bool) {
	if i >= len(q.nodes) {
		return 0, false
	}
	return q.nodes[i].score, true
}

func left(i int) int {
	return 2*i + 1
}

func right(i int) int {
	return 2*i + 2
}

type node struct {
	i           int
	lastUsed    time.Time
	blockCount  int64
	blockWeight int64
	score       int64
}

func (node *node) setLastUsed(lastUsed time.Time) {
	node.lastUsed = lastUsed
	node.updateScore()
}

func (node *node) incrBlock() {
	node.blockCount++
	node.updateScore()
}

func (node *node) updateScore() {
	blockScore := int64(math.Pow(float64(node.blockCount), float64(node.blockWeight)))
	node.score = node.lastUsed.Unix() + blockScore
}
