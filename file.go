package kit

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"reflect"
	"runtime"
	"strings"
)

const (
	DF = ":"
	PS = "/"
	PT = "."
	FS = ","
)

func UserName() string {
	if user, err := user.Current(); err == nil && user.Name != "" {
		return user.Name
	}
	return Select("root", Select(os.Getenv("LOGNAME"), os.Getenv("USER")))
}
func HomePath(str string, rest ...string) string {
	if user, err := user.Current(); err == nil {
		return Path(path.Join(user.HomeDir, str), rest...)
	}
	return Path(path.Join(os.Getenv("HOME"), str), rest...)
}
func Path(str string, rest ...string) string {
	if sep := string([]rune{os.PathSeparator}); strings.HasPrefix(str, sep) || strings.Contains(str, DF) {
		str = path.Join(append([]string{str}, rest...)...) + Select("", sep, len(rest) == 0 && strings.HasSuffix(str, sep))
	} else if wd, e := os.Getwd(); e == nil {
		str = path.Join(append([]string{wd, str}, rest...)...)
	}
	return strings.ReplaceAll(str, "\\", PS)
}
func Pwd() string {
	wd, _ := os.Getwd()
	return ReplaceAll(wd, "\\", PS)
}
func Env(key string) string { return os.Getenv(key) }
func EnvSimple(arg ...string) (res []string) {
	For(arg, func(k string) { res = append(res, k, Env(k)) })
	return
}
func IsDir(p string) bool {
	if _, e := ioutil.ReadDir(p); e == nil {
		return true
	}
	return false
}
func ExtReg(ext ...string) string { return Format(`.*\.(%s)$`, Join(ext, "|")) }
func ExtChange(file, ext string) string {
	if file == "" {
		return ""
	}
	return strings.TrimSuffix(file, PT+Ext(file)) + PT + ext
}
func Ext(str string) string {
	return strings.ToLower(path.Base(Select(str, strings.TrimPrefix(path.Ext(str), PT))))
}
func TrimPath(p string) string {
	return strings.TrimPrefix(p, Path("")+PS)
}
func TrimExt(str string, ext ...string) string {
	If(len(ext) == 0, func() { ext = []string{"go", "shy", "zip", "tgz", "tar.gz", "tar.xz", "tar.bz2"} })
	str = path.Base(str)
	for _, k := range ext {
		if k == "" {
			str = strings.Split(str, PT)[0]
		} else {
			str = strings.TrimSuffix(str, PT+k)
		}
	}
	return str
}
func ExtIsImage(str string) bool {
	switch strings.ToLower(Ext(str)) {
	case "png", "jpg", "jpeg":
		return true
	}
	return false
}
func ExtIsVideo(str string) bool {
	switch strings.ToLower(Ext(str)) {
	case "webm":
		return true
	}
	return false
}

func getFunc(p interface{}) (fun uintptr) {
	if p == nil {
		return 0
	}
	switch p := p.(type) {
	case uintptr:
		fun = p
	case int:
		fun, _, _, _ = runtime.Caller(p + 1)
	case nil:
		fun = 0
	default:
		fun = reflect.ValueOf(p).Pointer()
	}
	return fun
}
func ModName(p interface{}) string {
	ls := strings.Split(runtime.FuncForPC(getFunc(p)).Name(), PS)
	if strings.Contains(ls[0], PT) {
		return Select(ls[0], ls, 2)
	}
	return ls[0]
}
func PathName(p interface{}) string {
	fun := getFunc(p)
	if fun == 0 {
		return ""
	}
	file, _ := runtime.FuncForPC(fun).FileLine(fun)
	return path.Base(path.Dir(file))
}
func FileName(p interface{}) string {
	fun := getFunc(p)
	if fun == 0 {
		return ""
	}
	file, _ := runtime.FuncForPC(fun).FileLine(fun)
	return strings.Split(path.Base(file), PT)[0]
}
func FuncName(p interface{}) string {
	fun := getFunc(p)
	if fun == 0 {
		return ""
	}
	list := strings.Split(runtime.FuncForPC(fun).Name(), PT)
	return strings.TrimSuffix(list[len(list)-1], "-fm")
}
func FileLine(p interface{}, n int) string {
	fun := getFunc(p)
	if fun == 0 {
		return ""
	}
	file, line := runtime.FuncForPC(fun).FileLine(fun)
	list := strings.Split(file, PS)
	If(len(list) > n, func() { list = list[len(list)-n:] })
	return Format("%s:%d", strings.TrimPrefix(strings.Join(list, PS), Path("")+PS), line)
}
