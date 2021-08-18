package task

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"shylinux.com/x/toolkits/conf"
	log "shylinux.com/x/toolkits/logs"
)

type Pool struct {
	limit int64

	ID int64
	mu sync.Mutex

	workID  int64
	taskID  int64
	channel chan *Task

	Ctx    context.Context
	cancel context.CancelFunc
}

func (pool *Pool) Wait(args []interface{}, cb func(*Task, *Lock) error) *Pool {
	var lock Lock

	w := &sync.WaitGroup{}
	for _, arg := range args {
		w.Add(1)
		pool.Put(arg, func(task *Task) error {
			defer w.Done()
			return cb(task, &lock)
		})
	}
	w.Wait()
	return pool
}
func (pool *Pool) Put(arg interface{}, cb func(*Task) error) *Task {
	id := atomic.AddInt64(&pool.taskID, 1)
	// log.Show("task", "task put", log.FileLine(cb, 3), "arg", arg, "id", id, "pool", pool.ID)
	task := &Task{ID: id, Arg: arg, CB: cb, PrepareTime: time.Now()}

	if pool.channel <- task; len(pool.channel) > 0 && pool.workID < pool.limit {
		pool.Add(1)
	}
	return task
}
func (pool *Pool) Add(count int) *Pool {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	for i := 0; i < count && pool.workID < pool.limit; i++ {
		id := atomic.AddInt64(&pool.workID, 1)
		ctx := context.WithValue(pool.Ctx, "id", id)
		work := &Work{ID: id, pool: pool, Ctx: ctx}
		log.Show("work", "work add", log.FileLine(work.Run, 3), "id", id, "pool", pool.ID)
		go work.Run()
	}
	return pool
}
func (pool *Pool) Close() { pool.cancel() }

var poolID int64

func New(conf *conf.Conf) *Pool {
	id := atomic.AddInt64(&poolID, 1)
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{limit: int64(conf.GetInt("limit", 10)), ID: id, channel: make(chan *Task, 1024), Ctx: ctx, cancel: cancel}
	return p
}
