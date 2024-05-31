package pkg

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const maxPoolCap = 20

var ErrWorkerPoolFreed = errors.New("workerpool freed") //

type Task func()

type Pool struct {
	cap    int
	active chan struct{}
	tasks  chan Task
	quit   chan struct{}
	wg     sync.WaitGroup
}

func (p *Pool) run() {
	idx := 0

	for {
		select {
		case <-p.quit:
			return
		case p.active <- struct{}{}:
			idx++
			p.newWorker(idx)
		}
	}
}

func (p *Pool) newWorker(i int) {
	p.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover panic[%s] and exit\n", i, err)
				<-p.active
			}
			p.wg.Done()
		}()

		fmt.Printf("worker[%03d]: start\n", i)
		for {
			select {
			case <-p.quit:
				fmt.Printf("worker[%03d]: exit\n", i)
				<-p.active
				return
			case t := <-p.tasks:
				fmt.Printf("worker[%03d]: receive a task\n", i)
				t()
			}
		}
	}()
}

func (p *Pool) Schedule(t Task) error {
	for {
		select {
		case <-p.quit:
			return ErrWorkerPoolFreed
		case p.tasks <- t:
			return nil
		}
	}
}

func (p *Pool) Free() {
	close(p.quit) // make sure all worker and p.run exit and schedule return error
	p.wg.Wait()
	fmt.Printf("workerpool freed\n")
}

func New(cap int) *Pool {
	if cap > maxPoolCap {
		cap = maxPoolCap
	}

	p := &Pool{
		cap:    cap,
		active: make(chan struct{}, cap),
		tasks:  make(chan Task),
		quit:   make(chan struct{}),
	}
	fmt.Println("Pool start")
	go p.run()
	return p
}

func WorkerpoolTest() {
	p := New(3)

	time.Sleep(time.Second * 2)
	for i := 0; i < 10; i++ {
		fn := func(i int) func() {
			return func() {
				fmt.Println("task: ", i, "run")
				time.Sleep(time.Second * 3)
			}
		}(i)
		err := p.Schedule(fn)
		if err != nil {
			println("task: ", i, "err:", err)
		}
	}
	p.Free()
}
