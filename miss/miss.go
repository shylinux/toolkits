package miss

import (
	"github.com/shylinux/toolkits/conf"
	"time"
)

type Miss struct {
	store string
	fsize string
	limit string
	least string
}

func New(conf *conf.Conf) *Miss {
	return &Miss{
		store: conf.Get("store", "var/data"),
		fsize: conf.Get("fsize", "200000"),
		limit: conf.Get("limit", "110"),
		least: conf.Get("least", "10"),
	}
}
func (miss *Miss) now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
