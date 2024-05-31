package cache

import (
	"sync"

	"github.com/yuyuyue/cache/lru"
)

type mutexcache struct {
	mu       sync.Mutex
	lru      *lru.Lru
	maxBtyes int64
}

func (mc *mutexcache) get(key string) (v ByteView, ok bool) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if mc.lru == nil {
		return
	}
	if v, ok := mc.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}

func (mc *mutexcache) add(key string, value ByteView) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if mc.lru == nil {
		mc.lru = lru.New(mc.maxBtyes, nil)
	}
	mc.lru.Add(key, value)
}
