package concurrency

import (
	"fmt"
	"sync"
	"time"
)

// 优点：是官方出的，是亲儿子；通过读写分离，降低锁时间来提高效率；
// 缺点：不适用于大量写的场景，这样会导致read map读不到数据而进一步加锁读取，同时dirty map也会一直晋升为read map，整体性能较差。 适用场景：大量读，少量写
var s sync.RWMutex
var w sync.WaitGroup

func TestSyncMap() {
	mapTest()
	syncMapTest()
}

func mapTest() {
	m := map[int]int{1: 1}
	startTime := time.Now().Nanosecond()
	w.Add(1)
	go writeMap(m)
	w.Add(1)
	go writeMap(m)
	w.Add(1)
	go readMap(m)

	w.Wait()
	endTime := time.Now().Nanosecond()
	timeDiff := endTime - startTime
	fmt.Println("map:", timeDiff)
}

func writeMap(m map[int]int) {
	defer w.Done()
	i := 0
	for i < 10000 {
		// 加锁
		s.Lock()
		m[1] = 1
		// 解锁
		s.Unlock()
		i++
	}
}

func readMap(m map[int]int) {
	defer w.Done()
	i := 0
	for i < 10000 {
		s.RLock()
		_ = m[1]
		s.RUnlock()
		i++
	}
}

func syncMapTest() {
	m := sync.Map{}
	m.Store(1, 1)
	startTime := time.Now().Nanosecond()
	w.Add(1)
	go writeSyncMap(m)
	w.Add(1)
	go writeSyncMap(m)
	w.Add(1)
	go readSyncMap(m)

	w.Wait()
	endTime := time.Now().Nanosecond()
	timeDiff := endTime - startTime
	fmt.Println("sync.Map:", timeDiff)
}

func writeSyncMap(m sync.Map) {
	defer w.Done()
	i := 0
	for i < 10000 {
		m.Store(1, 1)
		i++
	}
}

func readSyncMap(m sync.Map) {
	defer w.Done()
	i := 0
	for i < 10000 {
		m.Load(1)
		i++
	}
}

// todo: ConcurrentMap
// 更好的优化：将map分片，变成行锁，只锁需要修改的部分
// package: orcaman/concurrent-map
