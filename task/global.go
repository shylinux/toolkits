package task

import (
	"shylinux.com/x/toolkits/conf"
)

var pool = New(conf.Sub("task"))

func Wait(arg []interface{}, cb func(*Task, *Lock) error) *Pool {
	return pool.Wait(arg, cb)
}
func Put(arg interface{}, cb func(*Task) error) *Task {
	return pool.Put(arg, cb)
}

func Init(conf *conf.Conf) { pool = New(conf) }
