package concurrency

import (
	"fmt"
	"sync"
	atomic "sync/atomic"
	"time"
)

type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func()) {
	// 防止当前协程重复执行
	if (atomic.LoadUint32(&o.done)) == 0 {
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	// o加锁
	o.m.Lock()
	defer o.m.Unlock()
	// 防止加锁期间其他协程执行
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

var once Once

func TestOnce() {
	var r int
	for i := 0; i < 5; i++ {
		f := func(i int) func() {
			return func() {
				time.Sleep(2 * time.Second)
				r = 1
			}
		}(i)
		once.Do(f)
	}
	fmt.Println("r", r)
}
