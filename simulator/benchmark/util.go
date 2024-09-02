package main

import (
	"sync"
	"time"
)

type CyclicBarrier struct {
	n          int
	count      int
	cond       *sync.Cond
	beforeFunc func()
	afterFunc  func()
}

func NewCyclicBarrier(n int, beforeFunc func(), afterFunc func()) *CyclicBarrier {
	c := sync.NewCond(&sync.Mutex{})
	return &CyclicBarrier{
		n:          n,
		count:      0,
		cond:       c,
		beforeFunc: beforeFunc,
		afterFunc:  afterFunc,
	}
}

func (b *CyclicBarrier) await() {
	b.cond.L.Lock()
	defer b.cond.L.Unlock()

	b.beforeFunc()

	b.count += 1

	if b.count == b.n {
		b.count = 0
		b.cond.Broadcast()
		b.afterFunc()
		return
	}

	b.cond.Wait()
}

func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	ch := make(chan bool, 1)

	go time.AfterFunc(timeout, func() {
		ch <- true
	})

	go func() {
		wg.Wait()
		ch <- false
	}()

	return <-ch
}
