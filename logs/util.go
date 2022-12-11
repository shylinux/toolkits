package logs

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	kit "shylinux.com/x/toolkits"
)

func Now() time.Time  { return time.Now() }
func NowStamp() int64 { return time.Now().UnixNano() }

func FmtTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000")
}
func FmtDuration(d time.Duration) string {
	if d > time.Second {
		return fmt.Sprintf("%.2fs", float64(d)/float64(time.Second))
	}
	if d > time.Millisecond {
		return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
	}
	if d > time.Microsecond {
		return fmt.Sprintf("%.2fus", float64(d)/float64(time.Microsecond))
	}
	return fmt.Sprintf("%dns", d)
}
func FmtInt(i int) string {
	return strconv.FormatInt(int64(i), 10)
}
func FmtSize(i uint64) string {
	if i > 100000000000 {
		return fmt.Sprintf("%.2fG", float64(i)/10000000000)
	}
	if i > 1000000 {
		return fmt.Sprintf("%.2fM", float64(i)/1000000)
	}
	if i > 1000 {
		return fmt.Sprintf("%.2fK", float64(i)/1000)
	}
	return fmt.Sprintf("%dB", i)
}
func FileLines(h interface{}) string {
	if h == nil {
		return ""
	}
	var line int
	var file string
	switch h := h.(type) {
	case nil:
		return ""
	case []string:
		return kit.Join(h, " ")
	case string:
		return h
	case int:
		if h < 0 {
			i := 1 - h
			call := path.Dir(FileLines(i))
			for i++; i < 10; i++ {
				if path.Dir(FileLines(i)) != call {
					h = i - 1
					break
				}
			}
		}
		_, file, line, _ = runtime.Caller(h)
	default:
		if t := reflect.TypeOf(h); t.Kind() != reflect.Func {
			return ""
		}
		p := reflect.ValueOf(h)
		if p.Pointer == nil || p.Pointer() == 0 {
			return ""
		}
		f := runtime.FuncForPC(p.Pointer())
		file, line = f.FileLine(p.Pointer())
	}
	return fmt.Sprintf("%s:%d", file, line)
}
func FileLine(h interface{}, arg ...string) string {
	switch n := h.(type) {
	case int:
		if n > 0 {
			h = n + 1
		} else {
			h = n - 1
		}
	}
	ls := strings.Split(FileLines(h), "/")
	if n := kit.Int(kit.Select("3", arg, 0)); len(ls) > n {
		ls = ls[len(ls)-n:]
	}
	return strings.Join(ls, "/")
}
func CostTime(cb func(time.Duration)) func() {
	begin := Now()
	return func() { cb(Now().Sub(begin)) }
}
func Println(arg ...Any) {
	if len(arg) == 0 {
		println()
	} else {
		println(FmtTime(Now()), kit.Format(arg[0], arg[1:]...), FileLine(2))
	}
}
func PrintStack() {
	Println(Stack(2, 100))
}
func Stack(skip, deep int) string {
	pc := make([]uintptr, deep+10)
	frames := runtime.CallersFrames(pc[:runtime.Callers(skip+1, pc)])

	list := []string{}
	for {
		frame, more := frames.Next()
		file := kit.Slice(kit.Split(frame.File, "/", "/"), -1)[0]
		name := kit.Slice(kit.Split(frame.Function, "/", "/"), -1)[0]
		list = append(list, kit.Format("%s:%d\t%s", file, frame.Line, name))

		if len(list) >= deep {
			break
		}
		if !more {
			break
		}
	}
	return kit.Join(list, "\n")
}
