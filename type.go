package kit

func Int64(arg interface{}) int64 {
	switch arg := arg.(type) {
	case int64:
		return arg
	case string:
	case []interface{}:
		return int64(len(arg))
	case map[string]interface{}:
		return int64(len(arg))
	}
	return 0
}
