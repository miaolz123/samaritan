package task

import "sync"

// Task ...
type Task struct {
	waitGroup sync.WaitGroup
	Queue     int
}

// New : create a Task
func New() *Task {
	return &Task{}
}

// Add : add tasks
func (t *Task) Add(delta int) {
	t.Queue++
	t.waitGroup.Add(delta)
}

// Wait : wait all tasks done
func (t *Task) Wait() {
	t.waitGroup.Wait()
}

// Done : a task has done
func (t *Task) Done() {
	if t.Queue > 0 {
		t.Queue--
		t.waitGroup.Done()
	}
}

// AllDone : all tasks have done
func (t *Task) AllDone() bool {
	if t.Queue > 0 {
		t.Queue -= t.Queue
		t.waitGroup.Add(-t.Queue)
		return true
	}
	return false
}
