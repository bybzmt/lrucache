package cache

import (
	"errors"
	"log"
	"regexp"
	"strconv"
)

var KeyNotExists = errors.New("KeyNotExists")

//反回状态
const (
	RET_SUCCESS          = 0 //0 成功
	RET_ERROR            = 1 //1 错误
	RET_GROUP_NOT_EXISTS = 2 //2 分组不存成
	RET_KEY_NOT_EXISTS   = 3 //3 Key不存在
	RET_GROUP_EXISTS     = 4 //4 分组己存成
)

type result struct {
	Ret  int         `json:"err"`
	Data interface{} `json:"data"`
}

type fetch map[string]string

var Groups *Hub

func Setup(saveDir string) {
	SaveDir = saveDir
	Groups = new(Hub).Init()
}

func doIncr(name, key, val string) *result {
	if key == "" {
		return &result{Ret: RET_ERROR, Data: "Key Name Not Empty!"}
	}
	if len(key) > 255 {
		return &result{Ret: RET_ERROR, Data: "Key Name Too Long!"}
	}

	g, ok := Groups.Get(name)
	if !ok {
		return &result{Ret: RET_GROUP_NOT_EXISTS, Data: GroupNotExists.Error()}
	}

	_val, _ := strconv.ParseInt(val, 10, 64)

	newVal := g.Incr(key, _val)

	return &result{Ret: RET_SUCCESS, Data: newVal}
}

func doSet(name, key, val string) *result {
	if key == "" {
		return &result{Ret: RET_ERROR, Data: "Key Name Not Empty!"}
	}
	if len(key) > 255 {
		return &result{Ret: RET_ERROR, Data: "Key Name Too Long!"}
	}

	g, ok := Groups.Get(name)
	if !ok {
		return &result{Ret: RET_GROUP_NOT_EXISTS, Data: GroupNotExists.Error()}
	}

	g.Set(key, val)

	return &result{Ret: RET_SUCCESS}
}

func doGet(name, key string) *result {
	g, ok := Groups.Get(name)
	if !ok {
		return &result{Ret: RET_GROUP_NOT_EXISTS, Data: GroupNotExists.Error()}
	}

	data, ok := g.Get(key)
	if ok {
		return &result{Ret: RET_SUCCESS, Data: data}
	} else {
		return &result{Ret: RET_KEY_NOT_EXISTS, Data: KeyNotExists.Error()}
	}
}

func doHot(name, num string) *result {
	g, ok := Groups.Get(name)
	if !ok {
		return &result{Ret: RET_GROUP_NOT_EXISTS, Data: GroupNotExists.Error()}
	}

	val, _ := strconv.ParseInt(num, 10, 32)
	if val > 100000 {
		val = 100000
	}

	if val < 1 {
		val = 1
	}

	hots := g.Hot(int(val))

	return &result{Ret: RET_SUCCESS, Data: hots}
}

func doDel(name, key string) *result {
	g, ok := Groups.Get(name)
	if !ok {
		return &result{Ret: RET_GROUP_NOT_EXISTS, Data: GroupNotExists.Error()}
	}

	g.Remove(key)

	return &result{Ret: RET_SUCCESS}
}

func doGroupCreate(name, num, save, status, expire string) *result {
	_cap, _ := strconv.ParseInt(num, 10, 32)
	_save, _ := strconv.ParseInt(save, 10, 32)
	_status, _ := strconv.ParseInt(status, 10, 32)
	_expire, _ := strconv.ParseInt(expire, 10, 32)

	if name == "" {
		return &result{Ret: RET_ERROR, Data: "Group Name Not Empty!"}
	}

	if len(name) > 255 {
		return &result{Ret: RET_ERROR, Data: "Group Name Too Long!"}
	}

	reg := "^[a-zA-Z0-9_]+$"
	ok, err := regexp.MatchString(reg, name)
	if !ok || err != nil {
		if err != nil {
			log.Println(err)
		}
		return &result{Ret: RET_ERROR, Data: "Group Name Not Valid! regexp:" + reg}
	}

	if _cap < 1 {
		return &result{Ret: RET_ERROR, Data: "Group Capacity Too Small!"}
	}

	err = Groups.Create(name, int(_cap), int(_save), int(_status), int(_expire))

	if err != nil {
		return &result{Ret: RET_GROUP_EXISTS, Data: err.Error()}
	} else {
		return &result{Ret: RET_SUCCESS}
	}
}

func doGroupDel(name string) *result {
	err := Groups.Remove(name)

	if err != nil {
		return &result{Ret: RET_ERROR, Data: err.Error()}
	} else {
		return &result{Ret: RET_SUCCESS}
	}
}
