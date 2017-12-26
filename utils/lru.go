package utils

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type Cache struct {
	size      int
	evictList *list.List
	items     map[interface{}]*list.Element
	lock      sync.RWMutex
	onEvicted func(key interface{}, value interface{})
}

type entry struct {
	key          interface{}
	value        interface{}
	expireAtSecs int64
}

func NewLru(size int) *Cache {
	cache, _ := NewLruWithEvict(size, nil)
	return cache
}

func NewLruWithEvict(size int, onEvicted func(key interface{}, value interface{})) (*Cache, error) {
	if size <= 0 {
		return nil, errors.New(T("utils.iru.with_evict"))
	}
	c := &Cache{
		size:      size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element, size),
		onEvicted: onEvicted,
	}
	return c, nil
}

func (c *Cache) Purge() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.onEvicted != nil {
		for k, v := range c.items {
			c.onEvicted(k, v.Value)
		}
	}

	c.evictList = list.New()
	c.items = make(map[interface{}]*list.Element, c.size)
}

func (c *Cache) Add(key, value interface{}) bool {
	return c.AddWithExpiresInSecs(key, value, 0)
}

func (c *Cache) AddWithExpiresInSecs(key, value interface{}, expireAtSecs int64) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if expireAtSecs > 0 {
		expireAtSecs = (time.Now().UnixNano() / int64(time.Second)) + expireAtSecs
	}

	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value.(*entry).value = value
		ent.Value.(*entry).expireAtSecs = expireAtSecs
		return false
	}

	ent := &entry{key, value, expireAtSecs}
	entry := c.evictList.PushFront(ent)
	c.items[key] = entry

	evict := c.evictList.Len() > c.size
	if evict {
		c.removeOldest()
	}
	return evict
}

func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ent, ok := c.items[key]; ok {

		if ent.Value.(*entry).expireAtSecs > 0 {
			if (time.Now().UnixNano() / int64(time.Second)) > ent.Value.(*entry).expireAtSecs {
				c.removeElement(ent)
				return nil, false
			}
		}

		c.evictList.MoveToFront(ent)
		return ent.Value.(*entry).value, true
	}
	return
}

func (c *Cache) Remove(key interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
	}
}

func (c *Cache) RemoveOldest() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.removeOldest()
}

func (c *Cache) Keys() []interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	keys := make([]interface{}, len(c.items))
	ent := c.evictList.Back()
	i := 0
	for ent != nil {
		keys[i] = ent.Value.(*entry).key
		ent = ent.Prev()
		i++
	}

	return keys
}

func (c *Cache) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.evictList.Len()
}

func (c *Cache) removeOldest() {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
	}
}

func (c *Cache) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
	if c.onEvicted != nil {
		c.onEvicted(kv.key, kv.value)
	}
}
