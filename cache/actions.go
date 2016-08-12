package cache

import (
	"encoding/json"
	"net/http"
)

func doMultiple(w http.ResponseWriter, r *http.Request, fn func(fe fetch) *result) {
	tmp := r.FormValue("reqs")

	var ks []fetch

	err := json.Unmarshal([]byte(tmp), ks)
	if err != nil {
		json.NewEncoder(w).Encode(result{Ret: RET_ERROR, Data: err})
		return
	}

	res := make([]*result, len(ks))

	for i, fe := range ks {
		res[i] = fn(fe)
	}

	json.NewEncoder(w).Encode(res)
}

func IncrAction(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("group")
	key := r.FormValue("key")
	val := r.FormValue("val")

	res := doIncr(name, key, val)

	json.NewEncoder(w).Encode(res)
}

func IncrsAction(w http.ResponseWriter, r *http.Request) {
	doMultiple(w, r, func(fe fetch) *result {
		g := fe["group"]
		k := fe["key"]
		v := fe["val"]

		return doIncr(g, k, v)
	})
}

func SetAction(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("group")
	key := r.FormValue("key")
	val := r.FormValue("val")

	res := doSet(name, key, val)

	json.NewEncoder(w).Encode(res)
}

func SetsAction(w http.ResponseWriter, r *http.Request) {
	doMultiple(w, r, func(fe fetch) *result {
		g := fe["group"]
		k := fe["key"]
		v := fe["val"]

		return doSet(g, k, v)
	})
}

func GetAction(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("group")
	key := r.FormValue("key")

	res := doGet(name, key)

	json.NewEncoder(w).Encode(res)
}

func GetsAction(w http.ResponseWriter, r *http.Request) {
	doMultiple(w, r, func(fe fetch) *result {
		g := fe["group"]
		k := fe["key"]

		return doGet(g, k)
	})
}

func HotAction(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("group")
	_len := r.FormValue("len")

	res := doHot(name, _len)

	json.NewEncoder(w).Encode(res)
}

func HotsAction(w http.ResponseWriter, r *http.Request) {
	doMultiple(w, r, func(fe fetch) *result {
		g := fe["group"]
		n := fe["len"]

		return doHot(g, n)
	})
}

func DelAction(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("group")
	key := r.FormValue("key")

	res := doDel(name, key)

	json.NewEncoder(w).Encode(res)
}

func DelsAction(w http.ResponseWriter, r *http.Request) {
	doMultiple(w, r, func(fe fetch) *result {
		g := fe["group"]
		k := fe["key"]

		return doDel(g, k)
	})
}

func GroupCreateAction(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("group")
	_len := r.FormValue("cap")
	_save := r.FormValue("saveTick")
	_status := r.FormValue("statusTick")
	_expire := r.FormValue("expire")

	res := doGroupCreate(name, _len, _save, _status, _expire)

	json.NewEncoder(w).Encode(res)
}

func GroupCreatesAction(w http.ResponseWriter, r *http.Request) {
	doMultiple(w, r, func(fe fetch) *result {
		name := fe["group"]
		_len := fe["cap"]
		_save := fe["saveTick"]
		_status := fe["statusTick"]
		_expire := fe["expire"]

		return doGroupCreate(name, _len, _save, _status, _expire)
	})
}

func GroupDelAction(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("group")

	res := doGroupDel(name)

	json.NewEncoder(w).Encode(res)
}

func GroupDelsAction(w http.ResponseWriter, r *http.Request) {
	doMultiple(w, r, func(fe fetch) *result {
		name := fe["group"]

		return doGroupDel(name)
	})
}
