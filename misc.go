package kit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
)

func Create(p string) (*os.File, string, error) {
	if dir, _ := path.Split(p); dir != "" {
		os.MkdirAll(dir, 0777)
	}
	f, e := os.Create(p)
	return f, p, e
}
func FmtSize(size int64) string {
	if size > 1<<30 {
		return fmt.Sprintf("%d.%dG", size>>30, (size>>20)%1024*100>>10)
	}

	if size > 1<<20 {
		return fmt.Sprintf("%d.%dM", size>>20, (size>>10)%1024*100>>10)
	}

	if size > 1<<10 {
		return fmt.Sprintf("%d.%dK", size>>10, size%1024*100>>10)
	}

	return fmt.Sprintf("%dB", size)
}
func FmtTime(t int64) string {
	sign, time := "", t
	if time < 0 {
		sign, time = "-", -t
	}

	list := []string{}
	if time > 24*3600*1000000000 {
		list = append(list, fmt.Sprintf("%s%dd", sign, time/(24*3600*1000000000)))
		time = time % (24 * 3600 * 1000000000)
	}
	if time > 3600*1000000000 {
		list = append(list, fmt.Sprintf("%s%dh", sign, time/3600/1000000000))
		time = time % (3600 * 1000000000)
	}
	if len(list) > 0 {
		return strings.Join(list, "")
	}
	if time > 1000000000 {
		return fmt.Sprintf("%s%d.%ds", sign, time/1000000000, (time/1000000)%1000*100/1000)
	}
	if time > 1000000 {
		return fmt.Sprintf("%s%d.%dms", sign, time/1000000, (time/1000)%1000*100/1000)
	}
	if time > 1000 {
		return fmt.Sprintf("%s%d.%dus", sign, time/1000, (time/1000)%1000*100/1000)
	}
	return fmt.Sprintf("%s%dns", sign, time)
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
func MergeURL(str string, arg ...interface{}) string {
	list := strings.Split(str, "?")
	res := list[0]

	args := map[string][]string{}
	if len(list) > 1 {
		for _, l := range strings.Split(list[1], "&") {
			ls := strings.SplitN(l, "=", 2)
			args[ls[0]] = append(args[ls[0]], ls[1])
		}
	}

	list = Simple(arg...)
	for i := 0; i < len(list)-1; i += 2 {
		args[list[i]] = append(args[list[i]], list[i+1])
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
	if strings.HasPrefix(uri, "http") {
		str, uri = uri, ""
	}

	if ls := strings.Split(str, "//"); len(ls) > 1 && len(uri) > 0 {
		list := strings.Split(ls[1], "/")
		if strings.HasPrefix(uri, "/") {
			str = ls[0] + "//" + list[0] + uri
		} else {
			str = ls[0] + "//" + strings.Join(list[:len(list)-1], "/") + uri
		}
	}

	list := strings.Split(str, "?")
	res := list[0]

	args := map[string][]string{}
	if len(list) > 1 {
		for _, l := range strings.Split(list[1], "&") {
			ls := strings.SplitN(l, "=", 2)
			args[ls[0]] = append(args[ls[0]], ls[1])
		}
	}

	list = Simple(arg...)
	for i := 0; i < len(list)-1; i += 2 {
		args[list[i]] = append(args[list[i]], list[i+1])
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

func Width(str string, mul int) int {
	return len([]rune(str)) + (len(str)-len([]rune(str)))/2/mul
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

func Right(str string) bool {
	return str != "" && str != "0" && str != "false" && str != "off" && str != "[]" && str != "{}"
}
func Short(arg interface{}) interface{} {
	switch arg := arg.(type) {
	case string:
		if len(arg) > 6 {
			return arg[:6]
		}
	}
	return arg
}

func KeyValue(res map[string]interface{}, key string, arg interface{}) map[string]interface{} {
	switch arg := arg.(type) {
	case map[string]interface{}:
		for k, v := range arg {
			KeyValue(res, Select(Keys(key, k), k, key == ""), v)
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
