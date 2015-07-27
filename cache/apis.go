package cache

import (
	"strconv"
	"errors"
)

var KeyNotExists = errors.New("KeyNotExists")

//反回状态
const (
	RET_SUCCESS = 0          //0 成功
	RET_ERROR = 1            //1 错误
	RET_GROUP_NOT_EXISTS = 2 //2 分组不存成
	RET_KEY_NOT_EXISTS = 3   //3 Key不存在
)

type result struct {
	Ret int          `json:"err"`
	Data interface{} `json:"data"`
}

type fetch map[string]string

var Groups *Hub

func Setup(saveDir string) {
	SaveDir = saveDir
	Groups = new(Hub).Init()
}

func doIncr(name, key, val string) *result {
	_val, _ := strconv.ParseInt(val, 10, 64)

	g, ok := Groups.Get(name)
	if !ok {
		return &result{Ret:RET_GROUP_NOT_EXISTS, Data:GroupNotExists}
	}

	newVal := g.Incr(key, _val)

	return &result{Ret:RET_SUCCESS, Data:newVal}
}

func doSet(name, key, val string) *result {
	g, ok := Groups.Get(name)
	if !ok {
		return &result{Ret:RET_GROUP_NOT_EXISTS, Data:GroupNotExists}
	}

	g.Set(key, val)

	return &result{Ret:RET_SUCCESS}
}

func doGet(name, key string) *result {
	g, ok := Groups.Get(name)
	if !ok {
		return &result{Ret:RET_GROUP_NOT_EXISTS, Data:GroupNotExists}
	}

	data, ok := g.Get(key)
	if ok {
		return &result{Ret:RET_SUCCESS, Data:data}
	} else {
		return &result{Ret:RET_KEY_NOT_EXISTS, Data:KeyNotExists}
	}
}

func doHot(name, num string) *result {
	g, ok := Groups.Get(name)
	if !ok {
		return &result{Ret:RET_GROUP_NOT_EXISTS, Data:GroupNotExists}
	}

	val, _ := strconv.ParseInt(num, 10, 32)
	if val > 10000 {
		val = 10000
	}

	hots := g.Hot(int(val))

	return &result{Ret:RET_SUCCESS, Data:hots}
}

func doDel(name, key string) *result {
	g, ok := Groups.Get(name)
	if !ok {
		return &result{Ret:RET_GROUP_NOT_EXISTS, Data:GroupNotExists}
	}

	g.Remove(key)

	return &result{Ret:RET_SUCCESS}
}

func doGroupCreate(name, num, save, status, miss, evict string) *result {
	_cap, _ := strconv.ParseInt(num, 10, 32)
	_save, _ := strconv.ParseInt(save, 10, 32)
	_status, _ := strconv.ParseInt(status, 10, 32)

	err := Groups.Create(name, int(_cap), int(_save), int(_status), miss, evict)

	if err != nil {
		return &result{Ret:RET_ERROR, Data:err}
	} else {
		return &result{Ret:RET_SUCCESS}
	}
}

func doGroupDel(name string) *result {
	err := Groups.Remove(name)

	if err != nil {
		return &result{Ret:RET_ERROR, Data:err}
	} else {
		return &result{Ret:RET_SUCCESS}
	}
}

