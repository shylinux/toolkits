package kit

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func Split(str string, arg ...string) (res []string) {
	sep := []rune("\t \n")
	if len(arg) > 0 && len(arg[0]) > 0 {
		sep = []rune(arg[0])
	}
	for i := len(sep); i < 5; i++ {
		sep = append(sep, sep[0])
	}

	trip := false
	list := []rune(str)
	left, space, begin := '\000', true, 0
	for i := 0; i < len(list); i++ {
		switch list[i] {
		case '\'', '"', '`':
			if left == '\000' {
				if space && (len(arg) == 0 || arg[0] != "\n") {
					begin = i + 1
				}
				left, space = list[i], false
			} else if left == list[i] {
				left, space = '\000', false
				trip = true
			}
		case sep[0], sep[1], sep[2], sep[3], sep[4]:
			if left != '\000' {
				break
			}
			if !space {
				if trip {
					res = append(res, string(list[begin:i-1]))
				} else if i > 0 {
					res = append(res, string(list[begin:i]))
				}
			}
			space, begin = true, i+1
		case '\\':
			space = false
		default:
			trip = false
			space = false
		}
	}

	if begin < len(list) {
		res = append(res, string(list[begin:]))
	}
	return res
}
func Parse(value interface{}, key string, val ...string) interface{} {
	list := []*interface{}{&value}
	data := &value

	last_key := ""
	for _, v := range val {
		var node interface{}
		switch v {
		case ":", ",":
			continue
		case "]", "}":
			if len(list) == 1 {
				return *list[0]
			}
			data, list = list[len(list)-2], list[:len(list)-1]
			continue
		case "{":
			node = map[string]interface{}{}
		case "[":
			node = []interface{}{}
		default:
			node = v
		}

		switch last := (*data).(type) {
		case map[string]interface{}:
			switch v {
			case "{", "[":
				list = append(list, &node)
				last[key], last_key, key, data = node, key, "", &node
			default:
				if key == "" {
					key = v
				} else {
					last[key], key = v, ""
				}
			}

		case []interface{}:
			last = append(last, node)
			*data = last
			list[len(list)-1] = data

			if len(list) > 1 {
				switch p := (*list[len(list)-2]).(type) {
				case map[string]interface{}:
					p[last_key] = last
				case []interface{}:
					p[len(p)-1] = last
				}
			}

			switch v {
			case "{", "[":
				list = append(list, &node)
				data = &node
			}

		case nil:
			switch v {
			case "{", "[":
				data, list[0] = &node, &node
			default:
				return v
			}
		}
	}
	return *list[0]
}
func Value(root interface{}, args ...interface{}) interface{} {
	for i := 0; i < len(args); i += 2 {
		if arg, ok := args[i].(map[string]interface{}); ok {
			argn := []interface{}{}
			for k, v := range arg {
				argn = append(argn, k, v)
			}
			argn = append(argn, args[i+1:])
			args, i = argn, -2
			continue
		}

		// 解析索引
		keys := []string{}
		for _, v := range Simple(args[i]) {
			keys = append(keys, strings.Split(v, ".")...)
		}

		var parent interface{}
		parent_key, parent_index := "", 0

		data := root
		for j, key := range keys {
			index, e := strconv.Atoi(key)

			var next interface{}
			switch value := data.(type) {
			case nil:
				if i == len(args)-1 {
					return nil
				}
				if j == len(keys)-1 {
					next = args[i+1]
				}

				// 创建数据
				if e == nil {
					data, index = []interface{}{next}, 0
				} else {
					data = map[string]interface{}{key: next}
				}
			case []string:
				index = (index+2+len(value)+2)%(len(value)+2) - 2

				if j == len(keys)-1 {
					if i == len(args)-1 {
						// 读取数据
						if index < 0 {
							return ""
						}
						return value[index]
					}
					next = args[i+1]
				}

				// 添加数据
				if index == -1 {
					data, index = append([]string{Format(next)}, value...), 0
				} else {
					data, index = append(value, Format(next)), len(value)
				}
				next = value[index]
			case map[string]string:
				if j == len(keys)-1 {
					if i == len(args)-1 {
						// 读取数据
						return value[key]
					}
					// 修改数据
					value[key] = Format(next)
				}
				next = value[key]
			case map[string]interface{}:
				if j == len(keys)-1 {
					if i == len(args)-1 {
						// 读取数据
						if key == "" {
							return root
						}
						return value[key]
					}
					// 修改数据
					value[key] = args[i+1]
					if s, ok := args[i+1].(string); ok && s == "" {
						// 删除数据
						delete(value, key)
					}
				}
				next = value[key]
			case []interface{}:
				index = (index+2+len(value)+2)%(len(value)+2) - 2

				if j == len(keys)-1 {
					if i == len(args)-1 {
						// 读取数据
						if index < 0 {
							return nil
						}
						return value[index]
					}
					next = args[i+1]
				}

				// 添加数据
				if index == -1 {
					value, index = append([]interface{}{next}, value...), 0
				} else if index == -2 {
					value, index = append(value, next), len(value)
				} else if j == len(keys)-1 {
					value[index] = next
				}
				data, next = value, value[index]
			}

			// 添加索引
			switch p := parent.(type) {
			case map[string]interface{}:
				p[parent_key] = data
			case []interface{}:
				p[parent_index] = data
			case nil:
				root = data
			}

			// 索引递进
			parent, data = data, next
			parent_key, parent_index = key, index
		}
	}

	return root
}
func Fetch(val interface{}, cbs interface{}) interface{} {
	switch val := val.(type) {
	case map[string]interface{}:
		switch cb := cbs.(type) {
		case func(value map[string]interface{}):
			cb(val)
		case func(key string, value string):
			ls := []string{}
			for k := range val {
				ls = append(ls, k)
			}

			sort.Strings(ls)
			for _, k := range ls {
				cb(k, Format(val[k]))
			}
		case func(key string, value map[string]interface{}):
			ls := []string{}
			for k := range val {
				ls = append(ls, k)
			}
			sort.Strings(ls)
			for _, k := range ls {
				cb(k, val[k].(map[string]interface{}))
			}
		}
	case []interface{}:
		switch cb := cbs.(type) {
		case func(index int, value string):
			for i, v := range val {
				cb(i, Format(v))
			}
		case func(index int, value map[string]interface{}):
			for i, v := range val {
				cb(i, v.(map[string]interface{}))
			}
		}
	}
	return val
}

