package kit

import (
	"strings"
)

const (
	MIME_FORM = "application/x-www-form-urlencoded"
	MIME_FILE = "application/octet-stream"
	MIME_JSON = "application/json"
	MIME_TEXT = "text/plain"
	MIME_HTML = "text/html"
	MIME_PNG  = "image/png"
)

const (
	META_SOURCE = "meta.source"
	META_FIELD  = "meta.field"
	META_PATH   = "meta.path"
)

const (
	// 结构
	MDB_DICT = "dict"
	MDB_META = "meta"
	MDB_HASH = "hash"
	MDB_LIST = "list"

	MDB_FOREACH = "*"
	MDB_RANDOMS = "%"

	// 存储
	MDB_INPUT = "_input"
	MDB_SHORT = "short"
	MDB_FIELD = "field"
	MDB_STORE = "store"
	MDB_FSIZE = "fsize"
	MDB_TOTAL = "total"
	MDB_COUNT = "count"
	MDB_LIMIT = "limit"
	MDB_LEAST = "least"
	MDB_TABLE = "table"
	MDB_INDEX = "index"

	// 显示
	MDB_TEMPLATE = "template"
	MDB_CONTENT  = "content"
	// MDB_DISPLAY  = "display"
	MDB_ACTION = "action"
	MDB_BUTTON = "button"
	MDB_TITLE  = "title"
	MDB_TRANS  = "trans"
	MDB_TOPIC  = "topic"
	MDB_STYLE  = "style"
	MDB_COLOR  = "color"
	MDB_IMAGE  = "image"

	// 程序
	MDB_DOMAIN = "domain"
	MDB_PREFIX = "prefix"
	MDB_SCRIPT = "script"
	MDB_STATUS = "status"
	MDB_STREAM = "stream"
	MDB_EXPIRE = "expire"
	MDB_REGEXP = "regexp"
	MDB_EVENT  = "event"
	MDB_ERROR  = "error"
	MDB_EXTRA  = "extra"
	MDB_VALUE  = "value"
	MDB_CODE   = "code"
	MDB_BACK   = "back"

	// 文件
	MDB_PROTO = "proto"
	MDB_LINK  = "link"
	MDB_PATH  = "path"
	MDB_FILE  = "file"
	MDB_LINE  = "line"
	MDB_DATA  = "data"
	MDB_SIZE  = "size"
	MDB_MAKE  = "make"
	MDB_FROM  = "from"
	MDB_MAIN  = "main"

	// 数据
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
func Keys(arg ...interface{}) string {
	return strings.TrimSuffix(strings.TrimPrefix(strings.Join(Simple(arg...), "."), "."), ".")
}
func Keym(arg ...interface{}) string {
	return Keys(MDB_META, Keys(arg))
}
func KeyHash(arg ...interface{}) string {
	return Keys(MDB_HASH, Hashs(arg[0]))
}
func Keycb(arg ...interface{}) string {
	return Keys(Keys(arg), "cb")
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

func GetMeta(value map[string]interface{}) map[string]interface{} {
	if value[MDB_META] != nil {
		value = value[MDB_META].(map[string]interface{})
	}
	return value
}
