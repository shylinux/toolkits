package miss

import (
	"fmt"
	"io/ioutil"
	"path"
	"time"

	kit "shylinux.com/x/toolkits"
	"shylinux.com/x/toolkits/conf"
	"shylinux.com/x/toolkits/file"
	"shylinux.com/x/toolkits/logs"
)

type Any = interface{}
type Map = map[string]Any

const (
	META = kit.MDB_META
	HASH = kit.MDB_HASH
	LIST = kit.MDB_LIST

	SHORT = kit.MDB_SHORT
	FIELD = kit.MDB_FIELD
	LIMIT = kit.MDB_LIMIT
	LEAST = kit.MDB_LEAST
	STORE = kit.MDB_STORE
	FSIZE = kit.MDB_FSIZE

	TIME = kit.MDB_TIME
	SIZE = kit.MDB_SIZE
	FILE = kit.MDB_FILE
)

const MISS = "miss"

type Miss struct {
	short int
	limit string
	least string
	store string
	fsize string

	Times  func(time.Time) string
	Hashs  func(string) string
	Logger func(...Any)

	file file.File
	conf *conf.Conf
}

func (miss *Miss) cache(cache Map, must bool) (Map, Map, []Any) {
	meta, ok := cache[META].(Map)
	if !ok && must {
		meta = Map{}
		cache[META] = meta
	}
	hash, ok := cache[HASH].(Map)
	if !ok && must {
		hash = Map{}
		cache[HASH] = hash
	}
	list, _ := cache[LIST].([]Any)
	return meta, hash, list
}
func (miss *Miss) meta(meta Map, key string) string {
	switch key {
	case LIMIT:
		return kit.Select(miss.limit, kit.Format(meta[key]))
	case LEAST:
		return kit.Select(miss.least, kit.Format(meta[key]))
	case STORE:
		return kit.Select(miss.store, kit.Format(meta[key]))
	case FSIZE:
		return kit.Select(miss.fsize, kit.Format(meta[key]))
	default:
		return kit.Format(meta[key])
	}
}
func (miss *Miss) data(data Any, arg ...Any) {
	meta := kit.GetMeta(kit.Dict(data))
	if kit.Value(meta, TIME) == nil {
		kit.Value(meta, TIME, miss.now())
	}
	kit.Value(meta, arg...)
}
func (miss *Miss) now() string { return miss.Times(time.Now()) }
func (miss *Miss) filename(meta Map, prefix, hash, ext string) string {
	if ext == CSV {
		if records, _ := meta[RECORDS].([]Any); len(records) > 0 {
			name := kit.Format(kit.Value(records, kit.Keys("-3", FILE)))
			if s, e := miss.file.StatFile(name); e == nil {
				if s.Size() > kit.Int64(miss.meta(meta, FSIZE)) {
					dir := path.Join(miss.meta(meta, STORE), prefix)
					name = path.Join(dir, kit.Keys(LIST+"_"+kit.Format(meta[OFFSET]), CSV))
				}
			}
			return name
		}
	}

	return path.Join(miss.meta(meta, STORE), prefix, kit.Keys(hash, ext))
}
func (miss *Miss) readfile(p string) ([]byte, error) {
	if f, e := miss.file.OpenFile(p); e == nil {
		defer f.Close()
		b, e := ioutil.ReadAll(f)
		miss.Logger("read file", p, SIZE, len(b), logs.FileLineMeta(kit.FileLine(2, 3)))
		return b, e
	} else {
		return nil, e
	}
}
func (miss *Miss) writefile(p string, data ...Any) error {
	if f, p, e := miss.file.CreateFile(p); e == nil {
		defer f.Close()
		n, e := fmt.Fprint(f, data...)
		miss.Logger("write file", p, SIZE, n, logs.FileLineMeta(kit.FileLine(2, 3)))
		return e
	} else {
		return e
	}
}
func New(conf *conf.Conf, file file.File) *Miss {
	return &Miss{
		short: conf.GetInt(SHORT, 6),
		store: conf.Get(STORE, "var/data"),
		fsize: conf.Get(FSIZE, "200000"),
		limit: conf.Get(LIMIT, "120"),
		least: conf.Get(LEAST, "30"),

		file: file, conf: conf,

		Times:  func(t time.Time) string { return t.Format(kit.MOD_TIME) },
		Hashs:  func(str string) string { return kit.Hashs(str) },
		Logger: logs.Logger(MISS),
	}
}
