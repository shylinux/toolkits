package task

import (
	"shylinux.com/x/toolkits/conf"
)

var pool = New(conf.Sub(TASK))

func Init(conf *conf.Conf) { pool = New(conf) }

func Wait(params []interface{}, action func(*Task, *Lock) error) {
	pool.Wait(params, action)
}
func Put(param interface{}, action func(*Task) error) {
	pool.Put(param, action)
}
