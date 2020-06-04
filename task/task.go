package task

import (
	"sync/atomic"
	"time"

	"github.com/shylinux/toolkits/log"
)

const (
	StatusWait = iota
	StatusDone
	StatusError
)

type Task struct {
	Arg interface{}
	CB  func(*Task) error
	ID  int64

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

var TaskID int64

func Run(arg interface{}, cb func(*Task) error) *Task {
	TaskID = atomic.AddInt64(&TaskID, 1)

	log.Show("task", "add", log.FileLine(cb, 3), "id", TaskID, "arg", arg)
	task := &Task{ID: TaskID, Arg: arg, CB: cb, BeginTime: time.Now()}
	go func() {
		task.Run()
	}()
	return task
}
