package logs

import (
	"io"
	"os"

	"shylinux.com/x/toolkits/conf"
	"shylinux.com/x/toolkits/file"
)

var log = New(conf.Sub(LOG), file.NewDiskFile())

func Init(conf *conf.Conf, file file.File) { log = New(conf, file) }

func Disable(status bool) { log.disable = status }

func Info(arg ...Any)  { log.Info(fileline(arg)...) }
func Warn(arg ...Any)  { log.Warn(fileline(arg)...) }
func Error(arg ...Any) { log.Error(fileline(arg)...) }
func Debug(arg ...Any) { log.Debug(fileline(arg)...) }
func Show(arg ...Any)  { log.Show(fileline(arg)...) }
func Cost(arg ...Any) func(...func() []Any) {
	return log.Cost(fileline(arg)...)
}
func Infof(str string, arg ...Any)   { log.Infof(str, fileline(arg)...) }
func Logger(key string) func(...Any) { return log.Logger(key) }

func CreateFile(p string) (io.WriteCloser, string, error) {
	if log.disable {
		return nil, "", os.ErrInvalid
	}
	return log.file.CreateFile(p)
}
func fileline(arg []Any) []Any { return append(arg, FileLineMeta(FileLine(3, 3))) }