func Hash(arg ...interface{}) (string, []string) {
	if len(arg) == 0 {
		arg = append(arg, "uniq")
	}
	args := []string{}
	for _, v := range Simple(arg...) {
		switch v {
		case "time":
			args = append(args, Format(time.Now()))
		case "rand":
			args = append(args, Format(rand.Int()))
		case "uniq":
			args = append(args, Format(time.Now()))
			args = append(args, Format(rand.Int()))
		default:
			args = append(args, v)
		}
	}

	h := md5.Sum([]byte(strings.Join(args, "")))
	return hex.EncodeToString(h[:]), args
}
func Hashs(arg ...interface{}) string {
	if len(arg) > 0 {
		switch arg := arg[0].(type) {
		case string:
		case io.Reader:
			md := md5.New()
			io.Copy(md, arg)
			h := md.Sum(nil)
			return hex.EncodeToString(h[:])
		}
	}
	h, _ := Hash(arg...)
	return h
}
func Render(str string, arg interface{}) (b []byte, e error) {
	t := template.New("render").Funcs(template.FuncMap{
		"Value": Value,
	})
	if strings.HasPrefix(str, "@") {
		if t, e = template.ParseFiles(str[1:]); e != nil {
			return nil, e
		}
	} else {
		if t, e = t.Parse(str); e != nil {
			return nil, e
		}
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	if e := t.Execute(buf, arg); e != nil {
		return nil, e
	}
	return buf.Bytes(), nil
}
