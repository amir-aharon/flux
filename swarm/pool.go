package swarm

import "sync"

type Task[T any] struct {
	ID   int
	Data T
}

type Result[T any] struct {
	TaskID int
	Value  T
	Err    error
}

type Pool[TaskT any, ResT any] struct {
	workersCount int
	in           chan Task[TaskT]
	out          chan Result[ResT]
	process      func(TaskT) (ResT, error)
	wg           sync.WaitGroup
}

func NewPool[TaskT any, ResT any](numWorkers, queueSize int, processOp func(TaskT) (ResT, error)) *Pool[TaskT, ResT] {
	return &Pool[TaskT, ResT]{
		workersCount: numWorkers,
		in:           make(chan Task[TaskT], queueSize),
		out:          make(chan Result[ResT], queueSize),
		process:      processOp,
		wg:           sync.WaitGroup{},
	}
}

func (p *Pool[TaskT, ResT]) Run() {
	for range p.workersCount {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for task := range p.in {
				processedValue, err := p.process(task.Data)
				p.out <- Result[ResT]{
					TaskID: task.ID,
					Value:  processedValue,
					Err:    err,
				}
			}
		}()
	}
	go func() {
		p.wg.Wait()
		close(p.out)
	}()
}

func (p *Pool[TaskT, ResT]) Submit(id int, data TaskT) { p.in <- Task[TaskT]{ID: id, Data: data} }

func (p *Pool[TaskT, ResT]) Results() <-chan Result[ResT] {
	return p.out
}

func (p *Pool[TaskT, ResT]) Shutdown() {
	close(p.in)
}
