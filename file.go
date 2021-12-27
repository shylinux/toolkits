package kit

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
)

func Create(p string) (*os.File, string, error) {
	switch p {
	case "", "null":
		p = "/dev/null"
	case "stdout", "stderr":
		p = "/dev/" + p
	}

	if dir, _ := path.Split(p); dir != "" {
		os.MkdirAll(dir, 0777)
	}
	f, e := os.Create(p)
	return f, p, e
}
func Rewrite(file string, cb func(string) string) error {
	f, e := os.Open(file)
	if e != nil {
		return e
	}
	defer f.Close()

	b, e := ioutil.ReadAll(f)
	if e != nil {
		return e
	}
	bio := bufio.NewScanner(bytes.NewBuffer(b))

	o, _, e := Create(file)
	if e != nil {
		return e
	}
	defer o.Close()

	for bio.Scan() {
		line := cb(bio.Text())
		o.WriteString(line)
		o.WriteString("\n")
	}
	return nil
}
func FileExists(name string) bool {
	if s, e := os.Stat(name); s != nil && e == nil {
		return true
	}
	return false
}
func TrimExt(str string, ext ...string) string {
	if len(ext) == 0 {
		ext = []string{".go", ".zip", ".tgz", ".tar.gz", ".tar.xz", ".tar.bz2"}
	}
	str = path.Base(str)
	for _, k := range ext {
		str = strings.TrimSuffix(str, k)
	}
	return str
}
func Path(str string, rest ...string) string {
	if strings.HasPrefix(str, "/") {
		return path.Join(append([]string{str}, rest...)...) + Select("", "/", len(rest) == 0 && strings.HasSuffix(str, "/"))
	}
	if wd, e := os.Getwd(); e == nil {
		return path.Join(append([]string{wd, str}, rest...)...)
	}
	return str
}
func Ext(str string) string {
	return strings.ToLower(path.Base(Select(str, strings.TrimPrefix(path.Ext(str), "."))))
}
func ExtIsImage(str string) bool {
	switch strings.ToLower(Ext(str)) {
	case "png", "jpg", "jpeg":
		return true
	}
	return false
}
func Pwd() string {
	wd, _ := os.Getwd()
	return wd
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
	ls := strings.Split(runtime.FuncForPC(getFunc(p)).Name(), "/")
	if strings.Contains(ls[0], ".") {
		return Select(ls[0], ls, 2)
	}
	return ls[0]
}
func ModPath(p interface{}, arg ...string) string {
	ls := strings.Split(runtime.FuncForPC(getFunc(p)).Name(), "/")
	ls[len(ls)-1] = strings.Split(ls[len(ls)-1], ".")[0]
	return path.Join(path.Join(ls...), path.Join(arg...))
}
func FuncName(p interface{}) string {
	fun := getFunc(p)
	if fun == 0 {
		return ""
	}
	list := strings.Split(runtime.FuncForPC(fun).Name(), ".")
	return strings.TrimSuffix(list[len(list)-1], "-fm")
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
	return strings.Split(path.Base(file), ".")[0]
}
func FileLine(p interface{}, n int) string {
	fun := getFunc(p)
	if fun == 0 {
		return ""
	}

	file, line := runtime.FuncForPC(fun).FileLine(fun)
	list := strings.Split(file, "/")
	if len(list) > n {
		list = list[len(list)-n:]
	}
	return Format("%s:%d", strings.Join(list, "/"), line)
}
