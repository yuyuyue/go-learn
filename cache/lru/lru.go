package lru

import "container/list"

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

type Lru struct {
	maxBtyes  int64
	usedBytes int64
	ll        *list.List                    // 双向链表用于操作
	cache     map[string]*list.Element      // map用于查找
	OnEvicted func(key string, value Value) // 移除时的回调
}

func New(maxBtyes int64, OnEvicted func(key string, value Value)) *Lru {
	return &Lru{maxBtyes: maxBtyes, OnEvicted: OnEvicted, ll: list.New(), cache: make(map[string]*list.Element)}
}

func (l *Lru) Get(key string) (value Value, ok bool) {
	if ele, ok := l.cache[key]; ok {
		l.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (l *Lru) Add(key string, value Value) {
	if ele, ok := l.cache[key]; ok {
		l.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		l.usedBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := l.ll.PushFront(&entry{key, value})
		l.cache[key] = ele
		l.usedBytes += int64(len(key)) + int64(value.Len())
	}
	for l.maxBtyes > 0 && l.maxBtyes < l.usedBytes {
		l.RemoveOldest()
	}
}

func (l *Lru) AddIgnore(key string, value Value, ignore bool) {
	if ele, ok := l.cache[key]; ok {
		l.ll.MoveToFront(ele)
		if ignore {
			return
		}
		kv := ele.Value.(*entry)
		l.usedBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := l.ll.PushFront(&entry{key, value})
		l.cache[key] = ele
		l.usedBytes += int64(len(key)) + int64(value.Len())
	}
	for l.maxBtyes > 0 && l.maxBtyes < l.usedBytes {
		l.RemoveOldest()
	}
}

func (l *Lru) RemoveOldest() {
	back := l.ll.Back()
	if back != nil {
		l.ll.Remove(back)
		kv := back.Value.(*entry)
		delete(l.cache, kv.key)
		l.usedBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if l.OnEvicted != nil {
			l.OnEvicted(kv.key, kv.value)
		}
	}
}
