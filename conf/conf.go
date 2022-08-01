package conf

import (
	"context"
	"sync"

	kit "shylinux.com/x/toolkits"
)

type Any = interface{}

type Conf struct {
	data Any

	cancel context.CancelFunc
	ctx    context.Context
	wg     sync.WaitGroup

	sup *Conf
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
		case kit.Map:
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
func (conf *Conf) GetVal(key string, def ...Any) Any {
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
func (conf *Conf) Sub(key string) *Conf {
	ctx, cancel := context.WithCancel(conf.ctx)
	return &Conf{data: kit.Value(conf.data, key), cancel: cancel, ctx: ctx, sup: conf}
}
func New(data Any) *Conf {
	ctx, cancel := context.WithCancel(context.TODO())
	return &Conf{data: data, cancel: cancel, ctx: ctx}
}
