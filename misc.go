package kit

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
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
func FormatKV(data map[string]interface{}, args ...string) string {
	list := []string{}
	for k, v := range data {
		list = append(list, Format("%v%v%v", k, Select(":", args, 0), v))
	}
	return strings.Join(list, Select(";", args, 1))
}

const (
	SSH_ROUTE  = "route"
	SSH_REPOS  = "repos"
	SSH_SOURCE = "source"
	SSH_BRANCH = "branch"
	SSH_MASTER = "master"
)
