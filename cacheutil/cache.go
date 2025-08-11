// Package cache provides various caching implementations.
package cacheutil

import (
	"sync"
	"time"
)

// Cache defines the interface for cache implementations
type Cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, bool)
	Delete(key string) error
	Clear() error
	Size() int
}

// MemoryCache implements an in-memory cache with TTL support
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]*cacheItem
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]*cacheItem),
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Set stores a value in the cache with the specified TTL
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	c.items[key] = &cacheItem{
		value:     value,
		expiresAt: expiresAt,
	}

	return nil
}

// Get retrieves a value from the cache
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Check if item has expired
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		// Item has expired, but we don't delete it here to avoid deadlock
		// The cleanup goroutine will handle expired items
		return nil, false
	}

	return item.value, true
}

// Delete removes a value from the cache
func (c *MemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	return nil
}

// Clear removes all values from the cache
func (c *MemoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*cacheItem)
	return nil
}

// Size returns the number of items in the cache
func (c *MemoryCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// cleanup removes expired items from the cache
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if !item.expiresAt.IsZero() && now.After(item.expiresAt) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// LRUCache implements a Least Recently Used cache
type LRUCache struct {
	mu       sync.RWMutex
	capacity int
	items    map[string]*lruNode
	head     *lruNode
	tail     *lruNode
}

type lruNode struct {
	key   string
	value interface{}
	prev  *lruNode
	next  *lruNode
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		capacity = 100 // default capacity
	}

	cache := &LRUCache{
		capacity: capacity,
		items:    make(map[string]*lruNode),
	}

	// Initialize dummy head and tail nodes
	cache.head = &lruNode{}
	cache.tail = &lruNode{}
	cache.head.next = cache.tail
	cache.tail.prev = cache.head

	return cache
}

// Set stores a value in the LRU cache
func (c *LRUCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.items[key]; exists {
		// Update existing node
		node.value = value
		c.moveToHead(node)
		return nil
	}

	// Create new node
	node := &lruNode{
		key:   key,
		value: value,
	}

	c.items[key] = node
	c.addToHead(node)

	// Check capacity
	if len(c.items) > c.capacity {
		// Remove least recently used item
		tail := c.removeTail()
		delete(c.items, tail.key)
	}

	return nil
}

// Get retrieves a value from the LRU cache
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Move to head (mark as recently used)
	c.moveToHead(node)

	return node.value, true
}

// Delete removes a value from the LRU cache
func (c *LRUCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.items[key]; exists {
		c.removeNode(node)
		delete(c.items, key)
	}

	return nil
}

// Clear removes all values from the LRU cache
func (c *LRUCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*lruNode)
	c.head.next = c.tail
	c.tail.prev = c.head

	return nil
}

// Size returns the number of items in the LRU cache
func (c *LRUCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// addToHead adds a node right after the head
func (c *LRUCache) addToHead(node *lruNode) {
	node.prev = c.head
	node.next = c.head.next

	c.head.next.prev = node
	c.head.next = node
}

// removeNode removes an existing node from the linked list
func (c *LRUCache) removeNode(node *lruNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// moveToHead moves a node to the head
func (c *LRUCache) moveToHead(node *lruNode) {
	c.removeNode(node)
	c.addToHead(node)
}

// removeTail removes the last node and returns it
func (c *LRUCache) removeTail() *lruNode {
	lastNode := c.tail.prev
	c.removeNode(lastNode)
	return lastNode
}
