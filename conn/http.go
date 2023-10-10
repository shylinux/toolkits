package conn

import (
	"context"
	"errors"
	"net"
	"net/http"

	kit "shylinux.com/x/toolkits"
	"shylinux.com/x/toolkits/logs"
	"shylinux.com/x/toolkits/task"
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
func (client *HttpClient) Do(meta task.Meta, req *http.Request) (*http.Response, error) {
	res, err := client.Client.Do(req)
	if client.DoCount++; err != nil {
		client.ErrCount++
	}
	client.Logger("HTTP DO", meta.FileLine, "id", client.conn.id, "count", client.DoCount, "url", req.URL, "err", err, logs.PrefixMeta(meta.Prefix))
	return res, err
}
func (client *HttpClient) Get(meta task.Meta, url string) (*http.Response, error) {
	res, err := client.Client.Get(url)
	if client.GetCount++; err != nil {
		client.ErrCount++
	}
	client.Logger("HTTP GET", meta.FileLine, "id", client.conn.id, "count", client.GetCount, "url", url, "err", err, logs.PrefixMeta(meta.Prefix))
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

func (pool *Pool) GetHttp(meta task.Meta, ctx context.Context) (*HttpClient, error) {
	if conn := pool.Get(meta, ctx); conn != nil {
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
