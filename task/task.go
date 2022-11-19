package task

import (
	"context"
	"strings"
	"time"

	kit "shylinux.com/x/toolkits"
	"shylinux.com/x/toolkits/logs"
)

const (
	StatusPrepare = iota
	StatusProcess
	StatusCancel
	StatusFinish
)

const TASK = "task"

type Task struct {
	id int64

	Action func(*Task) error
	Params Any

	Error  Any
	Status int

	PrepareTime time.Time
	ProcessTime time.Time
	FinishTime  time.Time

	Logger func(...Any)

	ctx  context.Context
	work *Work
	pool *Pool
}

func (task *Task) Pool() *Pool {
	return task.pool
}
func (task *Task) Context() context.Context {
	return task.ctx
}
func (task *Task) Info() string {
	if task.work == nil {
		return kit.FormatShow(TASK, task.id, WORK, 0, POOL, task.pool.id)
	}
	return kit.FormatShow(TASK, task.id, WORK, task.work.id, POOL, task.pool.id)
}

func (task *Task) Run(ctx context.Context) {
	task.Logger("task run", logs.FileLine(task.Action), "params", task.Params, task.Info())
	defer logs.CostTime(func(d time.Duration) {
		task.Logger("task end", logs.FileLine(task.Action), "params", task.Params, "cost", logs.FmtDuration(d), "err", task.Error, task.Info())
	})()

	task.Status, task.ProcessTime, task.ctx = StatusProcess, time.Now(), ctx
	defer func() {
		if e := recover(); e != nil {
			list := []string{}
			for i := 1; i < 10; i++ {
				list = append(list, logs.FileLine(i))
			}
			task.Logger("task err", e, "stack", "\n", strings.Join(list, "\n"), "\n")
			task.Status, task.Error = StatusCancel, e
		}
		task.FinishTime = logs.Now()
	}()

	if e := task.Action(task); e != nil {
		task.Status, task.Error = StatusCancel, e
	} else {
		task.Status = StatusFinish
	}
}
