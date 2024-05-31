package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("Get hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("Get hit not value key2")
	}
}

func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "v1", "v2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get(k1); ok {
		t.Fatalf("Test Removeoldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "v1", "v2", "v3"

	maps := map[string]String{}

	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), func(key string, value Value) {
		maps[key] = value.(String)
	})
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	expect := map[string]String{"key1": String("v1")}

	if !reflect.DeepEqual(expect, maps) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
