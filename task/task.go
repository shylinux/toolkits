package task

import (
	"context"
	"fmt"
	"time"

	log "github.com/shylinux/toolkits/logs"
)

const (
	StatusPrepare = iota
	StatusProcess
	StatusCancel
	StatusFinish
)

type Task struct {
	Arg interface{}
	CB  func(*Task) error

	ID     int64
	Error  interface{}
	Status int

	PrepareTime time.Time
	ProcessTime time.Time
	FinishTime  time.Time

	Ctx  context.Context
	work *Work
}

func (task *Task) Info() string {
	return fmt.Sprintf("poolID: %d workID: %d taskID: %d", task.work.pool.ID, task.work.ID, task.ID)
}
func (task *Task) Pool() *Pool {
	return task.work.pool
}
func (task *Task) Run() {
	log.Show("task", "task run", log.FileLine(task.CB, 3), "arg", task.Arg, "id", task.ID, "work", task.work.ID, "pool", task.work.pool.ID)
	defer log.Cost("task end: ", log.FileLine(task.CB, 3), " err: ", task.Error, " id: ", task.ID, " work: ", task.work.ID, " pool: ", task.work.pool.ID, " ")()

	task.Status = StatusProcess
	task.ProcessTime = time.Now()
	defer func() { task.FinishTime = time.Now() }()
	defer func() {
		if e := recover(); e != nil {
			for i := 1; i < 6; i++ {
				log.Show("task", "task err", e, "pos", log.FileLine(i, 3))
			}
			task.Status = StatusCancel
			task.Error = e
		}
	}()

	if e := task.CB(task); e != nil {
		task.Status = StatusCancel
		task.Error = e
	} else {
		task.Status = StatusFinish
	}
}
