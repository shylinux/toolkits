package kit

import (
	"strings"
)

const (
	MIME_FILE = "application/octet-stream"
	MIME_JSON = "application/json"
	MIME_TEXT = "text/plain"
	MIME_HTML = "text/html"

	MIME_TABLE = "table"
	MIME_ORDER = "order"
	MIME_CHAIN = "chain"
	MIME_TITLE = "title"
	MIME_SHELL = "shell"

	MIME_SPACE = "space"
	MIME_STORY = "story"
	MIME_FAVOR = "favor"
	MIME_SHARE = "share"

	MIME_MASTER = "master"
	MIME_MYSELF = "myself"
	MIME_SERVER = "server"
	MIME_WORKER = "worker"
)

const (
	MDB_SHORT = "short"

	MDB_META = "meta"
	MDB_LIST = "list"
	MDB_HASH = "hash"

	MDB_STATUS = "status"
	MDB_EXTRA  = "extra"
	MDB_VALUE  = "value"
	MDB_USER   = "user"

	MDB_SIZE = "size"
	MDB_FILE = "file"
	MDB_TEXT = "text"
	MDB_TYPE = "type"
	MDB_TIME = "time"
	MDB_NICK = "nick"
	MDB_NAME = "name"
	MDB_KEY  = "key"
	MDB_ID   = "id"

	MDB_INPUT = "_input"
)

func Keys(arg ...interface{}) string {
	return strings.Join(Simple(arg...), ".")
}
func Data(arg ...interface{}) map[string]interface{} {
	meta := map[string]interface{}{}
	data := map[string]interface{}{
		MDB_META: meta, MDB_LIST: []interface{}{}, MDB_HASH: map[string]interface{}{},
	}
	for i := 0; i < len(arg)-1; i += 2 {
		Value(meta, arg[i], arg[i+1])
	}
	return data
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
func Dict(arg ...interface{}) map[string]interface{} {
	dict := map[string]interface{}{}
	for i := 0; i < len(arg)-1; i += 2 {
		Value(dict, arg[i], arg[i+1])
	}
	return dict
}
