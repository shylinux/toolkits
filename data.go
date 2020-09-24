package kit

import "strings"

const (
	MIME_FILE = "application/octet-stream"
	MIME_JSON = "application/json"
	MIME_TEXT = "text/plain"
	MIME_HTML = "text/html"
)

const (
	META_SOURCE = "meta.source"
	META_FIELD  = "meta.field"
	META_PATH   = "meta.path"
)

const (
	MDB_FOREACH = "*"
	MDB_RANDOMS = "%"

	MDB_SHORT = "short"
	MDB_FIELD = "field"
	MDB_STORE = "store"
	MDB_FSIZE = "fsize"
	MDB_LIMIT = "limit"
	MDB_LEAST = "least"

	MDB_DICT = "dict"
	MDB_META = "meta"
	MDB_HASH = "hash"
	MDB_LIST = "list"

	MDB_ACTION = "action"
	MDB_STATUS = "status"
	MDB_ERROR  = "error"
	MDB_EXTRA  = "extra"
	MDB_VALUE  = "value"
	MDB_TOPIC  = "topic"
	MDB_STYLE  = "style"

	MDB_TOTAL = "total"
	MDB_COUNT = "count"
	MDB_STEP  = "step"

	MDB_ROUTE = "route"
	MDB_REPOS = "repos"
	MDB_GROUP = "group"
	MDB_USER  = "user"

	MDB_PORT = "port"
	MDB_HOST = "host"
	MDB_DIR  = "dir"
	MDB_ENV  = "env"
	MDB_ARG  = "arg"
	MDB_PID  = "pid"

	MDB_LINK = "link"
	MDB_SIZE = "size"
	MDB_LINE = "line"
	MDB_FILE = "file"
	MDB_PATH = "path"

	MDB_HELP = "help"
	MDB_TEXT = "text"
	MDB_NAME = "name"
	MDB_TYPE = "type"
	MDB_ZONE = "zone"
	MDB_TIME = "time"
	MDB_KEY  = "key"
	MDB_ID   = "id"
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

	for i := 0; i < len(arg); i += 2 {
		if i == len(arg)-1 {
			switch arg := arg[i].(type) {
			case []string:
				for i := 0; i < len(arg)-1; i += 2 {
					Value(meta, arg[i], arg[i+1])
				}
			case []interface{}:
				for i := 0; i < len(arg)-1; i += 2 {
					Value(meta, arg[i], arg[i+1])
				}
			case map[string]interface{}:
				for k, v := range arg {
					Value(meta, k, v)
				}
			}
		} else {
			Value(meta, arg[i], arg[i+1])
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
func Keys(arg ...interface{}) string {
	return strings.TrimSuffix(strings.TrimPrefix(strings.Join(Simple(arg...), "."), "."), ".")
}
func List(arg ...interface{}) []interface{} {
	list, data := []interface{}{}, map[string]interface{}{}
	for i := 0; i < len(arg)-1; i += 2 {
		if arg[i] == MDB_INPUT {
			data = map[string]interface{}{}
			list = append(list, data)
		}
		Value(data, arg[i], arg[i+1])
	}
	return list
}
