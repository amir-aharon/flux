package stash

import "sync"

type Node[K comparable, V any] struct {
	Key   K
	Value V
	Next  *Node[K, V]
	Prev  *Node[K, V]
}

func NewNode[K comparable, V any](k K, v V) *Node[K, V] {
	return &Node[K, V]{
		Key:   k,
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
}

type Cache[K comparable, V any] struct {
	Capacity int
	KV       map[K]*Node[K, V]
	MRU      *Node[K, V]
	LRU      *Node[K, V]
	mu       sync.Mutex
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		Capacity: capacity,
		KV:       make(map[K]*Node[K, V]),
		MRU:      nil,
		LRU:      nil,
		mu:       sync.Mutex{},
	}
}

func (c *Cache[K, V]) Put(k K, v V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.KV[k]; exists {
		node.Value = v
		c.removeNode(node)
		c.addToMRU(node)
		return
	}

	node := NewNode(k, v)
	c.KV[k] = node
	c.addToMRU(node)

	if len(c.KV) > c.Capacity {
		victim := c.LRU
		c.removeNode(victim)
		delete(c.KV, victim.Key)
	}
}

func (c *Cache[K, V]) Get(k K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.KV[k]
	if !exists {
		var zero V // each data type has a zero-value in go
		return zero, false
	}

	c.removeNode(node)
	c.addToMRU(node)
	return node.Value, true
}

func (c *Cache[K, V]) Remove(k K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.KV[k]; exists {
		c.removeNode(node)
		delete(c.KV, k)
	}
}

func (c *Cache[K, V]) addToMRU(node *Node[K, V]) {
	if c.MRU == nil {
		c.MRU = node
		c.LRU = node
		node.Next = nil
		node.Prev = nil
		return
	}

	node.Prev = c.MRU
	node.Next = nil
	c.MRU.Next = node
	c.MRU = node
}

func (c *Cache[K, V]) removeNode(node *Node[K, V]) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		c.LRU = node.Next
	}

	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		c.MRU = node.Prev
	}

	node.Next = nil
	node.Prev = nil
}
