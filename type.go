package kit

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"
)

func Int64(val interface{}) int64 {
	switch val := val.(type) {
	case int:
		return int64(val)
	case int64:
		return val
	case string:
		i, _ := strconv.ParseInt(val, 10, 64)
		return i
	case []interface{}:
		return int64(len(val))
	case map[string]interface{}:
		return int64(len(val))
	case time.Time:
		return val.Unix()
	}
	return 0
}
func Int(val interface{}) int {
	return int(Int64(val))
}
func Format(val interface{}) string {
	switch val := val.(type) {
	case nil:
		return ""
	case string:
		return val
	case []interface{}:
		if len(val) == 0 {
			return "[]"
		}
	case map[string]interface{}:
		if len(val) == 0 {
			return "{}"
		}
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
		if val != "" {
			return val
		}
	}
	return def
}
