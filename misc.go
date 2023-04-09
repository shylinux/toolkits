package kit

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"
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
		res["port"] = Select(Select("80", "443", u.Scheme == "https"), strings.Split(u.Host, DF), 1)
		res["hostname"] = strings.Split(u.Host, DF)[0]
		res["origin"] = u.Scheme + "://" + u.Host
	}
	return res
}
func MergeURL(str string, arg ...Any) string {
	list := strings.Split(strings.Split(str, "#")[0], "?")
	args := map[string][]string{}
	if len(list) > 1 && list[1] != "" {
		for _, l := range strings.Split(list[1], "&") {
			ls := strings.SplitN(l, "=", 2)
			ls[0], ls[1] = QueryUnescape(ls[0]), QueryUnescape(ls[1])
			args[ls[0]] = append(args[ls[0]], ls[1])
		}
	}
	For(Simple(arg...), func(k string, v string) {
		if k == "" {
			args = map[string][]string{}
		} else if v == "" {
			delete(args, k)
		} else {
			args[k] = []string{v}
		}
	})
	res := []string{}
	For(args, func(k string, v []string) {
		For(v, func(v string) { res = append(res, url.QueryEscape(k)+"="+url.QueryEscape(v)) })
	})
	If(len(res) > 0, func() { list[0] += "?" + strings.Join(res, "&") })
	return list[0]
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
	If(!strings.HasPrefix(p, PS), func() { p = path.Join(raw.Path, get.Path) })
	return MergeURL(Select(raw.Scheme, get.Scheme)+"://"+Select(raw.Host, get.Host)+p+"?"+Select(raw.RawQuery, get.RawQuery), arg...)
}
func QueryUnescape(value string) string {
	value, _ = url.QueryUnescape(value)
	return value
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
func SplitKV(inner, outer string, text string, cb func(text string, ls []string)) {
	for _, l := range strings.Split(strings.TrimSpace(text), outer) {
		if ls := Split(strings.TrimSpace(l), inner, inner); len(ls) > 1 {
			cb(l, ls)
		}
	}
}
func SplitLine(text string) []string {
	return strings.Split(strings.TrimSpace(text), "\n")
}
func SplitWord(text string) []string {
	return Split(text, "\t ", "\t ")
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
func FormatKV(data Map, args ...string) string {
	list := []string{}
	for k, v := range data {
		list = append(list, Format("%v%v%v", k, Select(DF, args, 0), v))
	}
	return strings.Join(list, Select(";", args, 1))
}
func FormatShow(arg ...Any) string {
	res := []string{}
	for i := 0; i < len(arg); i += 2 {
		if i+1 < len(arg) {
			res = append(res, Format(arg[i])+DF, Format(arg[i+1]))
		} else {
			res = append(res, Format(arg[i]))
		}

	}
	return Join(res, " ")
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
func GetValid(cb ...func() string) (res string) {
	for _, cb := range cb {
		if res = cb(); res != "" {
			return res
		}
	}
	return ""
}
func DayBegin(t string) time.Time {
	return time.Unix(Time(Split(t)[0])/int64(time.Second), 0)
}
func IfNoKey(list Map, p string, cb func(string)) {
	If(list[p] == nil, func() {
		list[p] = true
		cb(p)
	})
}
func DictList(arg ...string) Map {
	res := Map{}
	for _, k := range arg {
		res[k] = "true"
	}
	return res
}
func Reflect(obj Any, cb func(name string, value Any)) (reflect.Type, reflect.Value) {
	t, v := reflect.TypeOf(obj), reflect.ValueOf(obj)
	If(t.Kind() == reflect.Ptr, func() { t, v = t.Elem(), v.Elem() })
	for i := 0; i < v.NumMethod(); i++ {
		method := v.Method(i)
		cb(strings.ToLower(t.Method(i).Name), method.Interface())
	}
	return t, v
}
func BeginEnd(begin, end func()) func() {
	begin()
	return end
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
