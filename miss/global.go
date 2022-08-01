package miss

import (
	"io"

	"shylinux.com/x/toolkits/conf"
	"shylinux.com/x/toolkits/file"
)

var miss = New(conf.Sub(MISS), file.NewDiskFile())

func Init(conf *conf.Conf, file file.File) { miss = New(conf, file) }

func Richs(prefix string, cache Map, raw Any, cb Any) (res Map) {
	return miss.Richs(prefix, cache, raw, cb)
}
func Rich(prefix string, cache Map, data Any) string {
	return miss.Rich(prefix, cache, data)
}
func Grow(prefix string, cache Map, data Any) int {
	return miss.Grow(prefix, cache, data)
}
func Grows(prefix string, cache Map, offend, limit int, match string, value string, cb Any) Map {
	return miss.Grows(prefix, cache, offend, limit, match, value, cb)
}
func OpenFile(p string) (io.ReadCloser, error) {
	return miss.file.OpenFile(p)
}
func CreateFile(p string) (io.WriteCloser, string, error) {
	return miss.file.CreateFile(p)
}
