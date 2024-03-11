package lrucache

import (
	"container/list"
	"sync"
)

var LruCache *LRUCache

func init() {
	LruCache = NewLRUCache(20)
}

type LRUCache struct {
	capacity int
	items    map[string]*list.Element
	evict    *list.List
	mutex    sync.Mutex
}

type CacheItem struct {
	Key   string
	Value string
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		evict:    list.New(),
	}
}

func (c *LRUCache) Get(key string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.items[key]; ok {
		c.evict.MoveToFront(elem)
		return elem.Value.(*CacheItem).Value, true
	}
	return "", false
}

func (c *LRUCache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.items[key]; ok {
		elem.Value.(*CacheItem).Value = value
		c.evict.MoveToFront(elem)
		return
	}

	if c.evict.Len() > c.capacity {
		evicted := c.evict.Back()
		if evicted != nil {
			delete(c.items, evicted.Value.(*CacheItem).Key)
			c.evict.Remove(evicted)
		}
	}
}
