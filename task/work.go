package task

import (
	"context"

	kit "shylinux.com/x/toolkits"
)

const WORK = "work"

type Work struct {
	id int64

	Logger func(...Any)

	pool *Pool
}

func (work *Work) Info() string {
	return kit.FormatShow(WORK, work.id, POOL, work.pool.id)
}
func (work *Work) Run(ctx context.Context) {
	// work.Logger("work add", work.Info())
	// defer logs.CostTime(func(d time.Duration) {
	// 	work.Logger("work end", "err", ctx.Err(), work.Info())
	// })()
	for {
		select {
		case task, ok := <-work.pool.channel:
			if !ok {
				return
			}
			task.work = work
			task.Run(context.WithValue(ctx, TASK, task.id))
		}
	}
}
