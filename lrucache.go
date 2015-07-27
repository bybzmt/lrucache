package main

import (
	"net/http"
	"runtime"
	"flag"
	"time"
	"log"
	"os"
	"os/signal"
	"./cache"
)

var addr = flag.String("addr", ":80", "listen addr")
var dbdir = flag.String("dir", "./", "db file dir")

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Println("Runing.")

	cache.Setup(*dbdir)
	cache.Groups.RecoveryFromFile()

	go runHttp()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c
	cache.Groups.SaveToFile()

	log.Println("Closed.")
	os.Exit(0)
}

func runHttp() {
	mux := http.NewServeMux()
	mux.HandleFunc("/counter/hot", cache.HotAction)
	mux.HandleFunc("/counter/incr", cache.IncrAction)
	mux.HandleFunc("/cache/set", cache.SetAction)
	mux.HandleFunc("/cache/get", cache.GetAction)
	mux.HandleFunc("/cache/del", cache.DelAction)
	mux.HandleFunc("/group/create", cache.GroupCreateAction)
	mux.HandleFunc("/group/del", cache.GroupDelAction)

	mux.HandleFunc("/multiple/counter/hot", cache.HotsAction)
	mux.HandleFunc("/multiple/counter/incr", cache.IncrsAction)
	mux.HandleFunc("/multiple/cache/set", cache.SetAction)
	mux.HandleFunc("/multiple/cache/get", cache.GetAction)
	mux.HandleFunc("/multiple/cache/del", cache.DelAction)
	mux.HandleFunc("/multiple/group/create", cache.GroupCreatesAction)
	mux.HandleFunc("/multiple/group/del", cache.GroupDelsAction)

	log.Fatal(http.ListenAndServe(*addr, nil))

	s := http.Server{
		Addr: *addr,
		Handler:mux,
		MaxHeaderBytes: 1024 * 8,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatalln(s.ListenAndServe())
}


