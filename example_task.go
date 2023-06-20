package simple_pool

import "log"

type exampleTask struct {
}

func NewExampleTask() Task {
	return &exampleTask{}
}

// OnExecute implements Task
func (*exampleTask) OnExecute() error {
	log.Printf("Executing example successfully")
	return nil
}

// OnFailure implements Task
func (*exampleTask) OnFailure(err error) {
	log.Fatalf("Failed to execute %v", err)
}

var _ Task = (*exampleTask)(nil)
