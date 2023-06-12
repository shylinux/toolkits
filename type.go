package kit

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Any = interface{}
type Map = map[string]Any
type Maps = map[string]string

func Min(arg ...int) (min int) {
	for i, v := range arg {
		if i == 0 || v < min {
			min = v
		}
	}
	return
}
func Max(arg ...int) (max int) {
	for i, v := range arg {
		if i == 0 || v > max {
			max = v
		}
	}
	return
}
func Int(val Any) int { return int(Int64(val)) }
func Int64(val Any) int64 {
	switch val := val.(type) {
	case int:
		return int64(val)
	case int64:
		return val
	case float64:
		return int64(val)
	case string:
		val, unit := strings.ToLower(val), int64(1)
		If(strings.HasSuffix(val, "k") || strings.HasSuffix(val, "kb"), func() { val, unit = strings.Split(val, "k")[0], 1000 })
		If(strings.HasSuffix(val, "m") || strings.HasSuffix(val, "mb"), func() { val, unit = strings.Split(val, "k")[0], 1000000 })
		If(strings.HasSuffix(val, "g") || strings.HasSuffix(val, "gb"), func() { val, unit = strings.Split(val, "g")[0], 1000000000 })
		if i, e := strconv.ParseInt(val, 10, 64); e == nil {
			return i * unit
		}
		return 0
	case Map:
		return int64(len(val))
	case []Any:
		return int64(len(val))
	case time.Time:
		return val.UnixNano()
	case time.Duration:
		return val.Nanoseconds()
	}
	return 0
}
func Format(val Any, arg ...Any) string {
	switch val := val.(type) {
	case nil:
		return ""
	case string:
		if len(arg) > 0 {
			return fmt.Sprintf(val, arg...)
		}
		return val
	case Map:
		if len(val) == 0 {
			return "{}"
		}
	case []Any:
		if len(val) == 0 {
			return "[]"
		}
	case []byte:
		return hex.EncodeToString(val[:])
	case [16]byte:
		return hex.EncodeToString(val[:])
	case [20]byte:
		return hex.EncodeToString(val[:])
	case time.Time:
		return val.Format(MOD_TIME)
	case time.Duration:
		return val.String()
	case error:
		return val.Error()
	default:
		if t := reflect.TypeOf(val); t.Kind() == reflect.Func {
			return FuncName(val)
		}
	}
	b, _ := json.Marshal(val)
	return string(b)
}
func Formats(val Any) string {
	switch val := val.(type) {
	case nil:
		return ""
	case string:
		return val
	}
	b, _ := json.MarshalIndent(val, "", "  ")
	return string(b)
}
func SimpleKV(key string, arg ...Any) (res []string) {
	defs, args := Dict(), []string{}
	For(arg, func(v Any) {
		switch v := v.(type) {
		case Maps:
			For(v, func(k string, v string) { defs[k] = v })
		case Map:
			For(v, func(k string, v Any) { defs[k] = v })
		default:
			args = append(args, Simple(v)...)
		}
	})
	keys := Split(Select("type,name,text", key))
	for i, k := range keys {
		if v := Select(Format(defs[k]), args, i); v != "" {
			res = append(res, k, v)
		}
	}
	return append(res, Slice(args, len(keys))...)
}
func Simple(arg ...Any) (res []string) {
	for _, val := range arg {
		switch val := val.(type) {
		case nil:
		case []string:
			res = append(res, val...)
		case []Any:
			For(val, func(v Any) { res = append(res, Simple(v)...) })
		case Map:
			For(KeyValue(nil, "", val), func(k string, v string) { res = append(res, k, v) })
		case Maps:
			For(val, func(k string, v string) { res = append(res, k, v) })
		case map[string]int:
			For(val, func(k string, v int) { res = append(res, k, Format(v)) })
		case func(string) string:
			For(res, func(i int, v string) { res[i] = val(v) })
		case func(string, string) string:
			_res := []string{}
			For(val, func(k, v string) { _res = append(_res, val(k, v)) })
			res = _res
		case func(string) (string, error):
			For(res, func(i int, v string) { res[i], _ = val(v) })
		case func(string) bool:
			_res := []string{}
			For(res, func(v string) { If(val(v), func() { _res = append(_res, v) }) })
			res = _res
		case func(string):
			For(res, val)
		default:
			res = append(res, Format(val))
		}
	}
	return res
}
func Duration(val Any) time.Duration {
	switch val := val.(type) {
	case time.Duration:
		return val
	case string:
		d, _ := time.ParseDuration(val)
		return d
	}
	return time.Millisecond
}
func Time(arg ...string) int64 {
	if len(arg) == 0 {
		return Int64(time.Now())
	}
	if len(arg) > 1 {
		if t, e := time.ParseInLocation(arg[1], arg[0], time.Local); e == nil {
			return Int64(t)
		}
	}
	for _, v := range []string{MOD_DATE, MOD_TIME, MOD_TIMES} {
		if t, e := time.ParseInLocation(v, arg[0], time.Local); e == nil {
			return Int64(t)
		}
	}
	return Int64(time.Now())
}
func FmtSize(size int64) string {
	if size > 1<<30 {
		return fmt.Sprintf("%0.2fG", float64(size)/(1<<30))
	}
	if size > 1<<20 {
		return fmt.Sprintf("%0.2fM", float64(size)/(1<<20))
	}
	if size > 1<<10 {
		return fmt.Sprintf("%0.2fK", float64(size)/(1<<10))
	}
	return fmt.Sprintf("%dB", size)
}
func FmtDuration(t time.Duration) string {
	sign, tt := "", t
	If(tt < 0, func() { sign, tt = "-", -t })
	list, unit := []string{}, 24*3600*time.Second
	if tt > unit {
		list, tt = append(list, fmt.Sprintf("%s%dd", sign, tt/unit)), tt%unit
	}
	if unit = 3600 * time.Second; tt > unit {
		list, tt = append(list, fmt.Sprintf("%s%dh", sign, tt/unit)), tt%unit
	}
	if unit = 60 * time.Second; tt > unit {
		list, tt = append(list, fmt.Sprintf("%s%dm", sign, tt/unit)), tt%unit
	}
	if len(list) > 0 {
		if unit = time.Second; tt > unit {
			list, tt = append(list, fmt.Sprintf("%s%ds", sign, tt/unit)), tt%unit
		}
		return strings.Join(list, "")
	}
	if unit = time.Second; tt > unit {
		return fmt.Sprintf("%s%0.2fs", sign, float64(tt)/float64(unit))
	}
	if unit = time.Millisecond; tt > unit {
		return fmt.Sprintf("%s%0.2fms", sign, float64(tt)/float64(unit))
	}
	if unit = time.Microsecond; tt > unit {
		return fmt.Sprintf("%s%0.2fus", sign, float64(tt)/float64(unit))
	}
	return fmt.Sprintf("%s%dns", sign, tt)
}

