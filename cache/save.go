package cache

import (
	"encoding/gob"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

var SaveDir string

type SavedGroup struct {
	Name       string
	SaveTick   int
	StatusTick int
	MaxEntries int
	Expire     int64
	Entrys     []Entry
}

func copy_data(g *Group) (out *SavedGroup) {
	out = new(SavedGroup)
	out.Name = g.Name
	out.SaveTick = g.saveTick
	out.StatusTick = g.statusTick
	out.MaxEntries = g.cache.MaxEntries
	out.Expire = g.expire

	g.cache.lock.Lock()
	out.Entrys = make([]Entry, 0, g.cache.ll.Len())

	g.cache.Each(func(key string, value interface{}) bool {
		out.Entrys = append(out.Entrys, Entry{Key: key, Value: value})
		return true
	})
	g.cache.lock.Unlock()
	return
}

func from_data(out *SavedGroup) *Group {
	g := new(Group).Init(out.Name, out.MaxEntries, out.SaveTick, out.StatusTick, out.Expire)
	//反序添加
	for i := len(out.Entrys) - 1; i >= 0; i-- {
		g.cache.Add(out.Entrys[i].Key, out.Entrys[i].Value)
	}
	return g
}

var saveLock sync.Mutex

func delete_dbfile(g *Group) {
	saveLock.Lock()
	defer saveLock.Unlock()

	dbfile := path.Join(SaveDir, g.Name) + ".db"
	os.Remove(dbfile)
}

func save_to_dbfile(g *Group) {
	saveLock.Lock()
	defer saveLock.Unlock()

	t1 := time.Now()

	data := copy_data(g)
	dbfile := path.Join(SaveDir, g.Name) + ".db"

	os.Remove(dbfile + ".new")
	file, err := os.Create(dbfile + ".new")
	if err != nil {
		log.Println("db can not write.")
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(data)

	if err != nil {
		log.Println("data encode fail.")
		return
	}

	ok := rename_dbfile(dbfile)
	if ok {
		sec := float64(time.Now().Sub(t1)/time.Millisecond) / 1000
		log.Printf("group:%s saved in %.3fs\n", g.Name, sec)
	}
}

func rename_dbfile(dbfile string) bool {
	//先把原文件移到old
	os.Remove(dbfile + ".old")
	err := os.Rename(dbfile, dbfile+".old")
	if err != nil && !os.IsNotExist(err) {
		log.Println("move dbfile to old fail:", err)
		return false
	}
	//把现临时文件移到dbfile
	err = os.Rename(dbfile+".new", dbfile)
	if err != nil {
		log.Println("move dbfile.new to dbfile fail:", err)
		return false
	}
	//删除老文件
	err = os.Remove(dbfile + ".old")
	if err != nil && !os.IsNotExist(err) {
		log.Println("del old dbfile fail:", err)
		return false
	}

	return true
}

func Init_from_dbfile(dbfile string) (*Group, error) {
	data := new(SavedGroup)

	file, err := os.Open(dbfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(data)
	if err != nil {
		return nil, err
	}

	g := from_data(data)
	return g, nil
}
