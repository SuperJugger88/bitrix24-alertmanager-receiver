package worker

import (
	"context"
	"sync"
)

type Pool struct {
	workers int
	jobs    chan Job
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

type Job interface {
	Execute() error
}

func NewPool(workers int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		workers: workers,
		jobs:    make(chan Job, workers*2),
		ctx:     ctx,
		cancel:  cancel,
	}
	p.Start()
	return p
}

func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case job, ok := <-p.jobs:
			if !ok {
				return
			}
			if err := job.Execute(); err != nil {
				// обработка ошибки
			}
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Pool) Submit(job Job) {
	select {
	case p.jobs <- job:
	default:
		// очередь переполнена, выполняем в новой горутине
		go job.Execute()
	}
}

func (p *Pool) Stop() {
	p.cancel()
	close(p.jobs)
	p.wg.Wait()
}
