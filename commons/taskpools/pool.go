package taskpools

import (
	"errors"

	"github.com/panjf2000/ants/v2"
)

var ErrPoolWaitQueueFull = errors.New("task pool wait queue is full")

type AsyncTaskPool struct {
	pool  *ants.Pool
	tasks chan func()
	done  chan struct{}
}

func NewAsyncTaskPool(workerSize, queueSize int) (*AsyncTaskPool, error) {
	pool, err := ants.NewPool(workerSize, ants.WithPreAlloc(true))
	if err != nil {
		return nil, err
	}
	p := &AsyncTaskPool{
		pool:  pool,
		tasks: make(chan func(), queueSize),
		done:  make(chan struct{}),
	}
	go p.dispatch()
	return p, nil
}

func (p *AsyncTaskPool) dispatch() {
	for task := range p.tasks {
		_ = p.pool.Submit(task)
	}
	close(p.done)
}

func (p *AsyncTaskPool) Submit(task func()) error {
	select {
	case p.tasks <- task:
		return nil
	default:
		return ErrPoolWaitQueueFull
	}
}

func (p *AsyncTaskPool) Release() {
	close(p.tasks)
	<-p.done
	p.pool.Release()
}
