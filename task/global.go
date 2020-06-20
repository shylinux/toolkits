package task

import (
	"github.com/shylinux/toolkits/conf"
)

var pool = New(conf.Sub("task"))

func Sync(arg []interface{}, cb func(*Task, *Lock) error) *Pool {
	return pool.Sync(arg, cb)
}
func Put(arg interface{}, cb func(*Task) error) *Task {
	return pool.Put(arg, cb)
}

func Init(conf *conf.Conf) { pool = New(conf) }
