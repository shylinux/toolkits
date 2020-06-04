package task

import (
	"time"
)

const (
	StatusWait = iota
	StatusDone
	StatusError
)

type Task struct {
	Arg interface{}
	CB  func(*Task) error

	ID     int64
	Error  error
	Status int

	BeginTime  time.Time
	StartTime  time.Time
	FinishTime time.Time
}

func (task *Task) Run() {
	// log.Show("task", "run", log.FileLine(task.CB, 3), "id", task.ID, "arg", task.Arg)
	// defer log.Cost("task ", log.FileLine(task.CB, 3), " id: ", task.ID, " err: ", task.Error, " ")()

	task.StartTime = time.Now()
	defer func() { task.FinishTime = time.Now() }()

	if e := task.CB(task); e != nil {
		task.Status = StatusError
		task.Error = e
		return
	}

	task.Status = StatusDone
}
