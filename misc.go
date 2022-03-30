package kit

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path"
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
		res["hostname"] = strings.Split(u.Host, ":")[0]
		res["origin"] = u.Scheme + "://" + u.Host
	}
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
func MergeURL2(str string, uri string, arg ...interface{}) string {
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
func MergePOD(url string, pod string, arg ...interface{}) string {
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
func UnMarshal(data string) interface{} {
	var res interface{}
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
func TransArg(arg []string, key string, trans interface{}) []string {
	for i := 0; i < len(arg); i += 2 {
		if arg[i] == key {
			if val := Value(trans, arg[i+1]); val != "" {
				arg[i+1] = Format(val)
			}
		}
	}
	return arg
}
func FormatKV(data map[string]interface{}, args ...string) string {
	list := []string{}
	for k, v := range data {
		list = append(list, Format("%v%v%v", k, Select(":", args, 0), v))
	}
	return strings.Join(list, Select(";", args, 1))
}

func Now(arg ...string) string {
	return time.Now().Format(Select("2006-01-02 15:04:05", arg, 0))
}
