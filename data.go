package kit

import (
	"strings"
)

func _merge(meta Map, arg ...Any) Map {
	if len(arg) == 1 {
		switch arg := arg[0].(type) {
		case string:
			data, _ := UnMarshal(arg).(Map)
			return data
		case []string:
			if len(arg) == 1 {
				data, _ := UnMarshal(arg[0]).(Map)
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
		case []Any:
			for i := 0; i < len(args)-1; i += 2 {
				Value(meta, args[i], args[i+1])
			}
		case Map:
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
			if i < len(arg)-1 {
				Value(meta, arg[i], arg[i+1])
			}
			i++
		}
	}
	return meta
}
func Dict(arg ...Any) Map {
	if len(arg) == 1 {
		switch v := arg[0].(type) {
		case nil:
			return Map{}
		case Map:
			return v
		case string:
			res, _ := UnMarshal(v).(Map)
			return res
		case []byte:
			res, _ := UnMarshal(string(v)).(Map)
			return res
		}
	}
	return _merge(Map{}, arg...)
}
func Data(arg ...Any) Map {
	return Map{MDB_META: _merge(Map{}, arg...)}
}
func List(arg ...Any) []Any {
	if len(arg) == 0 {
		return []Any{}
	}
	if len(arg) == 1 {
		if arg[0] == nil {
			return []Any{}
		}
		if list, ok := arg[0].([]Any); ok {
			return list
		}
		return arg
	}
	if arg[0] != MDB_TYPE {
		res := []Any{}
		for _, v := range arg {
			switch v := v.(type) {
			case []Any:
				res = append(res, v...)
			default:
				res = append(res, v)
			}
		}
		return res
	}
	list, data := List(), Dict()
	for i := 0; i < len(arg)-1; i += 2 {
		if arg[i] == MDB_TYPE {
			data = Map{}
			list = append(list, data)
		} else if i == 0 {
			return arg
		}
		Value(data, arg[i], arg[i+1])
	}
	return list
}

func Keys(arg ...Any) string {
	return strings.TrimSuffix(strings.TrimPrefix(strings.Join(Simple(arg...), "."), "."), ".")
}
func Keym(arg ...Any) string {
	return Keys(MDB_META, Keys(arg))
}
func Keycb(arg ...Any) string {
	return Keys(Keys(arg), "cb")
}
func KeyHash(arg ...Any) string {
	return Keys(MDB_HASH, Hashs(arg[0]), arg[1:])
}
func KeyExtra(arg ...Any) string {
	return Keys(MDB_EXTRA, Keys(arg...))
}
func Fields(arg ...Any) string { return Join(Simple(arg...)) }
