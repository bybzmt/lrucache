package cache

import (
	"sync"
	"errors"
	"path/filepath"
	"log"
)

var GroupExists = errors.New("GroupExists")
var GroupNotExists = errors.New("GroupNotExists")

type Hub struct {
	l sync.Mutex
	hub map[string]*Group
}

func (h *Hub) Init() *Hub {
	h.hub = make(map[string]*Group)
	return h
}

func (h *Hub) Create(name string, maxEnteries, saveTick, statusTick int, miss, evict string) error {
	h.l.Lock()
	defer h.l.Unlock()

	if _, ok := h.hub[name]; ok {
		return GroupExists
	}

	g := new(Group).Init(name, maxEnteries, saveTick, statusTick, miss, evict)
	h.hub[name] = g

	g.Run()

	return nil
}

func (h *Hub) Get(name string) (*Group, bool) {
	h.l.Lock()
	defer h.l.Unlock()

	g, ok := h.hub[name]
	return g, ok
}

func (h *Hub) Remove(name string) error {
	h.l.Lock()
	defer h.l.Unlock()

	if g, ok := h.hub[name]; ok {
		g.Stop()
		delete(h.hub, name)
		return nil
	}

	return GroupNotExists
}

//保存到文件
func (h *Hub) SaveToFile() {
	h.l.Lock()
	defer h.l.Unlock()

	for _, g := range h.hub {
		save_to_dbfile(g)
	}
}

//从文件中恢复
func (h *Hub) RecoveryFromFile() {
	h.l.Lock()
	defer h.l.Unlock()

	fs, err := filepath.Glob(filepath.Join(SaveDir, "*.db"))
	if err != nil {
		log.Println("RecoveryFromFile Error:", err)
		return
	}

	for _, name := range fs {
		g, err := Init_from_dbfile(name)
		if err != nil {
			log.Println("RecoveryFromFile Error:", err)
		} else {
			h.hub[g.Name] = g
			g.Run()
			log.Println("Recovery Group", g.Name)
		}
	}
}


