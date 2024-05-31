package cache

import (
	"errors"
	"sync"
)

type Getter interface {
	Get(key string) (ByteView, error)
}

type Group struct {
	name   string
	getter Getter
	cache  mutexcache
}

var (
	mu     sync.RWMutex
	groups = map[string]*Group{}
)

func NewGroup(name string, cacheBtyes int64, getter Getter) *Group {
	if getter == nil {
		panic("getter can not nil")
	}
	mu.Lock()
	defer mu.Unlock()

	g := &Group{name: name, getter: getter, cache: mutexcache{maxBtyes: cacheBtyes}}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key is required")
	}
	if v, ok := g.cache.get(key); ok {
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.loadLocal(key)
}

func (g *Group) loadLocal(key string) (ByteView, error) {
	v, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	g.cache.add(key, v)
	return v, err
}
