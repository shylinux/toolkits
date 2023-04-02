package kit

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
	"time"
)

func ParseQuery(str string) url.Values {
	v, _ := url.ParseQuery(str)
	return v
}
func ParseURL(str string) *url.URL {
	u, _ := url.Parse(str)
	return u
}
func ParseURLMap(str string) map[string]string {
	res := map[string]string{}
	if u := ParseURL(str); u != nil {
		res["host"] = u.Host
		res["port"] = Select(Select("80", "443", u.Scheme == "https"), strings.Split(u.Host, ":"), 1)
		res["hostname"] = strings.Split(u.Host, ":")[0]
		res["origin"] = u.Scheme + "://" + u.Host
	}
	return res
}
func MergeURL(str string, arg ...Any) string {
	list := strings.Split(str, "#")
	list = strings.Split(list[0], "?")
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
		if list[i] == "" {
			args = map[string][]string{}
			break
		}
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
		}
	}
	if len(list) > 0 {
		res += "?" + strings.Join(list, "&")
	}
	return res
}
func MergeURL2(str string, uri string, arg ...Any) string {
	raw, err := url.Parse(str)
	if err != nil {
		return MergeURL(uri, arg...)
	}
	get, err := url.Parse(uri)
	if err != nil {
		return MergeURL(str, arg...)
	}
	p := get.Path
	if !strings.HasPrefix(p, "/") {
		p = path.Join(raw.Path, get.Path)
	}
	return MergeURL(Select(raw.Scheme, get.Scheme)+"://"+Select(raw.Host, get.Host)+p+"?"+Select(raw.RawQuery, get.RawQuery), arg...)
}
func MergePOD(url string, pod string, arg ...Any) string {
	uri := ParseURL(url)
	p := uri.Query().Get("pod")
	if strings.HasPrefix(uri.Path, "/chat/pod") {
		p = strings.Split(uri.Path, "/")[3]
	}
	return MergeURL2(url, "/chat/pod/"+Keys(p, pod), "pod", "", arg)
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
func UnMarshal(buf Any) (res Any) {
	switch buf := buf.(type) {
	case []byte:
		if json.Unmarshal(buf, &res) != nil {
			return Split(string(buf))
		}
	case string:
		if strings.HasSuffix(buf, ".json") {
			if b, e := ioutil.ReadFile(buf); e == nil {
				if json.Unmarshal(b, &res) != nil {
					return Split(string(b))
				}
			}
		} else {
			if json.Unmarshal([]byte(buf), &res) != nil {
				return Split(buf)
			}
		}
	case io.Reader:
		json.NewDecoder(buf).Decode(&res)
	case io.ReadCloser:
		defer buf.Close()
		json.NewDecoder(buf).Decode(&res)
	}
	return res
}
func TransArg(arg []string, key string, trans Any) []string {
	for i := 0; i < len(arg); i += 2 {
		if arg[i] == key {
			if val := Value(trans, arg[i+1]); val != "" {
				arg[i+1] = Format(val)
			}
		}
	}
	return arg
}
func FormatKV(data Map, args ...string) string {
	list := []string{}
	for k, v := range data {
		list = append(list, Format("%v%v%v", k, Select(":", args, 0), v))
	}
	return strings.Join(list, Select(";", args, 1))
}
func FormatShow(arg ...Any) string {
	res := []string{}
	for i := 0; i < len(arg); i += 2 {
		if i+1 < len(arg) {
			res = append(res, Format(arg[i])+":", Format(arg[i+1]))
		} else {
			res = append(res, Format(arg[i]))
		}

	}
	return Join(res, " ")
}

func Now(arg ...string) string {
	return time.Now().Format(Select("2006-01-02 15:04:05", arg, 0))
}

func Reflect(obj Any, cb func(name string, value Any)) (reflect.Type, reflect.Value) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t, v = t.Elem(), v.Elem()
	}
	for i := 0; i < v.NumMethod(); i++ {
		method := v.Method(i)
		cb(strings.ToLower(t.Method(i).Name), method.Interface())
	}
	return t, v
}

func Sort(list []string, cb ...func(a, b string) bool) []string {
	if len(cb) == 0 {
		sort.Strings(list)
		return list
	}
	for i := 0; i < len(list)-1; i++ {
		for j := i; j < len(list); j++ {
			if cb[0](list[i], list[j]) {
				list[j], list[i] = list[i], list[j]
			}
		}
	}
	return list
}
func Filter(arg []string, cb ...func(string) bool) (res []string) {
	for _, k := range arg {
		for _, cb := range cb {
			if cb(k) {
				res = append(res, k)
			}
		}
	}
	return res
}
func SortedKey(obj Any) (res []string) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			switch key := key.Interface().(type) {
			case string:
				res = append(res, key)
			}
		}
	}
	sort.Strings(res)
	return res
}
func SortedValue(obj Any) (res []string) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			res = append(res, Format(val.Interface()))
		}
	}
	sort.Strings(res)
	return res
}
func DayBegin(t string) time.Time {
	return time.Unix(Time(Split(t)[0])/int64(time.Second), 0)
}

