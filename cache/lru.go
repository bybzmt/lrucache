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

func (c *Cache) LockKey(key string) {
	l := c.getKeyLocker(key)
	l.Lock()
}

func (c *Cache) UnlockKey(key string) {
	c.keys.Lock()
	defer c.keys.Unlock()

	l, ok := c.keysLocker[key]
	if ok {
		defer l.Unlock()

		l.num--
		if l.num < 1 {
			delete(c.keysLocker, key)
		}
	}
}

func (c *Cache) getKeyLocker(key string) *locker {
	c.keys.Lock()
	defer c.keys.Unlock()

	l, ok := c.keysLocker[key]
	if !ok {
		l = &locker{}
		c.keysLocker[key] = l
	}
	l.num++
	return l
}

func (c *Cache) Add(key string, val interface{}) (evict string, has bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		ele.Value.(*Entry).Value = val
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
	c.lock.Lock()
	defer c.lock.Unlock()

	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*Entry).Value, true
	}

	return nil, false
}

func (c *Cache) Each(fn func(key string, val interface{}) bool) {
	c.lock.Lock()
	ele := c.ll.Front()
	c.lock.Unlock()

	for ele != nil {
		en := ele.Value.(*Entry)

		if !fn(en.Key, en.Value) {
			break
		}

		//c.lock.Lock()
		ele = ele.Next()
		//c.lock.Unlock()
	}
}

func (c *Cache) Remove(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ele, hit := c.cache[key]; hit {
		c.ll.Remove(ele)
		kv := ele.Value.(*Entry)
		delete(c.cache, kv.Key)
	}
}

func (c *Cache) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.ll.Len()
}
