package cache

import (
	"container/list"
	"sync"
)

type locker struct {
	sync.Mutex
	num int
}

type Cache struct {
	keys       sync.Mutex
	keysLocker map[string]*locker

	MaxEntries int
	ll         *list.List
	lock       sync.Mutex
	cache      map[string]*list.Element
}

type Entry struct {
	Key   string
	Value interface{}
}

func (c *Cache) Init(maxEntries int) *Cache {
	c.MaxEntries = maxEntries
	c.ll = list.New()
	c.cache = make(map[string]*list.Element)
	c.keysLocker = make(map[string]*locker)
	return c
}

func (c *Cache) Add(key string, val interface{}) (evict string, has bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		ele.Value.(*Entry).Value = val
		return
	}

	en := &Entry{Key: key, Value: val}
	ele := c.ll.PushFront(en)
	c.cache[key] = ele

	if c.ll.Len() > c.MaxEntries {
		ele := c.ll.Back()
		if ele != nil {
			evict = ele.Value.(*Entry).Key
			has = true
		}
	}

	return
}

func (c *Cache) Get(key string) (interface{}, bool) {
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*Entry).Value, true
	}

	return nil, false
}

func (c *Cache) Each(fn func(key string, val interface{}) bool) {
	for key, ele := range c.cache {
		en := ele.Value.(*Entry)

		if !fn(key, en.Value) {
			break
		}
	}
}

func (c *Cache) Remove(key string) {
	if ele, hit := c.cache[key]; hit {
		c.ll.Remove(ele)
		kv := ele.Value.(*Entry)
		delete(c.cache, kv.Key)
	}
}
