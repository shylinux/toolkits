package conf

import (
	"io/ioutil"
	"os"
	"strings"

	kit "shylinux.com/x/toolkits"
)

var conf = New(kit.Dict(
	"app", kit.Dict(
		"name", "demo",
	),
	"log", kit.Dict(
		"prefix", []string{"time"},
	),
	"file", kit.Dict(),
	"miss", kit.Dict(
		"store", "var/data",
		"fsize", "200000",
		"limit", "110",
		"least", "10",
	),
	"task", kit.Dict(
		"maxwork", 250,
		"maxtask", 1000,
	),
	"conn", kit.Dict(
		"limit", 30,
		"retry", 3,
	),
))

func Init(file string) {
	c, e := Open(file)
	if e != nil {
		panic(e)
	}
	conf = c
}
func Sub(key string) *Conf {
	return conf.Sub(key)
}
func Wait()  { conf.Wait() }
func Close() { conf.Close() }
func Open(file string) (*Conf, error) {
	f, e := os.Open(file)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	b, e := ioutil.ReadAll(f)
	if e != nil {
		return nil, e
	}
	return Parse(string(b))
}
func Parse(text string) (*Conf, error) {
	prefix, list := false, []string{}
	for _, v := range kit.Split(text, "\n", "\n") {
		if strings.TrimSpace(v) == "" {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(v), "#") {
			continue
		}
		if !prefix && strings.HasPrefix(strings.TrimSpace(v), "{") {
			prefix = true
		}
		list = append(list, kit.Split(v, "\t :=,;\n", "{[]}")...)
	}
	if !prefix {
		list = append([]string{"{"}, list...)
	}
	return New(kit.Parse(nil, "", list...)), nil
}
