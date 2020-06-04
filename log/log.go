package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shylinux/toolkits/conf"
)

type Log struct {
	Output   string
	Prefix   []string
	Filter   []string
	Notice   []string
	Out, Err io.Writer

	column int
	debug  bool
	color  bool
	colors map[string]string
	*conf.Conf
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
func (log *Log) prefix(level string, arg ...interface{}) (io.Writer, string) {
	if log.filter(level) {
		return nil, ""
	}

	file, list := log.Out, []string{}
	for _, k := range log.Notice {
		if k == level {
			file = log.Err
		}
	}

	for _, v := range log.Prefix {
		switch v {
		case "time":
			list = append(list, time.Now().Format("2006-01-02 15:04:05"))
		case "level":
			list = append(list, level)

		case "pid":
			list = append(list, strconv.FormatInt(int64(os.Getpid()), 10))

		case "fileline":
			_, file, line, _ := runtime.Caller(3)
			ls := strings.Split(file, "/")
			if len(ls) > log.column {
				ls = ls[len(ls)-log.column:]
			}
			list = append(list, fmt.Sprintf("%s:%d", strings.Join(ls, "/"), line))
		}
	}
	return file, strings.Join(list, " ")
}

var trans = map[string]string{
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
}

func (log *Log) output(level string, arg ...interface{}) {
	if log == nil {
		return
	}
	if output, prefix := log.prefix(level, arg...); output != nil {
		color, end := trans[log.colors[level]], ""
		if color != "" {
			end = "\033[0m"
		}
		fmt.Fprint(output, prefix, " ", color, fmt.Sprint(arg...), end, "\n")
	}
}

func (log *Log) Info(arg ...interface{})  { log.output("info", arg...) }
func (log *Log) Warn(arg ...interface{})  { log.output("warn", arg...) }
func (log *Log) Error(arg ...interface{}) { log.output("error", arg...) }

func (log *Log) Infof(str string, arg ...interface{})  { log.output("info", fmt.Sprintf(str, arg...)) }
func (log *Log) Warnf(str string, arg ...interface{})  { log.output("warn", fmt.Sprintf(str, arg...)) }
func (log *Log) Errorf(str string, arg ...interface{}) { log.output("error", fmt.Sprintf(str, arg...)) }

func (log *Log) Debugf(str string, arg ...interface{}) {
	if log.debug {
		log.output("debug", fmt.Sprintf(str, arg...))
	}
}
func (log *Log) Debug(arg ...interface{}) {
	if log.debug {
		log.output("debug", arg...)
	}
}
func (log *Log) Cost(arg ...interface{}) func() {
	begin := time.Now()
	return func() { log.output("cost", fmt.Sprint(arg...), "cost: ", FmtDuration(Now().Sub(begin))) }
}
func (log *Log) Show(arg ...interface{}) {
	level := fmt.Sprint(arg[0])
	list := []interface{}{}
	for i := 1; i < len(arg); i += 2 {
		if len(list) > 0 {
			list = append(list, " ")
		}
		list = append(list, arg[i], ": ", arg[i+1])
	}
	log.output(level, list...)
}

func New(conf *conf.Conf) *Log {
	return &Log{
		Out: os.Stdout, Err: os.Stderr,
		Conf: conf,
	}
}

var log *Log

func Init(conf *conf.Conf) {
	if log == nil {
		log = New(conf)
	}
}
