package cache

import (
	"testing"
)

func TestGet(t *testing.T) {
	var getTests = []struct {
		name       string
		keyToAdd   string
		keyToGet   string
		expectedOk bool
	}{
		{"string_hit", "myKey1", "myKey1", true},
		{"string_miss", "myKey2", "nonsense", false},
	}

	for _, tt := range getTests {
		lru := new(Cache).Init(10)
		lru.Add(tt.keyToAdd, 1234)
		val, ok := lru.Get(tt.keyToGet)
		if ok != tt.expectedOk {
			t.Fatalf("%s: cache hit = %v; want %v", tt.name, ok, !ok)
		} else if ok && val != 1234 {
			t.Fatalf("%s expected get to return 1234 but got %v", tt.name, val)
		}
	}
}

func TestRemove(t *testing.T) {
	lru := new(Cache).Init(10)
	lru.Add("myKey", 1234)
	if val, ok := lru.Get("myKey"); !ok {
		t.Fatal("TestRemove returned no match")
	} else if val != 1234 {
		t.Fatalf("TestRemove failed.  Expected %d, got %v", 1234, val)
	}

	lru.Remove("myKey")
	if _, ok := lru.Get("myKey"); ok {
		t.Fatal("TestRemove returned a removed entry")
	}
}

func TestEvict(t *testing.T) {
	lru := new(Cache).Init(1)
	//不驱逐
	e1, has := lru.Add("myKey", 1234)
	if has {
		t.Fatal("TestEvict returned no match")
	}

	//驱逐
	e1, has = lru.Add("myKey2", 1234)
	if !has {
		t.Fatal("TestEvict returned no match")
	} else if e1 != "myKey" {
		t.Fatal("TestEvict returned no match")
	}

	//驱逐
	lru = new(Cache).Init(0)
	e1, has = lru.Add("myKey", 1234)
	if !has {
		t.Fatal("TestEvict returned no match")
	} else if e1 != "myKey" {
		t.Fatal("TestEvict returned no match")
	}
}
