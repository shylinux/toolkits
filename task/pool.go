package task

import (
	"context"
	"sync"
	"sync/atomic"

	"shylinux.com/x/toolkits/conf"
	"shylinux.com/x/toolkits/logs"
)

type Any = interface{}

const POOL = "pool"

type Pool struct {
	id int64
	mu Lock

	ntask   int64
	nwork   int64
	maxwork int64
	channel chan *Task

	Logger func(...Any)

	closed bool
	cancel context.CancelFunc
	ctx    context.Context
	conf   *conf.Conf
}

func (pool *Pool) WaitN(n int, action func(*Task, *Lock) error) {
	args := []Any{}
	for i := 1; i <= n; i++ {
		args = append(args, i)
	}
	pool.Wait(args, action)
}
func (pool *Pool) Wait(args []Any, action func(*Task, *Lock) error) {
	wg, lock := &sync.WaitGroup{}, &Lock{}
	defer wg.Wait()

	for _, arg := range args {
		wg.Add(1)
		pool.Put(arg, func(task *Task) error {
			defer wg.Done()
			return action(task, lock)
		})
	}
}
func (pool *Pool) Put(params Any, action func(*Task) error) {
	if pool.closed {
		return
	}
	task := &Task{id: atomic.AddInt64(&pool.ntask, 1), Action: action, Params: params, PrepareTime: logs.Now(), pool: pool, Logger: pool.Logger}
	pool.Logger("task put", logs.FileLine(action), "params", task.Params, task.Info())
	pool.channel <- task
	pool.add()
}
func (pool *Pool) add() {
	defer pool.mu.Lock()()
	if len(pool.channel) > 0 && pool.nwork < pool.maxwork {
		pool.conf.Daemon("task work", func(_ context.Context) {
			work := &Work{id: atomic.AddInt64(&pool.nwork, 1), pool: pool, Logger: pool.Logger}
			work.Run(context.WithValue(pool.ctx, WORK, work.id))
		})
	}
}
func (pool *Pool) Close() {
	pool.closed = true
	close(pool.channel)
	pool.cancel()
}

var npool int64

func New(conf *conf.Conf) *Pool {
	ctx, cancel := context.WithCancel(conf.Context())
	pool := &Pool{id: atomic.AddInt64(&npool, 1),
		maxwork: int64(conf.GetInt("maxwork", 30)),
		channel: make(chan *Task, conf.GetInt("maxtask", 300)),
		cancel:  cancel, ctx: ctx, conf: conf,
		Logger: logs.Logger(TASK),
	}
	conf.Daemon("task pool", func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				pool.closed = true
				close(pool.channel)
				return
			}
		}
	})
	return pool
}
