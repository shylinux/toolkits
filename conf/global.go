package conf

import (
	kit "shylinux.com/x/toolkits"
)

var conf = New(kit.Dict(
	"app", kit.Dict(
		"name", "demo",
	),
	"log", kit.Dict(
		"prefix", []string{"time"},
	),
	"miss", kit.Dict(
		"store", "var/data",
		"fsize", "200000",
		"limit", "110",
		"least", "10",
	),
	"task", kit.Dict(
		"limit", 50,
	),
	"conn", kit.Dict(
		"limit", 30,
		"retry", 3,
	),
))

func GetVal(key string, def ...interface{}) interface{} {
	return conf.GetVal(key, def...)
}
func GetInt(key string, def ...int) int {
	return conf.GetInt(key, def...)
}
func Get(key string, def ...string) string {
	return conf.Get(key, def...)
}
func Sub(key string) *Conf {
	return &Conf{data: kit.Value(conf.data, key)}
}

func Init(file string) {
	c, e := Open(file)
	if e != nil {
		panic(e)
	}
	conf = c
}
