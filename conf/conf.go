package conf

import (
	"io/ioutil"
	"os"
	"strings"

	kit "github.com/shylinux/toolkits"
)

type Conf struct {
	data interface{}
}

func (conf *Conf) GetBool(key string, def ...bool) bool {
	if val := kit.Value(conf.data, key); val != nil {
		switch kit.Format(val) {
		case "true", "on", "1":
			return true
		default:
			return false
		}
	}
	for _, v := range def {
		if v == true {
			return v
		}
	}
	return false
}
func (conf *Conf) GetVal(key string, def ...interface{}) interface{} {
	if val := kit.Value(conf.data, key); val != nil {
		return val
	}
	for _, v := range def {
		if v != nil {
			return v
		}
	}
	return nil
}
func (conf *Conf) GetInt(key string, def ...int) int {
	if val := kit.Value(conf.data, key); val != nil {
		return kit.Int(val)
	}

	for _, v := range def {
		if v != 0 {
			return v
		}
	}
	return 0
}
func (conf *Conf) Get(key string, def ...string) string {
	if val := kit.Value(conf.data, key); val != nil {
		return kit.Format(val)
	}

	for _, v := range def {
		if v != "" {
			return v
		}
	}
	return ""
}
func (conf *Conf) GetList(key string, def ...[]string) []string {
	if val := kit.Value(conf.data, key); val != nil {
		return kit.Simple(val)
	}

	for _, v := range def {
		if len(v) > 0 {
			return v
		}
	}
	return nil
}
func (conf *Conf) GetDict(key string, def ...map[string]string) map[string]string {
	res := map[string]string{}
	if val := kit.Value(conf.data, key); val != nil {
		switch val := val.(type) {
		case map[string]string:
			return val
		case map[string]interface{}:
			for k, v := range val {
				res[k] = kit.Format(v)
			}
			return res
		}
	}

	for _, v := range def {
		if len(v) > 0 {
			return v
		}
	}
	return res
}

func New(data interface{}) *Conf {
	return &Conf{data: data}
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

var conf = &Conf{data: kit.Dict(
	"app.name", "demo",
	"log.name", "demo",
)}

func GetVal(key string, def ...interface{}) interface{} { return conf.GetVal(key, def...) }
func GetInt(key string, def ...int) int                 { return conf.GetInt(key, def...) }
func Get(key string, def ...string) string              { return conf.Get(key, def...) }

func Init(file string) {
	c, e := Open(file)
	if e != nil {
		panic(e)
	}
	conf = c
}
