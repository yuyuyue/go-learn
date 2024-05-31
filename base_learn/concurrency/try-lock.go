package concurrency

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	mutexLocked      = 1 << iota // 加锁标识位
	mutexWoken                   // 唤醒标识位
	mutexStarving                // 锁饥饿标识位
	mutexWaiterShift = 3         // 标识waiter的起始bit位
)

type Mutex struct {
	sync.Mutex
}

func (m *Mutex) TryLock() bool {
	// 如果能成功抢到锁
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked) {
		return true
	}
	// 如果处于唤醒、加锁或者饥饿状态，这次请求就不参与竞争了，返回false
	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexStarving|mutexWoken) != 0 {
		return false
	}
	// 尝试在竞争的状态下请求锁
	new := old | mutexLocked
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), old, new)
}

func (m *Mutex) Count() int {
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	v = v >> mutexWaiterShift // 得到等待者的数值
	v = v + (v & mutexLocked) // 再加上锁持有者的数量，0或者1
	return int(v)
}

// 锁是否被持有
func (m *Mutex) IsLocked() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexLocked == mutexLocked
}

func (m *Mutex) IsWoken() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexWoken == mutexWoken
}

func (m *Mutex) IsStarving() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexStarving == mutexStarving
}

func count() {
	var mu Mutex
	for i := 0; i < 1000; i++ { // 启动1000个goroutine
		go func() {
			mu.Lock()
			time.Sleep(time.Second)
			mu.Unlock()
		}()
	}
	time.Sleep(time.Second)
	// 输出锁的信息
	fmt.Printf("waitings: %d, isLocked: %t, woken: %t,  starving: %t\n", mu.Count(), mu.IsLocked(), mu.IsWoken(), mu.IsStarving())
}

func TestTryLock() {
	// var mu Mutex
	// go func() { // 启动一个goroutine持有一段时间的锁
	// 	mu.Lock()
	// 	time.Sleep(2 * time.Second)
	// 	mu.Unlock()
	// }()
	// time.Sleep(time.Second)
	// ok := mu.TryLock() // 尝试获取到锁
	// if ok {            // 获取成功
	// 	fmt.Println("got the lock")
	// 	// do something
	// 	mu.Unlock()
	// 	return
	// }
	// // 没有获取到
	// fmt.Println("can't get the lock")
	count()
}
