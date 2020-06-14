package conn

import (
	"context"
	"net"
	"sync"
	"sync/atomic"

	"github.com/shylinux/toolkits/conf"
	"github.com/shylinux/toolkits/log"
)

type Pool struct {
	limit int64
	retry int
	addrs []string

	ID int64
	mu sync.Mutex

	connID  int64
	channel chan *Conn

	sync.Pool
}

func (pool *Pool) Put(c *Conn) {
	log.Show("conn", "put", c.LocalAddr(), "id", c.ID)
	pool.channel <- c
}
func (pool *Pool) Get(ctx context.Context) *Conn {
	pool.add()

	for {
		select {
		case <-ctx.Done():
			log.Show("conn", "err", ctx.Err())
			return nil
		case c := <-pool.channel:
			log.Show("conn", "get", c.LocalAddr(), "id", c.ID)
			return c
		}
	}
	return nil
}
func (pool *Pool) add() {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if len(pool.channel) == 0 && pool.connID < pool.limit {
		if c, ok := pool.Pool.Get().(*Conn); ok {
			pool.channel <- c
		}
	}
}

var poolID int64

func New(conf *conf.Conf, addrs []string, limit int64) *Pool {
	pool := &Pool{
		channel: make(chan *Conn, limit),
		limit:   limit, retry: 3, addrs: addrs,
		ID: atomic.AddInt64(&poolID, 1),
	}
	pool.Pool = sync.Pool{New: func() interface{} {
		id := atomic.AddInt64(&pool.connID, 1)
		for i := 0; i < pool.retry; i++ {
			if c, e := net.Dial("tcp", addrs[id%int64(len(addrs))]); e == nil {
				log.Show("conn", "add", c.LocalAddr(), "id", id, "limit", limit)
				return &Conn{ID: id, Conn: c, pool: pool}
			} else {
				log.Warn("dial", addrs, i, e)
			}
		}
		return nil
	}}
	return pool
}
