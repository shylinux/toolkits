package redis

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/shylinux/toolkits/task"
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

	Cost time.Duration

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
	return ""
}

func _SET(i int64) []interface{} {
	return []interface{}{fmt.Sprintf("hi%d", i), "hello"}
}
func _GET(i int64) []interface{} {
	return []interface{}{fmt.Sprintf("hi%d", i)}
}

var trans = map[string]func(i int64) []interface{}{
	"GET": _GET,
	"SET": _SET,
}

func Redis(nconn, nreq int64, hosts []string, cmds []string, check func(string, []interface{}, interface{})) (*Stat, error) {
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

	// 连接池
	rp := redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", hosts[0])
	}, 10)

	// 协程池
	list := []interface{}{}
	for i := int64(0); i < nconn; i++ {
		list = append(list, i)
	}

	task.Sync(list, func(task *task.Task, lock *task.Lock) error {
		// 请求汇总
		var nerr, nok int64
		defer func() {
			atomic.AddInt64(&s.NReq, nreq)
			atomic.AddInt64(&s.NErr, nerr)
			atomic.AddInt64(&s.NOK, nok)
		}()

		conn := rp.Get()
		defer conn.Close()

		cmd := cmds[0]
		method := trans[cmd]

		for i := int64(0); i < nreq; i++ {
			func() {
				// 请求耗时
				begin := time.Now()
				defer func() {
					d := time.Now().Sub(begin)
					s.mu.Lock()
					defer s.mu.Unlock()
					s.Cost += d
				}()

				arg := method(i)
				if reply, err := conn.Do(cmd, arg...); err != nil {
					// 请求失败
					nerr++
				} else {
					// 请求成功
					if check != nil {
						check(cmd, arg, reply)
					}
					nok++
				}
			}()
		}
		return nil
	})
	return s, nil
}
