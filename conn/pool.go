package conn

import (
	"context"
	"net"
	"sync"
	"sync/atomic"

	"github.com/shylinux/toolkits/conf"
	log "github.com/shylinux/toolkits/logs"
)

type Pool struct {
	addrs []string
	limit int64
	retry int

	mu      sync.Mutex
	channel chan *Conn
	connID  int64

	ID int64
	sync.Pool
}

func (pool *Pool) Get(ctx context.Context) *Conn {
	pool.add()

	for {
		select {
		case <-ctx.Done():
			log.Show("conn", "conn err", ctx.Err(), "pool", pool.ID)
			return nil
		case c := <-pool.channel:
			log.Show("conn", "conn get", c.LocalAddr(), "id", c.ID, "pool", pool.ID)
			return c
		}
	}
	return nil
}
func (pool *Pool) Put(c *Conn) {
	log.Show("conn", "conn put", c.LocalAddr(), "id", c.ID, "pool", pool.ID)
	pool.channel <- c
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

func New(conf *conf.Conf, addrs []string, limit int64, retry int) *Pool {
	pool := &Pool{
		addrs: addrs, limit: limit, retry: retry,
		channel: make(chan *Conn, limit),
		ID:      atomic.AddInt64(&poolID, 1),
	}
	pool.Pool = sync.Pool{New: func() interface{} {
		id := atomic.AddInt64(&pool.connID, 1)
		for i := 0; i < pool.retry; i++ {
			if c, e := net.Dial("tcp", addrs[id%int64(len(addrs))]); e == nil {
				log.Show("conn", "conn add", c.LocalAddr(), "id", id, "pool", pool.ID)
				return &Conn{ID: id, Conn: c, pool: pool}
			} else {
				log.Warn("dial", addrs, i, e)
			}
		}
		return nil
	}}
	return pool
}