func HasPrefix(text string, arg ...string) bool {
	for _, v := range arg {
		if strings.HasPrefix(text, v) {
			return true
		}
	}
	return false
}
func HasSuffix(text string, arg ...string) bool {
	for _, v := range arg {
		if strings.HasSuffix(text, v) {
			return true
		}
	}
	return false
}
func SplitLine(text string) []string {
	return strings.Split(strings.TrimSpace(text), "\n")
}
func SplitWord(text string) []string {
	return Split(text, "\t ", "\t ")
}
func QueryUnescape(value string) string {
	value, _ = url.QueryUnescape(value)
	return value
}

func SliceRemove(list []string, key string) ([]string, int) {
	index := -1
	for i := 0; i < len(list); i++ {
		if list[i] == key {
			for index = i; i < len(list)-1; i++ {
				list[i] = list[i+1]
			}
			list = Slice(list, 0, -1)
		}
	}
	return list, index
}
func Min(arg ...int) (res int) {
	for i, v := range arg {
		if i == 0 || v < res {
			res = v
		}
	}
	return
}
func TrimPath(p string) string {
	return strings.TrimPrefix(p, Path("")+"/")
}
func AddUniq(list []string, arg ...string) []string {
	for _, k := range arg {
		if IndexOf(list, k) == -1 {
			list = append(list, k)
		}
	}
	return list
}
func ExtChange(file, ext string) string {
	if file == "" {
		return ""
	}
	return strings.TrimSuffix(file, "."+Ext(file)) + "." + ext
}
func Filters(list []string, arg ...string) (res []string) {
	for _, v := range list {
		if IndexOf(arg, v) == -1 {
			res = append(res, v)
		}
	}
	return res
}
func DictList(arg ...string) Map {
	res := Map{}
	for _, k := range arg {
		res[k] = "true"
	}
	return res
}

func GetValid(cb ...func() string) (res string) {
	for _, cb := range cb {
		if res = cb(); res != "" {
			return res
		}
	}
	return ""
}
func ExtReg(ext ...string) string { return Format(`.*\.(%s)$`, strings.Join(ext, "|")) }
func BeginEnd(begin, end func()) func() {
	begin()
	return end
}
func IfNoKey(list Map, p string, cb func(string)) {
	If(list[p] == nil, func() {
		list[p] = true
		cb(p)
	})
}
func For(val Any, cb Any) Any { return Fetch(val, cb) }
func If(exp Any, cb ...Any) {
	cbs := func(cb Any, exp Any) {
		switch cb := cb.(type) {
		case func(string):
			cb(Format(exp))
		case func(int):
			cb(Int(exp))
		case func():
			cb()
		}
	}
	switch exp := exp.(type) {
	case string:
		if exp != "" && exp != "false" {
			cbs(cb[0], exp)
		} else if len(cb) > 1 {
			cbs(cb[1], exp)
		}
	case bool:
		if exp {
			cbs(cb[0], exp)
		} else if len(cb) > 1 {
			cbs(cb[1], exp)
		}
	case int:
		if exp != 0 {
			cbs(cb[0], exp)
		} else if len(cb) > 1 {
			cbs(cb[1], exp)
		}
	default:
		if exp != nil {
			cbs(cb[0], exp)
		} else if len(cb) > 1 {
			cbs(cb[1], exp)
		}
	}
}
func Switch(exp Any, arg ...Any) Any {
	for i := 0; i < len(arg); i += 2 {
		switch val := arg[i].(type) {
		case []string:
			if IndexOf(val, Format(exp)) == -1 {
				continue
			}
		case string:
			if i < len(arg)-1 && Format(exp) != arg[i] {
				continue
			}
		default:
			if i < len(arg)-1 && exp != arg[i] {
				continue
			}
		}
		if i == len(arg)-1 {
			i--
		}
		switch cb := arg[i+1].(type) {
		case func(Any) Any:
			return cb(arg[i])
		case func(Any):
			cb(arg[i])
		case func():
			cb()
		}
		break
	}
	return nil
}
func Default(list []string, arg ...string) []string {
	if len(list) > 0 {
		return list
	}
	return arg
}
func SplitKV(inner, outer string, text string, cb func(text string, ls []string)) {
	for _, l := range strings.Split(strings.TrimSpace(text), outer) {
		if ls := Split(strings.TrimSpace(l), inner, inner); len(ls) > 1 {
			cb(l, ls)
		}
	}
}
func IsIn(v string, arg ...string) bool { return IndexOf(arg, v) > -1 }
