package kit

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
	"time"
)

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
func UnMarshal(data string) Any {
	var res Any
	if strings.HasSuffix(data, ".json") {
		if b, e := ioutil.ReadFile(data); e == nil {
			if json.Unmarshal(b, &res) != nil {
				return string(b)
			}
		}
	} else {
		if json.Unmarshal([]byte(data), &res) != nil {
			return data
		}
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

func HasPrefix(arg []string, args ...string) bool {
	if len(arg) >= len(args) {
		for i, v := range args {
			if v != arg[i] {
				return false
			}
		}
		return true
	}
	return false
}
