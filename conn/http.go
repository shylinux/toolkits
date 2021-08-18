package conn

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"shylinux.com/x/toolkits/logs"
)

var (
	ErrTypeError = errors.New("type error")
	ErrUnknown   = errors.New("unknown error")
)

type HttpClient struct {
	*http.Client

	PostCount int64
	GetCount  int64
	DoCount   int64
	ErrCount  int64

	conn *Conn
}

func (hc *HttpClient) Info() string {
	return fmt.Sprintf("connID: %d poolID: %d", hc.conn.ID, hc.conn.pool.ID)
}
func (hc *HttpClient) Get(url string) (*http.Response, error) {
	res, err := hc.Client.Get(url)
	if hc.GetCount++; err != nil {
		hc.ErrCount++
	}
	log.Show("http", "conn", hc.conn.ID, "count", hc.GetCount, "GET", url, "err", err)
	return res, err
}
func (hc *HttpClient) Do(req *http.Request) (*http.Response, error) {
	res, err := hc.Client.Do(req)
	if hc.DoCount++; err != nil {
		hc.ErrCount++
	}
	log.Show("http", "conn", hc.conn.ID, "count", hc.DoCount, "DO", req.URL, "err", err)
	return res, err
}
func (hc *HttpClient) NRead() int64 {
	return int64(hc.conn.nread)
}
func (hc *HttpClient) NWrite() int64 {
	return int64(hc.conn.nwrite)
}
func (hc *HttpClient) Close() {
	hc.conn.Close()
}

func (pool *Pool) GetHttp(ctx context.Context) (*HttpClient, error) {
	if c := pool.Get(ctx); c != nil {
		switch hc := c.client.(type) {
		case *HttpClient:
			return hc, nil
		case nil:
			client := &HttpClient{Client: &http.Client{Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return c, nil
				},
			}}, conn: c}
			c.client = client
			return client, nil
		default:
			return nil, ErrTypeError
		}
	}
	return nil, ErrUnknown
}
