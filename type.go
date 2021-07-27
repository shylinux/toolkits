package kit

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Int(val interface{}) int {
	return int(Int64(val))
}
func Int64(val interface{}) int64 {
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
	case []interface{}:
		return int64(len(val))
	case map[string]interface{}:
		return int64(len(val))
	case time.Time:
		return val.UnixNano()
	case time.Duration:
		return val.Nanoseconds()
	}
	return 0
}
func Format(val interface{}, arg ...interface{}) string {
	switch val := val.(type) {
	case nil:
		return ""
	case string:
		if len(arg) > 0 {
			return fmt.Sprintf(val, arg...)
		}
		return val
	case []interface{}:
		if len(val) == 0 {
			return "[]"
		}
	case map[string]interface{}:
		if len(val) == 0 {
			return "{}"
		}
	case time.Duration:
		return val.String()
	case time.Time:
		return val.Format("2006-01-02 15:04:05")
	}
	b, _ := json.Marshal(val)
	return string(b)
}
func Formats(val interface{}) string {
	switch val := val.(type) {
	case nil:
		return ""
	case string:
		return val
	}
	b, _ := json.MarshalIndent(val, "", "  ")
	return string(b)
}
func Duration(str interface{}) time.Duration {
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
func Width(str string, mul int) int {
	return len([]rune(str)) + (len(str)-len([]rune(str)))/2/mul
}
func Join(str []string, key string) string {
	return strings.Join(str, key)
}

func Select(def string, arg ...interface{}) string {
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
func Slice(list []string, index ...int) []string {
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
	if step == 1 {
		return list[begin:end]
	}
	return nil
}
func Simple(val ...interface{}) []string {
	res := []string{}
	for _, v := range val {
		switch val := v.(type) {
		case nil:
		case float64:
			res = append(res, fmt.Sprintf("%d", int64(val)))
		case []string:
			res = append(res, val...)
		case []interface{}:
			for _, v := range val {
				res = append(res, Simple(v)...)
			}
		case map[string]int:
			for k, v := range val {
				res = append(res, k, Format(v))
			}
		case map[string]interface{}:
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
func Revert(str []string) []string {
	for i := 0; i < len(str)/2; i++ {
		str[i], str[len(str)-1-i] = str[len(str)-1-i], str[i]
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
