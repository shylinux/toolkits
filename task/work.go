package task

import (
	"context"

	log "shylinux.com/x/toolkits/logs"
)

type Work struct {
	ID int64

	Ctx  context.Context
	pool *Pool
}

func (work *Work) Run() {
	defer log.Cost("work", " err: ", work.Ctx.Err(), " id: ", work.ID, " pool: ", work.pool.ID, " ")()

	for {
		select {
		case <-work.Ctx.Done():
			return
		case task := <-work.pool.channel:
			task.work = work
			task.Ctx = context.WithValue(work.Ctx, "id", task.ID)
			task.Run()
		}
	}
}
