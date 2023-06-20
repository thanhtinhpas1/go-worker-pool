package simple_pool

import (
	"log"
	"sync"
)

const (
	defaultNumWorkers = 10
	defaultChanSize   = 10
)

type SimplePool interface {
	// Start the pool to ready process tasks
	Start()
	// Stop the pool to not receive tasks
	Stop()
	// Add a tasks to the channel
	AddTask(task Task)
	// Get result channel to observe tasks
	GetResultChan() chan bool
}

type Task interface {
	// Handle execution of task
	OnExecute() error
	// Handle when the task executed failed
	OnFailure(err error)
}

type simplePool struct {
	numWorkers int
	taskChan   chan Task
	resultChan chan bool
	start      sync.Once
	stop       sync.Once
	quit       chan struct{}
}

func NewSimplePool(numWorkers, channSize int) SimplePool {
	finalNumWorkers := defaultNumWorkers
	finalChanSize := defaultChanSize

	if numWorkers > 0 {
		finalNumWorkers = numWorkers
	}

	if channSize > 0 {
		finalChanSize = channSize
	}

	taskChan := make(chan Task, finalChanSize)
	resultChan := make(chan bool, finalChanSize)

	return &simplePool{
		numWorkers: finalNumWorkers,
		taskChan:   taskChan,
		resultChan: resultChan,
		start:      sync.Once{},
		stop:       sync.Once{},
		quit:       make(chan struct{}),
	}
}

func (p *simplePool) Start() {
	p.start.Do(func() {
		p.startWorkers()
	})
}

// Stop implements SimplePool
func (p *simplePool) Stop() {
	p.stop.Do(func() {
		close(p.taskChan)
		close(p.resultChan)
	})
}

// AddTask implements SimplePool
func (p *simplePool) AddTask(task Task) {
	select {
	case p.taskChan <- task:
	case <-p.quit:
	}
}

// GetResultChan implements SimplePool
func (p *simplePool) GetResultChan() chan bool {
	return p.resultChan
}

func (p *simplePool) startWorkers() {
	for i := 0; i < p.numWorkers; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("SimplePool: recovered task panic: %v", r)
					return
				}
			}()

			for {
				select {
				case <-p.quit:
					return

				case task, ok := <-p.taskChan:
					if !ok {
						p.resultChan <- false
						return
					}

					if err := task.OnExecute(); err != nil {
						p.resultChan <- false
						task.OnFailure(err)
						continue
					}

					p.resultChan <- true
				}
			}
		}()
	}
}