func Select(def string, arg ...Any) string {
	if len(arg) == 0 {
		return def
	}
	switch val := arg[0].(type) {
	case []string:
		i := 0
		if len(arg) > 1 {
			switch v := arg[1].(type) {
			case int:
				i = v
			}
		}
		if i < 0 && 0 <= i+len(val) && i+len(val) < len(val) {
			return val[i+len(val)]
		}
		if 0 <= i && i < len(val) && val[i] != "" {
			return val[i]
		}
	case string:
		if len(arg) > 1 {
			switch v := arg[1].(type) {
			case bool:
				if v && val != "" {
					return val
				}
				return def
			case string:
				args := Simple(arg)
				for i := len(args) - 1; i >= 0; i-- {
					if args[i] != "" {
						return args[i]
					}
				}
			}
		}
		if val != "" {
			return val
		}
	}
	return def
}
func IsUpper(str string) bool        { return strings.ToUpper(str) == str }
func Capital(str string) string      { return strings.ToUpper(str[0:1]) + str[1:] }
func LowerCapital(str string) string { return strings.ToLower(str[0:1]) + str[1:] }
func Width(str string, mul int) int  { return len([]rune(str)) + (len(str)-len([]rune(str)))/2/mul }
func Sort(list []string) []string {
	sort.Strings(list)
	return list
}
func Contains(str string, arg ...string) bool {
	for _, v := range arg {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}
func HasPrefix(str string, arg ...string) bool {
	for _, v := range arg {
		if strings.HasPrefix(str, v) {
			return true
		}
	}
	return false
}
func HasSuffix(str string, arg ...string) bool {
	for _, v := range arg {
		if strings.HasSuffix(str, v) {
			return true
		}
	}
	return false
}
func ReplaceAll(str string, arg ...string) string {
	For(arg, func(from, to string) { str = strings.Replace(str, from, to, -1) })
	return str
}
func Replace(str string, arg ...string) string {
	For(arg, func(from, to string) { str = strings.Replace(str, from, to, 1) })
	return str
}

func IsIn(str string, arg ...string) bool { return IndexOf(arg, str) > -1 }
func HasPrefixList(args []string, arg ...string) bool {
	for i, v := range arg {
		if len(args) < i+1 || args[i] != v {
			return false
		}
	}
	return true
}
func IndexOf(str []string, sub string) int {
	for i, v := range str {
		if v == sub {
			return i
		}
	}
	return -1
}
func Filters(list []string, arg ...string) (res []string) {
	for _, v := range list {
		if IndexOf(arg, v) == -1 {
			res = append(res, v)
		}
	}
	return res
}
func Reverse(str []string) []string {
	for i := 0; i < len(str)/2; i++ {
		str[i], str[len(str)-1-i] = str[len(str)-1-i], str[i]
	}
	return str
}
func Slice(list []string, index ...int) []string {
	if len(list) == 0 {
		return []string{}
	}
	begin := 0
	if len(index) > 0 {
		if index[0] < 0 {
			begin = index[0] + len(list)
		} else {
			begin = index[0]
		}
	}
	If(begin > len(list), func() { begin = len(list) })
	end := len(list)
	if len(index) > 1 {
		if index[1] < 0 {
			end = index[1] + len(list)
		} else if index[1] < end {
			end = index[1]
		}
	}
	step := 1
	If(len(index) > 2, func() { step = index[3] })
	for ; end > 0; end-- {
		if list[end-1] != "" {
			break
		}
	}
	If(begin > end, func() { begin = end })
	if step == 1 {
		return list[begin:end]
	}
	return nil
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
func AddUniq(list []string, arg ...string) []string {
	For(arg, func(k string) { If(IndexOf(list, k) == -1, func() { list = append(list, k) }) })
	return list
}
func Join(str []string, arg ...string) string {
	return strings.Join(str, Select(FS, arg, 0))
}
func JoinKV(inner, outer string, arg ...string) string {
	res := []string{}
	for i := 0; i < len(arg)-1; i += 2 {
		if i == len(arg)-1 {
			res = append(res, arg[i])
			continue
		} else if arg[i+1] == "" {
			continue
		}
		res = append(res, arg[i]+inner+arg[i+1])
	}
	return strings.Join(res, outer)
}
