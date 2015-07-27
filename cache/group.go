package cache

import (
	"sync/atomic"
	"strconv"
	"time"
	"log"
)


type Group struct {
	Name string
	saveTick int
	statusTick int
	isClosed bool
	miss_callback_url string
	evict_callback_url string

	count_miss int32
	count_hit int32
	count_evicted int32
	count_set int32
	count_incr int32
	count_remove int32

	cache *Cache
}

func (g *Group) Init(name string, maxEntries, saveTick, statusTick int, miss, evict string) *Group {
	g.Name = name
	g.saveTick = saveTick
	g.statusTick = statusTick
	g.miss_callback_url = miss
	g.evict_callback_url = evict
	g.cache = new(Cache).Init(maxEntries)
	return g
}

func (g *Group) Run() {
	if g.saveTick > 0 {
		go g.SaveTick()
	}
	if g.statusTick > 0 {
		go g.StatusTick()
	}
}

func (g *Group) Stop() {
	g.isClosed = true
}

func (g *Group) SaveTick() {
	c := time.Tick(time.Duration(g.saveTick) * time.Second)
	for _ = range c {
		if g.isClosed {
			return
		}

		save_to_dbfile(g)
	}
}

func (g *Group) StatusTick() {
	c := time.Tick(time.Duration(g.saveTick) * time.Second)
	for _ = range c {
		if g.isClosed {
			return
		}

		miss := atomic.SwapInt32(&g.count_miss, 0)
		hit := atomic.SwapInt32(&g.count_hit, 0)
		evict := atomic.SwapInt32(&g.count_evicted, 0)
		set := atomic.SwapInt32(&g.count_set, 0)
		incr := atomic.SwapInt32(&g.count_incr, 0)
		remove := atomic.SwapInt32(&g.count_remove, 0)

		log.Printf("group:%s status get:%d hit:%d%% set:%d incr:%d del:%d evicted:%d\n",
			g.Name, miss+hit, hit*100/(miss+hit), set, incr, remove, evict)
	}
}

func (g *Group) Get(key string) (interface{}, bool) {
	g.cache.LockKey(key)
	defer g.cache.UnlockKey(key)

	val, ok := g.cache.Get(key)
	if ok {
		atomic.AddInt32(&g.count_hit, 1)
	} else {
		atomic.AddInt32(&g.count_miss, 1)

		if g.miss_callback_url != "" {
			val, ok = g.OnMiss(key)
			if ok {
				g.cache.Add(key, val)
			}
		}
	}

	return val, ok
}

func (g *Group) Incr(key string, val int64) int64 {
	g.cache.LockKey(key)
	defer g.cache.UnlockKey(key)

	atomic.AddInt32(&g.count_incr, 1)

	var old int64

	_old, _ := g.cache.Get(key)
	switch tmp := _old.(type) {
	case int64:
		old = tmp
	case string:
		old, _ = strconv.ParseInt(tmp, 10, 64)
	}

	evict, has := g.cache.Add(key, old + val)
	if has {
		atomic.AddInt32(&g.count_evicted, 1)
		go g.Remove(evict)
	}

	return old+val
}

func (g *Group) Set(key string, val interface{}) {
	g.cache.LockKey(key)
	defer g.cache.UnlockKey(key)

	atomic.AddInt32(&g.count_set, 1)

	evict, has := g.cache.Add(key, val)
	if has {
		atomic.AddInt32(&g.count_evicted, 1)
		go g.Remove(evict)
	}
}

type HotVal struct {
	name string
	val int64
}

func (g *Group) Hot(num int) []HotVal {
	hot := make([]HotVal, num)

	g.cache.Each(func(key string, value interface{})bool{
		var val int64
		switch tmp := value.(type) {
		case int64:
			val = tmp
		case string:
			val, _ = strconv.ParseInt(tmp, 10, 64)
		}

		if val > hot[num-1].val {
			hot[num-1] = HotVal{name:key, val:val}
			for i:=num-2; i>0; i-- {
				if hot[i+1].val > hot[i].val {
					hot[i+1], hot[i] = hot[i], hot[i+1]
				} else {
					break
				}
			}
		}
		return true
	})

	for i:=0; i<num; i++ {
		if hot[i].name != "" {
			return hot[i:]
		}
	}

	return nil
}

func (g *Group) Remove(key string) {
	g.cache.LockKey(key)
	defer g.cache.UnlockKey(key)

	atomic.AddInt32(&g.count_remove, 1)

	if g.evict_callback_url != "" {
		val, ok := g.cache.Get(key)
		if ok {
			g.OnEvicted(key, val)
		}
	}

	g.cache.Remove(key)
}

func (g *Group) OnMiss(key string) (val interface{}, ok bool) {
	return
}

func (g *Group) OnEvicted(key string, val interface{}) {
}
