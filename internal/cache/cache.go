package cache

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	data            map[string]*list.Element
	lru             *list.List
	ttl             time.Duration
	maxSize         int
	cleanupInterval time.Duration
}

type cacheItem struct {
	key       string
	value     interface{}
	timestamp time.Time
}

func NewCache(ttl time.Duration, maxSize int) *Cache {
	c := &Cache{
		data:            make(map[string]*list.Element),
		lru:             list.New(),
		ttl:             ttl,
		maxSize:         maxSize,
		cleanupInterval: time.Minute,
	}
	go c.startCleanup()
	return c
}

func (c *Cache) Set(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()

	// Если ключ уже существует, обновляем значение
	if elem, exists := c.data[key]; exists {
		c.lru.MoveToFront(elem)
		item := elem.Value.(*cacheItem)
		item.value = value
		item.timestamp = time.Now()
		return
	}

	// Если достигнут максимальный размер, удаляем старые записи
	if c.lru.Len() >= c.maxSize {
		oldest := c.lru.Back()
		if oldest != nil {
			item := oldest.Value.(*cacheItem)
			delete(c.data, item.key)
			c.lru.Remove(oldest)
		}
	}

	// Добавляем новую запись
	item := &cacheItem{
		key:       key,
		value:     value,
		timestamp: time.Now(),
	}
	elem := c.lru.PushFront(item)
	c.data[key] = elem
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	if elem, exists := c.data[key]; exists {
		item := elem.Value.(*cacheItem)
		if time.Since(item.timestamp) < c.ttl {
			c.lru.MoveToFront(elem)
			return item.value, true
		}
		// Запись устарела, удаляем её
		delete(c.data, key)
		c.lru.Remove(elem)
	}
	return nil, false
}

func (c *Cache) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	for range ticker.C {
		c.cleanup()
	}
}

func (c *Cache) cleanup() {
	c.Lock()
	defer c.Unlock()

	for elem := c.lru.Back(); elem != nil; {
		item := elem.Value.(*cacheItem)
		if time.Since(item.timestamp) > c.ttl {
			nextElem := elem.Prev()
			delete(c.data, item.key)
			c.lru.Remove(elem)
			elem = nextElem
		} else {
			break // Остальные записи свежие
		}
	}
}
