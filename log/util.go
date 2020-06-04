package log

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func FmtDuration(d time.Duration) string {
	if d > time.Second {
		return fmt.Sprintf("%d.%03ds", d/time.Second, d/time.Millisecond%1000)
	}
	if d > time.Millisecond {
		return fmt.Sprintf("%d.%03dms", d/time.Millisecond, d/time.Microsecond%1000)
	}
	if d > time.Microsecond {
		return fmt.Sprintf("%d.%03dus", d/time.Microsecond, d%1000)
	}
	return fmt.Sprintf("%dns", d)
}
func Now() time.Time {
	return time.Now()
}
func NowStamp() int64 {
	return time.Now().UnixNano()
}
func FmtSize(i uint64) string {
	if i > 1<<30 {
		return fmt.Sprintf("%d.%03dG", i>>30, (i>>20)%(1<<10))
	}
	if i > 1<<20 {
		return fmt.Sprintf("%d.%03dM", i>>20, (i>>10)%(1<<10))
	}
	if i > 1<<10 {
		return fmt.Sprintf("%d.%03dK", i>>10, i%(1<<10))
	}
	return fmt.Sprintf("%dB", i)
}
func FileLine(h interface{}, n int) string {
	if h == nil || n == 0 {
		return ""
	}

	p := reflect.ValueOf(h)
	f := runtime.FuncForPC(p.Pointer())
	file, line := f.FileLine(p.Pointer())
	ls := strings.Split(file, "/")
	if len(ls) > n {
		ls = ls[len(ls)-n:]
	}
	return fmt.Sprintf("%s:%d", strings.Join(ls, "/"), line)

}
