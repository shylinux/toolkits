package conn

import (
	"fmt"
	"net"
)

type Conn struct {
	client interface{}
	net.Conn

	ID int64

	pool *Pool

	nread  int
	nwrite int
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
func (conn *Conn) Info() string {
	return fmt.Sprintf("connID: %d poolID: %d", conn.ID, conn.pool.ID)
}
func (conn *Conn) Close() error {
	conn.pool.Put(conn)
	return nil
}
