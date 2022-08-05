package kit

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Any = interface{}
type Map = map[string]Any
type Maps = map[string]string

func Max(list ...int) (max int) {
	for i := 0; i < len(list); i++ {
		if i == 0 {
			max = list[i]
			continue
		}
		if list[i] > max {
			max = list[i]
		}
	}
	return max
}
func Int(val Any) int {
	return int(Int64(val))
}
func Int64(val Any) int64 {
	switch val := val.(type) {
	case int:
		return int64(val)
	case int64:
		return val
	case float64:
		return int64(val)
	case string:
		i, _ := strconv.ParseInt(val, 10, 64)
		return i
	case []Any:
		return int64(len(val))
	case Map:
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
	case []Any:
		if len(val) == 0 {
			return "[]"
		}
	case Map:
		if len(val) == 0 {
			return "{}"
		}
	case []byte:
		return hex.EncodeToString(val[:])
	case [20]byte:
		return hex.EncodeToString(val[:])
	case time.Duration:
		return val.String()
	case time.Time:
		return val.Format("2006-01-02 15:04:05")
	default:
		if t := reflect.TypeOf(val); t.Kind() == reflect.Func {
			return LowerCapital(FuncName(val))
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
	for _, v := range arg {
		switch v := v.(type) {
		case Maps:
			for k, v := range v {
				defs[k] = v
			}
		case Map:
			for k, v := range v {
				defs[k] = v
			}
		default:
			args = append(args, Format(v))
		}
	}

	keys := Split(Select("type,name,text", key))
	for i, k := range keys {
		if v := Select(Format(defs[k]), args, i); v != "" {
			res = append(res, k, v)
		}
	}
	res = append(res, Slice(args, len(keys))...)
	return
}
func Simple(val ...Any) []string {
	res := []string{}
	for _, v := range val {
		switch val := v.(type) {
		case nil:
		case float64:
			res = append(res, fmt.Sprintf("%d", int64(val)))
		case []string:
			res = append(res, val...)
		case []Any:
			for _, v := range val {
				res = append(res, Simple(v)...)
			}
		case map[string]int:
			for k, v := range val {
				res = append(res, k, Format(v))
			}
		case Map:
			keys := []string{}
			for k := range val {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				res = append(res, k, Format(val[k]))
			}
		default:
			res = append(res, Format(val))
		}
	}
	return res
}
func Regexp(arg string) *regexp.Regexp {
	reg, _ := regexp.Compile(arg)
	return reg
}
func Duration(str Any) time.Duration {
	switch str := str.(type) {
	case string:
		d, _ := time.ParseDuration(str)
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

	for _, v := range []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"01-02 15:04",
		"2006/01/02",
		"2006-01-02",
		"2006-01",
		"15:04:05",
		"15:04",
	} {
		if t, e := time.ParseInLocation(v, arg[0], time.Local); e == nil {
			return Int64(t)
		}
	}
	return 0
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
		return fmt.Sprintf("%s%0.2fs", sign, float64(time)/1000000000)
	}
	if time > 1000000 {
		return fmt.Sprintf("%s%0.2fms", sign, float64(time)/1000000)
	}
	if time > 1000 {
		return fmt.Sprintf("%s%0.2fus", sign, float64(time)/1000)
	}
	return fmt.Sprintf("%s%dns", sign, time)
}

func Contains(str, sub Any) bool {
	return strings.Contains(Format(str), Format(sub))
}
func Capital(str string) string {
	return strings.ToUpper(str[0:1]) + str[1:]
}
func LowerCapital(str string) string {
	return strings.ToLower(str[0:1]) + str[1:]
}
func Select(def string, arg ...Any) string {
	if len(arg) == 0 {
		return def
	}

	switch val := arg[0].(type) {
	case []string:
		index := 0
		if len(arg) > 1 {
			switch v := arg[1].(type) {
			case int:
				index = v
			}
		}
		if index >= 0 && index < len(val) && val[index] != "" {
			return val[index]
		}
	case string:
		if len(arg) > 1 {
			switch v := arg[1].(type) {
			case bool:
				if v && val != "" {
					return val
				}
				return def
			}
		}
		if val != "" {
			return val
		}
	}
	return def
}
func Width(str string, mul int) int {
	return len([]rune(str)) + (len(str)-len([]rune(str)))/2/mul
}

func Replace(str string, arg ...string) string {
	for i := 0; i < len(arg); i += 2 {
		str = strings.Replace(str, arg[i], arg[i+1], 1)
	}
	return str
}
func ReplaceAll(str string, arg ...string) string {
	for i := 0; i < len(arg); i += 2 {
		str = strings.ReplaceAll(str, arg[i], arg[i+1])
	}
	return str
}
func IndexOf(str []string, sub string) int {
	for i, v := range str {
		if v == sub {
			return i
		}
	}
	return -1
}
func ForEach(arg []string, cb func(string)) {
	for _, v := range arg {
		cb(v)
	}
}
func Revert(str []string) []string {
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
	if begin > len(list) {
		begin = len(list)
	}

	end := len(list)
	if len(index) > 1 {
		if index[1] < 0 {
			end = index[1] + len(list)
		} else if index[1] < end {
			end = index[1]
		}
	}

	step := 1
	if len(index) > 2 {
		step = index[3]
	}
	for ; end > 0; end-- {
		if list[end-1] != "" {
			break
		}
	}
	if begin > end {
		begin = end
	}
	if step == 1 {
		return list[begin:end]
	}
	return nil
}
func Join(str []string, arg ...string) string {
	return strings.Join(str, Select(",", arg, 0))
}
