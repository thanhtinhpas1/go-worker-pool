package main

import (
	"time"

	simple_pool "github.com/thanhtinhpas1/go-worker-pool"
)

func main() {
	pool := simple_pool.NewSimplePool(50, 50)
	pool.Start()
	task := simple_pool.NewExampleTask()

	pool.AddTask(task)

	time.Sleep(5 * time.Second)
	pool.Stop()
}
