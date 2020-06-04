package task

import (
	"sync/atomic"
	"time"

	"github.com/shylinux/toolkits/conf"
)

type Pool struct {
	TaskID int64
}

func (pool *Pool) Run(arg interface{}, cb func(*Task) error) *Task {
	atomic.AddInt64(&pool.TaskID, 1)

	// log.Show("task", "add", log.FileLine(cb, 3), "id", TaskID, "arg", arg)
	task := &Task{ID: pool.TaskID, Arg: arg, CB: cb, BeginTime: time.Now()}
	go func() {
		task.Run()
	}()
	return task
}

func New(conf *conf.Conf) *Pool {
	return &Pool{}
}
