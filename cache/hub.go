package cache

import (
	"errors"
	"log"
	"path/filepath"
	"sync"
	"time"
)

var GroupExists = errors.New("GroupExists")
var GroupNotExists = errors.New("GroupNotExists")

type Hub struct {
	l   sync.Mutex
	hub map[string]*Group
}

func (h *Hub) Init() *Hub {
	h.hub = make(map[string]*Group)
	return h
}

func (h *Hub) Create(name string, maxEnteries, saveTick, statusTick, expire_num int) error {
	h.l.Lock()
	defer h.l.Unlock()

	if _, ok := h.hub[name]; ok {
		return GroupExists
	}

	expire := int64(expire_num)
	if expire > 0 {
		expire += time.Now().Unix()
	}

	g := new(Group).Init(name, maxEnteries, saveTick, statusTick, expire)
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
		delete_dbfile(g)
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

func (h *Hub) Status(sec time.Duration) {
	c := time.Tick(sec * time.Second)
	for now := range c {
		unix := now.Unix()

		h.l.Lock()
		num := len(h.hub)

		//查找过期的组
		var expireName []string
		for _, g := range h.hub {
			if g.expire > 0 && g.expire < unix {
				expireName = append(expireName, g.Name)
			}
		}
		h.l.Unlock()

		//删除过期的组
		for _, name := range expireName {
			h.Remove(name)
		}

		log.Println("Status Groups num:", num)
	}
}
