package conn

import (
	"context"
	"errors"
	"net"
	"net/http"

	kit "shylinux.com/x/toolkits"
)

const HTTP = "http"

var (
	ErrTypeError = errors.New("type error")
	ErrUnknown   = errors.New("unknown error")
)

type HttpClient struct {
	*http.Client

	DoCount   int
	GetCount  int
	PostCount int
	ErrCount  int

	Logger func(...Any)

	conn *Conn
}

func (client *HttpClient) Info() string {
	return kit.FormatShow(CONN, client.conn.id, POOL, client.conn.pool.id)
}
func (client *HttpClient) Do(req *http.Request) (*http.Response, error) {
	res, err := client.Client.Do(req)
	if client.DoCount++; err != nil {
		client.ErrCount++
	}
	client.Logger(CONN, client.conn.id, "count", client.DoCount, "DO", req.URL, "err", err)
	return res, err
}
func (client *HttpClient) Get(url string) (*http.Response, error) {
	res, err := client.Client.Get(url)
	if client.GetCount++; err != nil {
		client.ErrCount++
	}
	client.Logger(CONN, client.conn.id, "count", client.GetCount, "GET", url, "err", err)
	return res, err
}
func (client *HttpClient) NWrite() int {
	return client.conn.nwrite
}
func (client *HttpClient) NRead() int {
	return client.conn.nread
}
func (client *HttpClient) Release() {
	client.conn.Close()
}

func (pool *Pool) GetHttp(ctx context.Context) (*HttpClient, error) {
	if conn := pool.Get(ctx); conn != nil {
		switch client := conn.client.(type) {
		case nil:
		case *HttpClient:
			return client, nil
		default:
			return nil, ErrTypeError
		}
		client := &HttpClient{Client: &http.Client{Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) { return conn, nil },
		}}, conn: conn, Logger: pool.Logger}
		conn.client = client
		return client, nil
	}
	return nil, ErrUnknown
}
