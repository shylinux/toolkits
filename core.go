package kit

import (
	"strconv"
	"strings"
)

// 字符串切分
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

// 字符串解析
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

// 数据读写
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

				if e == nil {
					data, index = []interface{}{next}, 0
				} else {
					data = map[string]interface{}{key: next}
				}
			case []string:
				index = (index+2+len(value)+2)%(len(value)+2) - 2

				if j == len(keys)-1 {
					if i == len(args)-1 {
						if index < 0 {
							return ""
						}
						return value[index]
					}
					next = args[i+1]
				}

				if index == -1 {
					data, index = append([]string{Format(next)}, value...), 0
				} else {
					data, index = append(value, Format(next)), len(value)
				}
				next = value[index]
			case map[string]string:
				if j == len(keys)-1 {
					if i == len(args)-1 {
						return value[key] // 读取数据
					}
					value[key] = Format(next) // 修改数据
				}
				next = value[key]
			case map[string]interface{}:
				if j == len(keys)-1 {
					if i == len(args)-1 {
						return value[key] // 读取数据
					}
					value[key] = args[i+1] // 修改数据
					if s, ok := args[i+1].(string); ok && s == "" {
						delete(value, key)
					}
				}
				next = value[key]
			case []interface{}:
				index = (index+2+len(value)+2)%(len(value)+2) - 2

				if j == len(keys)-1 {
					if i == len(args)-1 {
						if index < 0 {
							return nil
						}
						return value[index] // 读取数据
					}
					next = args[i+1] // 修改数据
				}

				if index == -1 {
					value, index = append([]interface{}{next}, value...), 0
				} else if index == -2 {
					value, index = append(value, next), len(value)
				} else if j == len(keys)-1 {
					value[index] = next
				}
				data, next = value, value[index]
			}

			switch p := parent.(type) {
			case map[string]interface{}:
				p[parent_key] = data
			case []interface{}:
				p[parent_index] = data
			case nil:
				root = data
			}

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
			for k, v := range val {
				cb(k, Format(v))
			}
		case func(key string, value map[string]interface{}):
			for k, v := range val {
				cb(k, v.(map[string]interface{}))
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

func Table(data interface{}, offset, limit int, cb interface{}) interface{} {
	switch data := data.(type) {
	case map[string]interface{}:
		switch cb := cb.(type) {
		case func(key string, value string):
			for k, v := range data {
				cb(k, Format(v))
			}
		}
	case []interface{}:
		switch cb := cb.(type) {
		case func(key int, value string):
			for i, v := range data {
				cb(i, Format(v))
			}
		}
	}
	return data
}
