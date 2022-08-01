package bench

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	kit "shylinux.com/x/toolkits"
	"shylinux.com/x/toolkits/conf"
	"shylinux.com/x/toolkits/conn"
	"shylinux.com/x/toolkits/file"
	"shylinux.com/x/toolkits/logs"
	"shylinux.com/x/toolkits/task"
)

type Stat struct {
	NOK  int64
	NErr int64
	NReq int64

	NRead  int64
	NWrite int64

	BeginTime time.Time
	EndTime   time.Time

	List map[int64]int64
	mu   sync.Mutex

	Up   float64
	Down float64
	QPS  float64
	AVG  float64

	Max int64
	Min int64
	Avg int64
	Sum map[int64]int64

	Wait func()
}

func _span(sum, p, num int64) int64 {
	if sum*100/num < 2 {
		return 0
	}
	if sum*100/num > 98 {
		return 0
	}
	if sum > (p+1)/100*num {
		return 5
	}
	return 0
}

func (s *Stat) Show() string {
	list := []string{}
	list = append(list, fmt.Sprintf("beginTime: %s duration: %s\n", s.BeginTime.Format("2006-01-02 15:04:05"), s.EndTime.Sub(s.BeginTime)))
	list = append(list, fmt.Sprintf("QPS: %.2f n/s AVG: %.2fms Up: %.2f MB/s Down: %.2f MB/s\n", s.QPS, s.AVG,
		s.Up/float64((1<<20)), s.Down/float64((1<<20))))
	list = append(list, fmt.Sprintf("nreq: %d nerr: %d nok: %d\n", s.NReq, s.NErr, s.NOK))
	list = append(list, fmt.Sprintf("nread: %d nwrite: %d\n", s.NRead, s.NWrite))

	var max, min, sum, avg, num int64
	min = 100000000
	for i, v := range s.List {
		if i > max {
			max = i
		}
		if i < min && i != 0 {
			min = i
		}
		sum += i * v
		num += v
	}
	avg = sum / num
	list = append(list, fmt.Sprintf("max: %dms min: %dms avg: %dms\n", max, min, avg))
	s.Max, s.Min, s.Avg = max, min, avg

	last := int64(0)
	res := map[int64]int64{}
	for i, p := int64(1), int64(0); i < max+1; i += 1 {
		if last == 0 {
			last = i
		}
		res[last] += s.List[i]
		if n := _span(res[last], p, num); n > 0 && i != last {
			res[i] += res[last]
			last = i
			p += n
		}
	}

	last = 0
	for i := int64(0); i < max+1; i += 1 {
		if res[i] > 0 {
			if last > 0 {
				list = append(list, fmt.Sprintf("%dms: %3d %d%%\n", last, res[last], res[last]*100/num))
			}
			last = i
		}
	}
	s.Sum = res

	return strings.Join(list, "")
}

func HTTP(nconn, nreq int64, req []*http.Request, check func(*http.Request, *http.Response)) (*Stat, error) {
	// 输出流
	f := file.NewVoidFile()
	defer f.Close()

	// 日志流
	l := logs.New(conf.Sub(logs.LOG), file.NewDiskFile())

	// 协程池
	p := task.New(conf.Sub(task.TASK))
	p.Logger = l.Logger(task.TASK)
	defer p.Close()

	// 连接池
	hosts := []string{}
	for _, v := range req {
		hosts = append(hosts, v.Host)
	}
	c := conn.New(conf.Sub(conn.CONN), nconn, hosts, 3)
	c.Logger = l.Logger(task.TASK)
	defer c.Close()

	// 请求统计
	s := &Stat{BeginTime: time.Now(), List: make(map[int64]int64, 100)}
	defer func() {
		if s.EndTime = time.Now(); s.BeginTime != s.EndTime {
			d := float64(s.EndTime.Sub(s.BeginTime)) / float64(time.Second)

			s.QPS = float64(s.NReq) / d
			s.AVG = float64(s.EndTime.Sub(s.BeginTime)) / float64(time.Millisecond) / float64(s.NReq)
			s.Down = float64(s.NRead) / d
			s.Up = float64(s.NWrite) / d
		}
	}()

	p.WaitN(int(nconn), func(task *task.Task, lock *task.Lock) error {
		order := kit.Int(task.Params)
		hc, e := c.GetHttp(task.Context())
		if e != nil {
			return e
		}
		defer hc.Release()

		var nerr, nok int64
		defer func() { // 请求计数
			atomic.AddInt64(&s.NReq, nreq)
			atomic.AddInt64(&s.NErr, nerr)
			atomic.AddInt64(&s.NOK, nok)
			atomic.AddInt64(&s.NRead, int64(hc.NRead()))
			atomic.AddInt64(&s.NWrite, int64(hc.NWrite()))
		}()

		for i := int64(0); i < nreq; i++ {
			func() {
				begin := time.Now()
				defer func() { // 请求耗时
					d := int64(time.Now().Sub(begin) / time.Millisecond)
					defer lock.Lock()()
					s.List[d]++
				}()

				req := req[nreq%int64(len(req))]
				if res, err := hc.Do(req); err != nil { // 请求失败
					nerr++
				} else { // 请求完成
					defer res.Body.Close()
					if check == nil {
						f, _, e := f.CreateFile(kit.Format("var/bench/some/%s-%s-res-body.txt", order, i))
						if e != nil {
							return
						}
						defer f.Close()

						io.Copy(f, res.Body)
					} else {
						check(req, res)
					}

					switch res.StatusCode {
					case http.StatusOK: // 请求成功
						nok++
					}
				}
			}()
		}
		return nil
	})
	return s, nil
}
