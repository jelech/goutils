package cacheutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryCache_SetAndGet(t *testing.T) {
	cache := NewMemoryCache()

	err := cache.Set("key1", "value1", time.Minute)
	assert.NoError(t, err)

	value, exists := cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)
}

func TestMemoryCache_GetNonExistent(t *testing.T) {
	cache := NewMemoryCache()

	value, exists := cache.Get("nonexistent")
	assert.False(t, exists)
	assert.Nil(t, value)
}

func TestMemoryCache_Delete(t *testing.T) {
	cache := NewMemoryCache()

	cache.Set("key1", "value1", time.Minute)

	value, exists := cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)

	err := cache.Delete("key1")
	assert.NoError(t, err)

	value, exists = cache.Get("key1")
	assert.False(t, exists)
	assert.Nil(t, value)
}

func TestMemoryCache_Clear(t *testing.T) {
	cache := NewMemoryCache()

	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)

	assert.Equal(t, 2, cache.Size())

	err := cache.Clear()
	assert.NoError(t, err)
	assert.Equal(t, 0, cache.Size())
}

func TestMemoryCache_TTL(t *testing.T) {
	cache := NewMemoryCache()

	// Set with short TTL
	cache.Set("key1", "value1", time.Millisecond*50)

	// Should exist immediately
	value, exists := cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)

	// Wait for expiration
	time.Sleep(time.Millisecond * 100)

	// Should not exist after expiration
	value, exists = cache.Get("key1")
	assert.False(t, exists)
	assert.Nil(t, value)
}

func TestMemoryCache_NoTTL(t *testing.T) {
	cache := NewMemoryCache()

	// Set without TTL (ttl = 0)
	cache.Set("key1", "value1", 0)

	// Should exist
	value, exists := cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)

	// Should still exist after some time
	time.Sleep(time.Millisecond * 50)
	value, exists = cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)
}

func TestLRUCache_SetAndGet(t *testing.T) {
	cache := NewLRUCache(2)

	err := cache.Set("key1", "value1", 0)
	assert.NoError(t, err)

	value, exists := cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)
}

func TestLRUCache_Capacity(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)
	cache.Set("key3", "value3", 0) // This should evict key1

	assert.Equal(t, 2, cache.Size())

	// key1 should be evicted
	_, exists := cache.Get("key1")
	assert.False(t, exists)

	// key2 and key3 should exist
	value, exists := cache.Get("key2")
	assert.True(t, exists)
	assert.Equal(t, "value2", value)

	value, exists = cache.Get("key3")
	assert.True(t, exists)
	assert.Equal(t, "value3", value)
}

func TestLRUCache_UpdateExisting(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)

	// Update existing key
	cache.Set("key1", "new_value1", 0)

	assert.Equal(t, 2, cache.Size())

	value, exists := cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "new_value1", value)
}

func TestLRUCache_LRUEviction(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)

	// Access key1 to make it recently used
	cache.Get("key1")

	// Add key3, which should evict key2 (least recently used)
	cache.Set("key3", "value3", 0)

	assert.Equal(t, 2, cache.Size())

	// key2 should be evicted
	_, exists := cache.Get("key2")
	assert.False(t, exists)

	// key1 and key3 should exist
	_, exists = cache.Get("key1")
	assert.True(t, exists)

	_, exists = cache.Get("key3")
	assert.True(t, exists)
}

func TestLRUCache_Delete(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)

	assert.Equal(t, 2, cache.Size())

	err := cache.Delete("key1")
	assert.NoError(t, err)
	assert.Equal(t, 1, cache.Size())

	_, exists := cache.Get("key1")
	assert.False(t, exists)
}

func TestLRUCache_Clear(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)

	assert.Equal(t, 2, cache.Size())

	err := cache.Clear()
	assert.NoError(t, err)
	assert.Equal(t, 0, cache.Size())
}

func BenchmarkMemoryCache_Set(b *testing.B) {
	cache := NewMemoryCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", time.Minute)
	}
}

func BenchmarkMemoryCache_Get(b *testing.B) {
	cache := NewMemoryCache()
	cache.Set("key", "value", time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}

func BenchmarkLRUCache_Set(b *testing.B) {
	cache := NewLRUCache(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", 0)
	}
}

func BenchmarkLRUCache_Get(b *testing.B) {
	cache := NewLRUCache(1000)
	cache.Set("key", "value", 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}
