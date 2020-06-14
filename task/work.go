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
	defer log.Cost("work", " id: ", work.ID, " err: ", work.Ctx.Err(), " ")()
	defer func() {
		if e := recover(); e != nil {
			for i := 1; i < 6; i++ {
				log.Show("work", "err", e, "pos", log.FileLine(i, 3))
			}
			work.pool.Add(1)
		}
	}()

	for {
		select {
		case task := <-work.pool.channel:
			task.work = work
			task.Ctx = context.WithValue(work.Ctx, "id", task.ID)
			log.Show("task", "run", log.FileLine(task.CB, 3), "id", task.ID, "arg", task.Arg, "work", work.ID, "pool", work.pool.ID)
			task.Run()
		case <-work.Ctx.Done():
			return
		}
	}
}
