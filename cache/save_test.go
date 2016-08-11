package cache

import (
	"os"
	"testing"
	"path"
)


func TestSaveToFile(t *testing.T) {
	SaveDir = os.TempDir()

	var data = []string{"key1", "key2", "key3", "key4"}

	g := new(Group).Init("test", 10, 0, 0, "", "")
	for _, k := range data {
		g.Set(k, k)
	}

	save_to_dbfile(g)

	dbfile := path.Join(SaveDir, "test.db")

	g2, err := Init_from_dbfile(dbfile)

	if err != nil {
		t.Fatal("TestSaveToFile Error:%s", err)
		return
	}

	var data2 []string

	g2.cache.Each(func(key string, value interface{})bool{
		val, ok := value.(string)
		if !ok || key != val {
			t.Fatalf("%s: cache val not macth", key)
		}
		var tmp []string
		tmp = append(tmp, key)
		data2 = append(tmp, data2...)
		return true
	})

	if len(data) != len(data2) {
		t.Fatal("TestSaveToFile Data len error")
		return
	}

	for i:=0; i<len(data); i++ {
		if data[i] != data2[i] {
			t.Fatalf("TestSaveToFile %s: order not macth", data[i])
			return
		}
	}
}

