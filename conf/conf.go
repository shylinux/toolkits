package conf

import kit "github.com/shylinux/toolkits"

type Conf struct {
	data interface{}
}

func (conf *Conf) Get(key string, def ...interface{}) interface{} {
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
func (conf *Conf) GetStr(key string, def ...string) string {
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

func New(data interface{}) *Conf {
	return &Conf{data: data}
}

var conf *Conf

func Init(file string) {
	if conf == nil {
		conf = New(nil)
	}
}
