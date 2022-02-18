package kit

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"math/rand"
	"path"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func _list(str string) map[rune]bool {
	space := map[rune]bool{}
	for _, c := range str {
		space[c] = true
	}
	return space
}
func Split(str string, arg ...string) (res []string) {
	space := _list(Select("\t ,\n", arg, 0)) // 空白符
	block := _list(Select("{[()]}", arg, 1)) // 分隔符
	quote := _list(Select("\"'`", arg, 2))   // 引用符
	trans := _list(Select("\\", arg, 3))     // 转义符

	list := []rune(str)
	left, void, begin := '\000', true, 0
	for i := 0; i < len(list); i++ {
		switch {
		case space[list[i]]: // 空白符
			if left == '\000' {
				if !void {
					res = append(res, string(list[begin:i]))
				}
				void, begin = true, i+1
			}
		case block[list[i]]: // 分隔符
			if left == '\000' {
				if !void {
					res = append(res, string(list[begin:i]))
				}
				res = append(res, string(list[i:i+1]))
				void, begin = true, i+1
			}
		case quote[list[i]]: // 引用符
			if left == '\000' {
				left, void, begin = list[i], false, i+1
			} else if left == list[i] {
				res = append(res, string(list[begin:i]))
				left, void, begin = '\000', true, i+1
			}
		case trans[list[i]]: // 转义符
			for i := i; i < len(list)-1; i++ {
				list[i] = list[i+1]
			}
			list = list[:len(list)-1]
			void = false
		default:
			void = false
		}
	}

	if begin < len(list) { // 末尾字符
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
				Value(last, key, node)
				last_key, key, data = key, "", &node
			default:
				if key == "" {
					key = v
				} else {
					Value(last, key, v)
					key = ""
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
		case func(key string, value interface{}):
			ls := []string{}
			for k := range val {
				ls = append(ls, k)
			}

			sort.Strings(ls)
			for _, k := range ls {
				cb(k, val[k])
			}
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
				if v, ok := val[k].(map[string]interface{}); ok {
					cb(k, v)
				}
			}
		}
	case []interface{}:
		switch cb := cbs.(type) {
		case func(index int, value interface{}):
			for i, v := range val {
				cb(i, v)
			}
		case func(index int, value string):
			for i, v := range val {
				cb(i, Format(v))
			}
		case func(index int, value map[string]interface{}):
			for i, v := range val {
				cb(i, v.(map[string]interface{}))
			}
		}
	case []string:
		switch cb := cbs.(type) {
		case func(index int, value string):
			for i, v := range val {
				cb(i, v)
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
		case "uniq":
			args = append(args, Format(time.Now()))
			args = append(args, Format(rand.Int()))
		case "time":
			args = append(args, Format(time.Now()))
		case "rand":
			args = append(args, Format(rand.Int()))
		default:
			args = append(args, v)
		}
	}

	h := md5.Sum([]byte(strings.Join(args, "")))
	return hex.EncodeToString(h[:]), args
}
func HashsPath(arg ...interface{}) string {
	h := Hashs(arg...)
	return path.Join(h[:2], h)
}
func Hashs(arg ...interface{}) string {
	if len(arg) > 0 {
		switch arg := arg[0].(type) {
		case []byte:
			md := md5.New()
			md.Write(arg)
			h := md.Sum(nil)
			return hex.EncodeToString(h[:])
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
		"Format": Format, "Value": Value,
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

func GetMeta(value map[string]interface{}) map[string]interface{} {
	if value[MDB_META] != nil {
		value = value[MDB_META].(map[string]interface{})
	}
	return value
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
func ShortKey(list map[string]interface{}, min int, arg ...interface{}) string {
	h := Hashs(arg...)
	for i := min; i < len(h); i++ {
		if _, ok := list[h[:i]]; !ok {
			return h[:i]
		}
	}
	return h
}
func Fields(arg ...interface{}) string { return Join(Simple(arg...)) }
