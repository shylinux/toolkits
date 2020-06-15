package task

import (
	"github.com/shylinux/toolkits/conf"
)

var pool = New(conf.New(nil), 10)

func Put(arg interface{}, cb func(*Task) error) *Task { return pool.Put(arg, cb) }
