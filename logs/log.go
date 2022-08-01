package logs

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	kit "shylinux.com/x/toolkits"
	"shylinux.com/x/toolkits/conf"
	"shylinux.com/x/toolkits/file"
)

type Any = interface{}

const (
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
	DEBUG = "debug"
	SHOW  = "show"
	COST  = "cost"
)

const LOG = "log"

type Log struct {
	Output   string
	Filter   []string
	Prefix   []string
	Notice   []string
	Out, Err io.WriteCloser

	debug   bool
	color   bool
	colors  map[string]string
	channel chan string
	disable bool

	file file.File
	conf *conf.Conf
}

func (log *Log) filter(key string) bool {
	if log.debug {
		return false
	}
	for _, v := range log.Filter {
		if v == key {
			return true
		}
	}
	return false
}
func (log *Log) prefix(level string) (io.Writer, string) {
	if log.filter(level) {
		return nil, ""
	}

	list := []string{}
	for _, v := range log.Prefix {
		switch v {
		case "time":
			list = append(list, FmtTime(Now()))
		case "pid":
			list = append(list, FmtInt(os.Getpid()))
		case "level":
			list = append(list, level)
		}
	}
	return log.notice(level), strings.Join(list, " ")
}
func (log *Log) notice(key string) io.Writer {
	for _, k := range log.Notice {
		if k == key {
			return log.Err
		}
	}
	return log.Out
}

var trans = map[string]string{"red": "\033[31m", "green": "\033[32m", "yellow": "\033[33m"}

func (log *Log) output(level string, text string) bool {
	if log == nil || log.disable {
		return false
	}
	if output, prefix := log.prefix(level); output != nil {
		color, end := trans[log.colors[level]], ""
		if color != "" {
			end = "\033[0m"
		}
		if log.channel == nil {
			fmt.Fprintln(output, prefix, color, text, end)
		} else {
			log.channel <- fmt.Sprintln(prefix, color, text, end)
		}
	}
	return true
}

func (log *Log) Info(arg ...Any)  { log.output(INFO, Format("", arg...)) }
func (log *Log) Warn(arg ...Any)  { log.output(WARN, Format("", arg...)) }
func (log *Log) Error(arg ...Any) { log.output(ERROR, Format("", arg...)) }

func (log *Log) Infof(str string, arg ...Any)  { log.output(INFO, Format(str, arg...)) }
func (log *Log) Warnf(str string, arg ...Any)  { log.output(WARN, Format(str, arg...)) }
func (log *Log) Errorf(str string, arg ...Any) { log.output(ERROR, Format(str, arg...)) }

func (log *Log) Debugf(str string, arg ...Any) bool {
	return log.debug && log.output(DEBUG, Format(str, arg...))
}
func (log *Log) Debug(arg ...Any) bool {
	return log.debug && log.output(DEBUG, Format("", arg...))
}
func (log *Log) Show(arg ...Any) {
	list, meta := []Any{}, []Any{}
	for i := 1; i < len(arg); i++ {
		switch v := arg[i].(type) {
		case Meta:
			meta = append(meta, v)
		default:
			list = append(list, v)
		}
	}
	log.output(fmt.Sprint(arg[0]), Format(kit.FormatShow(list...), meta...))
}
func (log *Log) Cost(arg ...Any) func(...func() []Any) {
	begin := Now()
	return func(cbs ...func() []Any) {
		arg = append(arg, COST, FmtDuration(Now().Sub(begin)))
		for _, cb := range cbs {
			arg = append(arg, cb()...)
		}
		log.Show(append([]Any{COST}, arg...)...)
	}
}

func (log *Log) Logger(key string) func(...Any) {
	return func(arg ...Any) {
		log.Show(kit.List(key, arg, FileLineMeta(kit.FileLine(2, 3)))...)
	}
}

func New(conf *conf.Conf, file file.File) *Log {
	log := &Log{
		Output: conf.Get("name", "stderr"),
		Filter: conf.GetList("filter"),
		Prefix: conf.GetList("prefix", kit.Simple("time")),
		Notice: conf.GetList("notice"),
		Out:    os.Stderr, Err: os.Stderr,

		debug:  conf.GetBool("debug"),
		color:  conf.GetBool("color"),
		colors: conf.GetDict("colors"),

		file: file, conf: conf,
	}

	if f, e := file.AppendFile(log.Output + ".log"); e == nil {
		log.Out = f
	}
	if f, e := file.AppendFile(log.Output + ".err.log"); e == nil {
		log.Err = f
	}

	return log
	log.channel = make(chan string, 1024)
	conf.Daemon(LOG, func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				close(log.channel)
			case str, ok := <-log.channel:
				if ok {
					fmt.Fprint(log.Out, str)
				} else {
					return
				}
			}
		}
	})
	return log
}
