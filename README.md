# Go Simple Worker Pool
Example of a simple worker pool in Golang

## Overview
![overview](https://miro.medium.com/v2/resize:fit:720/format:webp/1*v9znVKY3FnBKGk5HkoSSNw.jpeg)

For controlling:
- How many tasks was executed at them same time?
- How quickly can tasks come in?
- How quickly can it handle the results
- How dynamically change the number of workers or tasks?

This simple solution will help you, work with a `tasks channel` as discover tasks coming to workers which manage by a pool.

```golang
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
```

Listening for tasks coming:

```golang
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
```