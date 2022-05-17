package log

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
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
func FileLine(h interface{}, n int) string {
	if h == nil || n == 0 {
		return ""
	}

	var line int
	var file string
	switch h := h.(type) {
	case int:
		_, file, line, _ = runtime.Caller(h)
	default:
		p := reflect.ValueOf(h)
		f := runtime.FuncForPC(p.Pointer())
		file, line = f.FileLine(p.Pointer())
	}

	ls := strings.Split(file, "/")
	if len(ls) > n {
		ls = ls[len(ls)-n:]
	}
	return fmt.Sprintf("%s:%d", strings.Join(ls, "/"), line)
}
