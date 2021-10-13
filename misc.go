package kit

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"unicode"
)

func ParseURL(str string) *url.URL {
	u, _ := url.Parse(str)
	return u
}
func ParseURLMap(str string) map[string]string {
	u := ParseURL(str)
	res := map[string]string{}
	res["host"] = u.Host
	return res
}
func MergeURL(str string, arg ...interface{}) string {
	list := strings.Split(str, "?")
	res := list[0]

	args := map[string][]string{}
	if len(list) > 1 && list[1] != "" {
		for _, l := range strings.Split(list[1], "&") {
			ls := strings.SplitN(l, "=", 2)
			ls[0], _ = url.QueryUnescape(ls[0])
			ls[1], _ = url.QueryUnescape(ls[1])
			args[ls[0]] = append(args[ls[0]], ls[1])
		}
	}

	list = Simple(arg...)
	for i := 0; i < len(list)-1; i += 2 {
		// args[list[i]] = append(args[list[i]], list[i+1])
		if list[i+1] == "" {
			delete(args, list[i])
			continue
		}
		args[list[i]] = []string{list[i+1]}
	}

	list = []string{}
	for k, v := range args {
		for _, v := range v {
			list = append(list, url.QueryEscape(k)+"="+url.QueryEscape(v))
			// list = append(list, k+"="+v)
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
	p := get.Path
	if !strings.HasPrefix(p, "/") {
		p = path.Join(raw.Path, get.Path)
	}
	return MergeURL(Select(raw.Scheme, get.Scheme)+"://"+Select(raw.Host, get.Host)+path.Join("/", p)+"?"+Select(raw.RawQuery, get.RawQuery), arg...)
}
func MergePOD(str string, pod string) string {
	return MergeURL(str, "pod", Keys(ParseURL(str).Query().Get("pod"), pod))
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
		ext = []string{".zip", ".tgz", ".tar.xz", ".tar.gz", ".tar.bz2"}
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

func SourcePath(arg ...string) string {
	pp := path.Join(arg...)
	if strings.HasPrefix(pp, "/") {
		return pp
	}

	ls := strings.Split(FileLine(2, 100), "usr")
	return path.Join("/require/shylinux.com/x", path.Dir(ls[len(ls)-1]), pp)
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
	return Split(reflect.TypeOf(p).PkgPath(), "/-")[2]
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

func TransArg(arg []string, key string, trans map[string]string) []string {
	for i := 0; i < len(arg); i += 2 {
		if arg[i] == key {
			if val, ok := trans[arg[i+1]]; ok {
				arg[i+1] = val
			}
		}
	}
	return arg
}
func Capital(str string) string {
	return string(unicode.ToUpper(rune(str[0]))) + str[1:]
}
func SubKey(name string) string {
	return Keys(MDB_HASH, Hashs(name))
}
func FormatKV(data map[string]interface{}, args ...string) string {
	list := []string{}
	for k, v := range data {
		list = append(list, Format("%v%v%v", k, Select(":", args, 0), v))
	}
	return strings.Join(list, Select(";", args, 1))
}

func Sort(list []string) []string {
	sort.Strings(list)
	return list
}
func ForEach(arg []string, cb func(string)) {
	for _, v := range arg {
		cb(v)
	}
}
func Contains(str, sub interface{}) bool {
	return strings.Contains(Format(str), Format(sub))
}
func Regexp(arg string) *regexp.Regexp {
	reg, _ := regexp.Compile(arg)
	return reg
}

func Replace(str string, from string, to string) string {
	trans := map[rune]rune{}
	for i, c := range []rune(from) {
		switch to := []rune(to); len(to) {
		case 0:
			trans[c] = '\000'
		case 1:
			trans[c] = to[0]
		default:
			if i < len(to) {
				trans[c] = to[i]
			} else {
				trans[c] = '\000'
			}
		}
	}

	res := []rune{}
	for _, c := range str {
		switch c := trans[c]; c {
		case '\000':
			continue
		default:
			res = append(res, trans[c])
		}
	}
	return string(res)
}
func Join(str []string, arg ...string) string {
	return strings.Join(str, Select(",", arg, 0))
}

type ReadCloser struct {
	r io.Reader
}

func (r *ReadCloser) Read(buf []byte) (int, error) {
	return r.r.Read(buf)
}
func (r *ReadCloser) Close() error {
	if c, ok := r.r.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
func NewReadCloser(r io.Reader) *ReadCloser {
	return &ReadCloser{r: r}
}
