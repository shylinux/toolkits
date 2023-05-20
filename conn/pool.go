package conn

import (
	"context"
	"net"
	"sync"
	"sync/atomic"

	"shylinux.com/x/toolkits/conf"
	"shylinux.com/x/toolkits/logs"
	"shylinux.com/x/toolkits/task"
)

type Any = interface{}

const (
	CONN_ERR = "conn err"
	CONN_ADD = "conn add"
	CONN_GET = "conn get"
	CONN_PUT = "conn put"
	CONN_END = "conn end"
)
const POOL = "pool"

type Pool struct {
	id int64
	mu task.Lock

	nconn   int64
	maxconn int64
	channel chan *Conn

	target []string
	retry  int

	Logger func(...Any)

	sync.Pool
	conf *conf.Conf
}

func (pool *Pool) Get(ctx context.Context) *Conn {
	pool.add()
	for {
		select {
		case <-ctx.Done():
			pool.Logger(CONN_ERR, ctx.Err(), POOL, pool.id)
			return nil
		case c := <-pool.channel:
			pool.Logger(CONN_GET, c.LocalAddr(), c.Info())
			return c
		}
	}
	return nil
}
func (pool *Pool) add() {
	defer pool.mu.Lock()()
	if len(pool.channel) == 0 && pool.nconn < pool.maxconn {
		if c, ok := pool.Pool.Get().(*Conn); ok {
			pool.Logger(CONN_ADD, c.LocalAddr(), c.Info())
			pool.channel <- c
		}
	}
}
func (pool *Pool) Put(c *Conn) {
	pool.Logger(CONN_PUT, c.LocalAddr(), c.Info())
	pool.channel <- c
}

func (pool *Pool) Close() {
	close(pool.channel)
	for c := range pool.channel {
		pool.Logger(CONN_END, c.LocalAddr(), c.Info())
		c.closed = true
		c.Close()
	}
}

var npool int64

func New(conf *conf.Conf, maxconn int64, target []string, retry int) *Pool {
	pool := &Pool{id: atomic.AddInt64(&npool, 1),
		maxconn: maxconn, channel: make(chan *Conn, maxconn),
		target: target, retry: retry, conf: conf,
		Logger: logs.Logger(CONN),
	}
	pool.Pool = sync.Pool{New: func() Any {
		id := atomic.AddInt64(&pool.nconn, 1)
		for i := 0; i < pool.retry; i++ {
			if c, e := net.Dial("tcp", target[id%int64(len(target))]); e == nil {
				return &Conn{id: id, Conn: c, pool: pool}
			} else {
				logs.Warn(CONN, " dial: ", target, " ", e)
			}
		}
		return nil
	}}
	return pool
}
