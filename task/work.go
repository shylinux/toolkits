package task

import (
	"context"
	"github.com/shylinux/toolkits/log"
)

type Work struct {
	ID int64

	Ctx  context.Context
	pool *Pool
}

func (work *Work) Run() {
	defer log.Cost("work", " err: ", work.Ctx.Err(), " id: ", work.ID, " work: ", work.ID, " pool: ", work.pool.ID, " ")()
	defer func() {
		if e := recover(); e != nil {
			for i := 1; i < 6; i++ {
				log.Show("work", "work err", e, "pos", log.FileLine(i, 3))
			}
			work.pool.Add(1)
		}
	}()

	for {
		select {
		case task := <-work.pool.channel:
			task.work = work
			task.Ctx = context.WithValue(work.Ctx, "id", task.ID)
			log.Show("task", "task run", log.FileLine(task.CB, 3), "arg", task.Arg, "id", task.ID, "work", work.ID, "pool", work.pool.ID)
			task.Run()
		case <-work.Ctx.Done():
			return
		}
	}
}
