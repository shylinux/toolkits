package log

import (
	"github.com/shylinux/toolkits/conf"
)

var log, _ = Open(conf.New(nil))

func init() { log.stack = STACK + 1 }

func Info(arg ...interface{})  { log.Info(arg...) }
func Warn(arg ...interface{})  { log.Warn(arg...) }
func Error(arg ...interface{}) { log.Error(arg...) }

func Show(arg ...interface{}) { log.Show(arg...) }
func Cost(arg ...interface{}) func(...func() []interface{}) {
	return log.Cost(arg...)
}

func Init(conf *conf.Conf) {
	l, e := Open(conf)
	if e != nil {
		panic(e)
	}
	l.stack = STACK + 1
	log = l
}
