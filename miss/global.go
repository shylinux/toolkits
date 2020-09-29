package miss

import (
	"github.com/shylinux/toolkits/conf"
)

var miss = New(conf.Sub("miss"))

func Init(conf *conf.Conf) { miss = New(conf) }

func Richs(prefix string, cache map[string]interface{}, raw interface{}, cb interface{}) (res map[string]interface{}) {
	return miss.Richs(prefix, cache, raw, cb)
}
func Rich(prefix string, cache map[string]interface{}, data interface{}) string {
	return miss.Rich(prefix, cache, data)
}
func Grow(prefix string, cache map[string]interface{}, data interface{}) int {
	return miss.Grow(prefix, cache, data)
}
func Grows(prefix string, cache map[string]interface{}, offend, limit int, match string, value string, cb interface{}) map[string]interface{} {
	return miss.Grows(prefix, cache, offend, limit, match, value, cb)
}
