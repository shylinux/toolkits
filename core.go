package kit

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode"
)

func _list(str string) map[rune]bool {
	space := map[rune]bool{}
	for _, c := range str {
		space[c] = true
	}
	return space
}
func Split(str string, arg ...string) (res []string) {
	space := _list(Select("\t ,\r\n", arg, 0)) // 空白符
	block := _list(Select("{[()]}", arg, 1))   // 分隔符
	quote := _list(Select("\"'`", arg, 2))     // 引用符
	trans := _list(Select("\\", arg, 3))       // 转义符
	raw := Select("", arg, 4) == "true"        // 转义符
	list := []rune(str)
	left, void, begin := '\000', true, 0
	for i := 0; i < len(list); i++ {
		switch {
		case space[list[i]]: // 空白符
			if left == '\000' {
				If(!void, func() { res = append(res, string(list[begin:i])) })
				void, begin = true, i+1
			}
		case block[list[i]]: // 分隔符
			if left == '\000' {
				If(!void, func() { res = append(res, string(list[begin:i])) })
				res = append(res, string(list[i:i+1]))
				void, begin = true, i+1
			}
		case quote[list[i]]: // 引用符
			if left == '\000' {
				left, void, begin = list[i], false, i+1
			} else if left == list[i] {
				if raw {
					res = append(res, string(list[begin-1:i+1]))
				} else {
					res = append(res, string(list[begin:i]))
				}
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
	If(begin < len(list), func() { res = append(res, string(list[begin:])) })
	return res
}
func Parse(value Any, key string, val ...string) Any {
	list := []*Any{&value}
	data := &value
	last_key := ""
	for _, v := range val {
		var node Any
		switch v {
		case DF, FS:
			continue
		case "]", "}":
			if len(list) == 1 {
				return *list[0]
			}
			data, list = list[len(list)-2], list[:len(list)-1]
			continue
		case "{":
			node = Map{}
		case "[":
			node = []Any{}
		default:
			node = v
		}
		switch last := (*data).(type) {
		case Map:
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
		case []Any:
			last = append(last, node)
			*data = last
			list[len(list)-1] = data
			if len(list) > 1 {
				switch p := (*list[len(list)-2]).(type) {
				case Map:
					p[last_key] = last
				case []Any:
					p[len(p)-1] = last
				}
			}
			switch v {
			case "{", "[":
				data, list = &node, append(list, &node)
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
func Value(root Any, args ...Any) Any {
	switch val := root.(type) {
	case string:
		root = UnMarshal(val)
	}
	for i := 0; i < len(args); i += 2 {
		if arg, ok := args[i].(Map); ok {
			argn := []Any{}
			for k, v := range arg {
				argn = append(argn, k, v)
			}
			argn = append(argn, args[i+1:])
			args, i = argn, -2
			continue
		}
		keys := []string{}
		for _, v := range Simple(args[i]) {
			keys = append(keys, strings.Split(v, PT)...)
		}
		var parent Any
		parent_key, parent_index := "", 0
		data := root
		for j, key := range keys {
			index, e := strconv.Atoi(key)
			var next Any
			switch value := data.(type) {
			case nil:
				if i == len(args)-1 {
					return nil
				}
				If(j == len(keys)-1, func() { next = args[i+1] })
				if e == nil {
					data, index = []Any{next}, 0
				} else {
					data = Map{key: next}
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
					value, index = append([]string{Format(next)}, value...), 0
				} else if index == -2 {
					value, index = append(value, Format(next)), len(value)
				} else if j == len(keys)-1 {
					value[index] = Format(next)
				}
				data, next = value, value[index]
			case map[string]string:
				if j == len(keys)-1 {
					if i == len(args)-1 {
						return value[key]
					}
					value[key] = Format(next)
				}
				next = value[key]
			case Map:
				if j == len(keys)-1 {
					if i == len(args)-1 {
						if key == "" {
							return root
						}
						return value[key]
					}
					if args[i+1] == nil {
						delete(value, key)
					} else if s, ok := args[i+1].(string); ok && s == "" {
						delete(value, key)
					} else {
						value[key] = args[i+1]
					}
				}
				next = value[key]
			case []Any:
				index = (index+2+len(value)+2)%(len(value)+2) - 2
				if j == len(keys)-1 {
					if i == len(args)-1 {
						if index < 0 {
							return nil
						}
						return value[index]
					}
					next = args[i+1]
				}
				if index == -1 {
					value, index = append([]Any{next}, value...), 0
				} else if index == -2 {
					value, index = append(value, next), len(value)
				} else if j == len(keys)-1 {
					value[index] = next
				}
				data, next = value, value[index]
			}
			switch p := parent.(type) {
			case Map:
				p[parent_key] = data
			case []Any:
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
func Fetch(val Any, cb Any) Any {
	switch val := val.(type) {
	case map[string]Any:
		for _, k := range SortedKey(val) {
			switch cb := cb.(type) {
			case func(k string):
				cb(k)
			case func(k, v string):
				cb(k, Format(val[k]))
			case func(k string, v Any):
				cb(k, val[k])
			case func(k string, v Map):
				if v, ok := val[k].(Map); ok {
					cb(k, v)
				}
			}
		}
	case map[string][]string:
		for _, k := range SortedKey(val) {
			switch cb := cb.(type) {
			case func(k string, v []string):
				cb(k, val[k])
			}
		}
	case map[string]string:
		for _, k := range SortedKey(val) {
			switch cb := cb.(type) {
			case func(k, v string):
				cb(k, val[k])
			}
		}
	case map[string]int:
		for _, k := range SortedKey(val) {
			switch cb := cb.(type) {
			case func(k string, v int):
				cb(k, val[k])
			}
		}
	case []string:
		switch cb := cb.(type) {
		case func(v string):
			for _, v := range val {
				cb(v)
			}
		case func(k, v string):
			for i := 0; i < len(val)-1; i += 2 {
				cb(val[i], val[i+1])
			}
		case func(i int, v string):
			for i, v := range val {
				cb(i, v)
			}
		}
	case []Any:
		switch cb := cb.(type) {
		case func(v Any):
			for _, v := range val {
				cb(v)
			}
		case func(k, v Any):
			for i := 0; i < len(val)-1; i += 2 {
				cb(val[i], val[i+1])
			}
		case func(i int, v Any):
			for i, v := range val {
				cb(i, v)
			}
		case func(i int, v Map):
			for i, v := range val {
				cb(i, v.(Map))
			}
		case func(v Map):
			for _, v := range val {
				cb(v.(Map))
			}
		case func(v string):
			for _, v := range val {
				cb(Format(v))
			}
		case func(i int, v string):
			for i, v := range val {
				cb(i, Format(v))
			}
		}
	case url.Values:
		for _, k := range SortedKey(val) {
			switch cb := cb.(type) {
			case func(k string, v []string):
				cb(k, val[k])
			}
		}
	case http.Header:
		for _, k := range SortedKey(val) {
			switch cb := cb.(type) {
			case func(k string, v []string):
				cb(k, val[k])
			}
		}
	case []*http.Cookie:
		for _, v := range val {
			switch cb := cb.(type) {
			case func(k, v string):
				cb(v.Name, v.Value)
			}
		}
	case *bufio.Scanner:
		for bio, i := val, 0; bio.Scan(); i++ {
			switch cb := cb.(type) {
			case func(s string, i int):
				cb(bio.Text(), i)
			case func(s string):
				cb(bio.Text())
			case func(ls []string, s string):
				cb(Split(bio.Text()), bio.Text())
			}
		}
	case io.Reader:
		Fetch(bufio.NewScanner(val), cb)
	case []os.FileInfo:
		for _, s := range val {
			switch cb := cb.(type) {
			case func(os.FileInfo):
				cb(s)
			}
		}
	case int:
		switch cb := cb.(type) {
		case func(int):
			for i := 0; i < val; i++ {
				cb(i)
			}
		case func():
			for i := 0; i < val; i++ {
				cb()
			}
		}
	case interface{ Operate(string, Any) Any }:
		list := []Any{}
		for {
			if list, _ = val.Operate("range", list).([]Any); list == nil {
				break
			}
			switch cb := cb.(type) {
			case func(string, Any):
				cb(Format(list[0]), list[1])
			}
		}
	case nil:
	default:
		panic(Format("not implements: %#v %v", val, FileLine(cb, 3)))
	}
	return val
}

func Hash(arg ...Any) (string, []string) {
	If(len(arg) == 0, func() { arg = append(arg, MDB_UNIQ) })
	args := []string{}
	for _, v := range Simple(arg...) {
		switch v {
		case MDB_UNIQ:
			args = append(args, Format(time.Now()))
			args = append(args, Format(rand.Int()))
		case MDB_TIME:
			args = append(args, Format(time.Now()))
		case MDB_RAND:
			args = append(args, Format(rand.Int()))
		default:
			args = append(args, v)
		}
	}
	return Format(md5.Sum([]byte(strings.Join(args, "")))), args
}
func HashsPath(arg ...Any) string {
	h := Hashs(arg...)
	return path.Join(h[:2], h)
}
func Hashs(arg ...Any) string {
	if len(arg) > 0 {
		switch arg := arg[0].(type) {
		case []byte:
			md := md5.New()
			md.Write(arg)
			return Format(md.Sum(nil))
		case io.Reader:
			md := md5.New()
			io.Copy(md, arg)
			return Format(md.Sum(nil))
		}
	}
	h, _ := Hash(arg...)
	return h
}
func Renders(str string, arg Any) string {
	if b, e := Render(str, arg); e != nil {
		panic(e)
	} else {
		return string(b)
	}
}
func Render(str string, arg Any) (b []byte, e error) {
	t := template.New("render").Funcs(template.FuncMap{
		"Format":  Format,
		"FmtSize": FmtSize,
		"Value":   Value,
		"Base":    func(p string) string { return path.Base(p) },
		"Capital": func(p string) string {
			if p != "" {
				return string(unicode.ToUpper(rune(p[0]))) + p[1:]
			}
			return path.Base(p)
		},
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
