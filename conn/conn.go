package conn

import (
	"net"

	kit "shylinux.com/x/toolkits"
)

const CONN = "conn"

type Conn struct {
	id int64

	client kit.Any
	net.Conn

	closed bool
	nwrite int
	nread  int

	pool *Pool
}

func (conn *Conn) Info() string {
	return kit.FormatShow(CONN, conn.id, POOL, conn.pool.id)
}
func (conn *Conn) Read(b []byte) (int, error) {
	n, e := conn.Conn.Read(b)
	conn.nread += n
	return n, e
}
func (conn *Conn) Write(b []byte) (int, error) {
	n, e := conn.Conn.Write(b)
	conn.nwrite += n
	return n, e
}
func (conn *Conn) Close() error {
	if !conn.closed {
		conn.pool.Put(conn)
	} else {
		conn.Conn.Close()
	}
	return nil
}
