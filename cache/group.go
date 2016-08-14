package cache

import (
	"log"
	"sort"
	"strconv"
	"sync/atomic"
	"time"
)

type Group struct {
	Name       string
	saveTick   int
	statusTick int
	isClosed   bool
	expire     int64

	count_miss    int32
	count_hit     int32
	count_evicted int32
	count_set     int32
	count_incr    int32
	count_remove  int32

	cache *Cache
}

func (g *Group) Init(name string, maxEntries, saveTick, statusTick int, expire int64) *Group {
	g.Name = name
	g.saveTick = saveTick
	g.statusTick = statusTick
	g.cache = new(Cache).Init(maxEntries)
	g.expire = expire
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
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()

	c := time.Tick(time.Duration(g.saveTick) * time.Second)
	for _ = range c {
		if g.isClosed {
			return
		}

		save_to_dbfile(g)
	}
}

func (g *Group) StatusTick() {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()

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

		var _hit int32
		if miss+hit > 0 {
			_hit = hit * 100 / (miss + hit)
		}

		log.Printf("group:%s status %ds, get:%d hit:%d%% set:%d incr:%d del:%d evicted:%d\n",
			g.Name, g.saveTick, miss+hit, _hit, set, incr, remove, evict)
	}
}

func (g *Group) Get(key string) (interface{}, bool) {
	g.cache.lock.Lock()
	defer g.cache.lock.Unlock()

	val, ok := g.cache.Get(key)
	if ok {
		atomic.AddInt32(&g.count_hit, 1)
	} else {
		atomic.AddInt32(&g.count_miss, 1)
	}

	return val, ok
}

func (g *Group) Incr(key string, val int64) int64 {
	g.cache.lock.Lock()
	defer g.cache.lock.Unlock()

	atomic.AddInt32(&g.count_incr, 1)

	var old int64

	_old, _ := g.cache.Get(key)
	switch tmp := _old.(type) {
	case int64:
		old = tmp
	case string:
		old, _ = strconv.ParseInt(tmp, 10, 64)
	}

	evict, has := g.cache.Add(key, old+val)
	if has {
		atomic.AddInt32(&g.count_evicted, 1)
		g.cache.Remove(evict)
	}

	return old + val
}

func (g *Group) Set(key string, val interface{}) {
	g.cache.lock.Lock()
	defer g.cache.lock.Unlock()

	atomic.AddInt32(&g.count_set, 1)

	evict, has := g.cache.Add(key, val)
	if has {
		atomic.AddInt32(&g.count_evicted, 1)
		g.cache.Remove(evict)
	}
}

type HotVal struct {
	Name string `json:"name"`
	Val  int64 `json:"val"`
}

type HotVals []HotVal

func (h HotVals) Len() int {
	return len(h)
}

func (h HotVals) Less(i, j int) bool {
	return h[i].Val < h[j].Val
}

func (h HotVals) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (g *Group) Hot(num int) []HotVal {

	g.cache.lock.Lock()
	hot := make(HotVals, 0, g.cache.ll.Len())

	g.cache.Each(func(key string, value interface{}) bool {
		var val int64
		switch tmp := value.(type) {
		case int64:
			val = tmp
		case string:
			val, _ = strconv.ParseInt(tmp, 10, 64)
		}

		hot = append(hot, HotVal{Name: key, Val: val})

		return true
	})
	g.cache.lock.Unlock()

	sort.Sort(sort.Reverse(hot))

	if len(hot) < num {
		num = len(hot)
	}

	return hot[0:num]
}

func (g *Group) Remove(key string) {
	g.cache.lock.Lock()
	defer g.cache.lock.Unlock()

	atomic.AddInt32(&g.count_remove, 1)

	g.cache.Remove(key)
}
