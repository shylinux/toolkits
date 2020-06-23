package task

import (
	"context"

	log "github.com/shylinux/toolkits/logs"
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
			task.Run()
		case <-work.Ctx.Done():
			return
		}
	}
}
