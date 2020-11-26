package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/shylinux/toolkits/conf"
)

const (
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
	DEBUG = "debug"
	COST  = "cost"
	SHOW  = "show"
	STACK = 3
)

type Log struct {
	Filter   []string
	Prefix   []string
	Notice   []string
	Out, Err *os.File

	column  int
	stack   int
	debug   bool
	color   bool
	colors  map[string]string
	channel chan string
	Output  string

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
func (log *Log) prefix(level string, stack int, arg ...interface{}) (io.Writer, string) {
	if log.filter(level) {
		return nil, ""
	}

	list := []string{}
	for _, v := range log.Prefix {
		switch v {
		case "time":
			list = append(list, FmtTime(time.Now()))
		case "pid":
			list = append(list, FmtInt(os.Getpid()))
		case "fileline":
			list = append(list, FileLine(stack, log.column))
		case "level":
			list = append(list, level)
		}
	}
	return log.notice(level), strings.Join(list, " ")
}
func (log *Log) notice(key string) *os.File {
	for _, k := range log.Notice {
		if k == key {
			return log.Err
		}
	}
	return log.Out
}

var trans = map[string]string{
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
}

var LogDisable = false

func (log *Log) output(level string, arg ...interface{}) bool {
	if LogDisable {
		return false
	}
	if log == nil {
		return false
	}
	if output, prefix := log.prefix(level, log.stack, arg...); output != nil {
		color, end := trans[log.colors[level]], ""
		if color != "" {
			end = "\033[0m"
		}
		if log.channel == nil {
			fmt.Fprint(output, prefix, " ", color, fmt.Sprint(arg...), end, "\n")
		} else {
			log.channel <- fmt.Sprint(prefix, " ", color, fmt.Sprint(arg...), end, "\n")
		}
	}
	return true
}

func (log *Log) Info(arg ...interface{})  { log.output(INFO, arg...) }
func (log *Log) Warn(arg ...interface{})  { log.output(WARN, arg...) }
func (log *Log) Error(arg ...interface{}) { log.output(ERROR, arg...) }

func (log *Log) Infof(str string, arg ...interface{})  { log.output(INFO, fmt.Sprintf(str, arg...)) }
func (log *Log) Warnf(str string, arg ...interface{})  { log.output(WARN, fmt.Sprintf(str, arg...)) }
func (log *Log) Errorf(str string, arg ...interface{}) { log.output(ERROR, fmt.Sprintf(str, arg...)) }

func (log *Log) Debugf(str string, arg ...interface{}) bool {
	return log.debug && log.output(DEBUG, fmt.Sprintf(str, arg...))
}
func (log *Log) Debug(arg ...interface{}) bool {
	return log.debug && log.output(DEBUG, arg...)
}

func (log *Log) Show(arg ...interface{}) {
	list := []interface{}{}
	for i := 1; i < len(arg)-1; i += 2 {
		if len(list) > 0 {
			list = append(list, " ")
		}
		list = append(list, arg[i], ": ", arg[i+1])
	}
	log.output(fmt.Sprint(arg[0]), list...)
}
func (log *Log) Cost(arg ...interface{}) func(...func() []interface{}) {
	begin := time.Now()
	return func(cbs ...func() []interface{}) {
		list := []interface{}{fmt.Sprint(arg...)}
		for _, cb := range cbs {
			list = append(list, cb()...)
		}
		list = append(list, "cost: ", FmtDuration(Now().Sub(begin)))
		log.output(COST, list...)
	}
}

func Open(conf *conf.Conf) (*Log, error) {
	out, err := os.Stderr, os.Stderr
	switch conf.Get("log.name", "stderr") {
	case "stdout":
	case "stderr":
	default:
		if f, e := os.OpenFile(conf.Get("log.name")+".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660); e == nil {
			out = f
		}
		if f, e := os.OpenFile(conf.Get("log.name")+".log.err", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660); e == nil {
			err = f
		}
	}

	log := &Log{
		Output:  conf.Get("log.name", "stderr"),
		Prefix:  conf.GetList("log.prefix", []string{"time"}),
		Filter:  conf.GetList("log.filter"),
		Notice:  conf.GetList("log.notice"),
		column:  conf.GetInt("log.column", 3),
		debug:   conf.GetBool("log.debug"),
		color:   conf.GetBool("log.color"),
		colors:  conf.GetDict("log.colors"),
		channel: make(chan string, conf.GetInt("log.nchan", 1024)),
		Out:     out, Err: err, Conf: conf, stack: STACK,
	}

	go func() {
		for {
			select {
			case str := <-log.channel:
				fmt.Fprint(log.Out, str)
			}
		}
	}()
	return log, nil
}
