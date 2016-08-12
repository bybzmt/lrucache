package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

var server = flag.String("url", "http://127.0.0.1", "listen addr")
var num = flag.Int("num", 10, "conns num")

var all_num1 int64
var all_num2 int64
var cc http.Client

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Println("Runing.")

	cc.Transport = &http.Transport{
		MaxIdleConnsPerHost: int(*num),
	}

	createGroup()

	for i := 0; i < *num; i++ {
		go func() {
			for {
				runtest()
			}
		}()
	}

	c := time.Tick(time.Second * 5)
	for _ = range c {
		n1 := atomic.LoadInt64(&all_num1)
		n2 := atomic.LoadInt64(&all_num2)

		log.Println("request all:", n1, "sec:", n2/5)

		atomic.StoreInt64(&all_num2, 0)
	}
}

func createGroup() {
	data := url.Values{}
	data.Add("group", "g1")
	data.Add("cap", "10000")
	data.Add("saveTick", "60")
	data.Add("statusTick", "60")
	data.Add("expire", "1800")

	resp, err := cc.PostForm(*server+"/group/create", data)

	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	ioutil.ReadAll(resp.Body)
}

func runtest() {
	data := url.Values{}
	data.Add("group", "g1")
	data.Add("key", "key"+strconv.Itoa(rand.Int()))
	data.Add("val", "1")

	resp, err := cc.PostForm(*server+"/counter/incr", data)

	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	ioutil.ReadAll(resp.Body)
	//d, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(d))

	atomic.AddInt64(&all_num1, 1)
	atomic.AddInt64(&all_num2, 1)
}
