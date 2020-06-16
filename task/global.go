package task

import (
	"github.com/shylinux/toolkits/conf"
)

var pool = New(conf.New(nil), 10)

func Sync(arg []interface{}, cb func(*Task) error) *Pool { return pool.Sync(arg, cb) }
func Put(arg interface{}, cb func(*Task) error) *Task    { return pool.Put(arg, cb) }
