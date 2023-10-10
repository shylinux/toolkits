package task

import (
	"shylinux.com/x/toolkits/conf"
)

var pool = New(conf.Sub(TASK))

func Init(conf *conf.Conf) { pool = New(conf) }

func WaitN(meta Meta, n int, action func(*Task, *Lock) error) {
	pool.WaitN(meta, n, action)
}
func Wait(meta Meta, params []Any, action func(*Task, *Lock) error) {
	pool.Wait(meta, params, action)
}
func Put(meta Meta, param Any, action func(*Task)) {
	pool.Put(meta, param, action)
}
