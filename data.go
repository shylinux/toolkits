package kit

import (
	"strings"
)

func _merge(meta map[string]interface{}, arg ...interface{}) map[string]interface{} {
	if len(arg) == 1 {
		switch arg := arg[0].(type) {
		case string:
			data, _ := UnMarshal(arg).(map[string]interface{})
			return data
		case []string:
			if len(arg) == 1 {
				data, _ := UnMarshal(arg[0]).(map[string]interface{})
				return data
			}
		}
	}

	for i := 0; i < len(arg); i++ {
		switch args := arg[i].(type) {
		case []string:
			for i := 0; i < len(args)-1; i += 2 {
				Value(meta, args[i], args[i+1])
			}
		case []interface{}:
			for i := 0; i < len(args)-1; i += 2 {
				Value(meta, args[i], args[i+1])
			}
		case map[string]interface{}:
			for k, v := range args {
				if Value(meta, k) == nil {
					Value(meta, k, v)
				}
			}
		case map[string]string:
			for k, v := range args {
				Value(meta, k, v)
			}
		default:
			Value(meta, arg[i], arg[i+1])
			i++
		}
	}
	return meta
}
func Dict(arg ...interface{}) map[string]interface{} {
	return _merge(map[string]interface{}{}, arg...)
}
func Data(arg ...interface{}) map[string]interface{} {
	return map[string]interface{}{
		MDB_META: _merge(map[string]interface{}{}, arg...),
		MDB_HASH: map[string]interface{}{},
		MDB_LIST: []interface{}{},
	}
}
func List(arg ...interface{}) []interface{} {
	if len(arg) == 0 || arg[0] != MDB_TYPE {
		return arg
	}
	list, data := List(), Dict()
	for i := 0; i < len(arg)-1; i += 2 {
		if arg[i] == MDB_TYPE {
			data = map[string]interface{}{}
			list = append(list, data)
		} else if i == 0 {
			return arg
		}
		Value(data, arg[i], arg[i+1])
	}
	return list
}

func Keys(arg ...interface{}) string {
	return strings.TrimSuffix(strings.TrimPrefix(strings.Join(Simple(arg...), "."), "."), ".")
}
func Keym(arg ...interface{}) string {
	return Keys(MDB_META, Keys(arg))
}
func Keycb(arg ...interface{}) string {
	return Keys(Keys(arg), "cb")
}
func KeyHash(arg ...interface{}) string {
	return Keys(MDB_HASH, Hashs(arg[0]), arg[1:])
}
func KeyExtra(arg ...interface{}) string {
	return Keys(MDB_EXTRA, Keys(arg...))
}
