package kit

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"reflect"
	"runtime"
	"sort"
	"strings"
)

func ParseURL(str string) *url.URL {
	u, _ := url.Parse(str)
	return u
}
func MergeURL(str string, arg ...interface{}) string {
	list := strings.Split(str, "?")
	res := list[0]

	args := map[string][]string{}
	if len(list) > 1 && list[1] != "" {
		for _, l := range strings.Split(list[1], "&") {
			ls := strings.SplitN(l, "=", 2)
			args[ls[0]] = append(args[ls[0]], ls[1])
		}
	}

	list = Simple(arg...)
	for i := 0; i < len(list)-1; i += 2 {
		// args[list[i]] = append(args[list[i]], list[i+1])
		args[list[i]] = []string{list[i+1]}
	}

	list = []string{}
	for k, v := range args {
		for _, v := range v {
			// list = append(list, url.QueryEscape(k)+"="+url.QueryEscape(v))
			list = append(list, k+"="+v)
		}
	}
	if len(list) > 0 {
		res += "?" + strings.Join(list, "&")
	}
	return res
}
func MergeURL2(str string, uri string, arg ...interface{}) string {
	raw, _ := url.Parse(str)
	get, _ := url.Parse(uri)
	return MergeURL(Select(raw.Scheme, get.Scheme)+"://"+Select(raw.Host, get.Host)+""+Select(raw.Path, get.Path)+"?"+Select(raw.RawQuery, get.RawQuery), arg...)
}

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
		ext = []string{".zip", ".tar.xz", ".tar.gz", ".tar.bz2"}
	}
	str = path.Base(str)
	for _, k := range ext {
		str = strings.TrimSuffix(str, k)
	}
	return str
}
func Path(str string, rest ...string) string {
	if strings.HasPrefix(str, "/") {
		return path.Join(append([]string{str}, rest...)...)
	}
	if wd, e := os.Getwd(); e == nil {
		return path.Join(append([]string{wd, str}, rest...)...)
	}
	return str
}

func CSV(file string, limit int, cb func(index int, value map[string]string, head []string)) error {
	f, e := os.Open(file)
	if e != nil {
		return e
	}
	defer f.Close()

	r := csv.NewReader(f)
	head, e := r.Read()
	if e != nil {
		return e
	}

	for i := 0; i < limit; i++ {
		line, e := r.Read()
		if e != nil {
			break
		}

		value := map[string]string{}
		for i, k := range head {
			value[k] = line[i]
		}
		cb(i, value, head)
	}
	return nil
}
func UnMarshal(data string) interface{} {
	var res interface{}
	if strings.HasSuffix(data, ".json") {
		if b, e := ioutil.ReadFile(data); e == nil {
			json.Unmarshal(b, &res)
		}
	} else {
		json.Unmarshal([]byte(data), &res)
	}
	return res
}
func ShortKey(list map[string]interface{}, min int, arg ...interface{}) string {
	h := Hashs(arg...)
	for i := min; i < len(h); i++ {
		if _, ok := list[h[:i]]; !ok {
			return h[:i]
		}
	}
	return h
}
func KeyValue(res map[string]interface{}, key string, arg interface{}) map[string]interface{} {
	if res == nil {
		res = map[string]interface{}{}
	}
	switch arg := arg.(type) {
	case map[string]interface{}:
		for k, v := range arg {
			KeyValue(res, Keys(key, k), v)
		}

	case []interface{}:
		for i, v := range arg {
			KeyValue(res, Keys(key, i), v)
		}
	default:
		res[key] = arg
	}
	return res
}

func FileLine(p interface{}, n int) string {
	if p == nil {
		return ""
	}
	var fun uintptr
	switch p := p.(type) {
	case uintptr:
		fun = p
	case int:
		fun, _, _, _ = runtime.Caller(p)
	case nil:
		return ""
	default:
		f := reflect.ValueOf(p)
		fun = f.Pointer()
	}
	if fun == 0 {
		return ""
	}

	f := runtime.FuncForPC(fun)
	file, line := f.FileLine(fun)
	ls := strings.Split(file, "/")
	if len(ls) > n {
		ls = ls[len(ls)-n:]
	}
	return Format("%s:%d", strings.Join(ls, "/"), line)
}

func Sort(list []string) []string {
	sort.Strings(list)
	return list
}
