package task

import (
	"shylinux.com/x/toolkits/conf"
)

var pool = New(conf.Sub(TASK))

func Init(conf *conf.Conf) { pool = New(conf) }

func WaitN(n int, action func(*Task, *Lock) error) {
	pool.WaitN(n, action)
}
func Wait(params []Any, action func(*Task, *Lock) error) {
	pool.Wait(params, action)
}
func Put(param Any, action func(*Task) error) {
	pool.Put(param, action)
}
