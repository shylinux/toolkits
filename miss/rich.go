package miss

import (
	kit "github.com/shylinux/toolkits"
	"github.com/shylinux/toolkits/logs"

	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"sort"
	"time"
)

func (miss *Miss) Richs(prefix string, cache map[string]interface{}, raw interface{}, cb interface{}) (res map[string]interface{}) {
	// 数据结构
	meta, ok := cache[kit.MDB_META].(map[string]interface{})
	hash, ok := cache[kit.MDB_HASH].(map[string]interface{})
	if !ok {
		return nil
	}

	h := kit.Format(raw)
	switch h {
	case kit.MDB_FOREACH:
		// 全部遍历
		switch cb := cb.(type) {
		case func(string, string):
			for k, v := range hash {
				cb(k, kit.Format(v))
			}
		case func(string, map[string]interface{}):
			for k, v := range hash {
				res = v.(map[string]interface{})
				cb(k, res)
			}
		}
		return res
	case kit.MDB_RANDOMS:
		// 随机选取
		if len(hash) > 0 {
			list := []string{}
			for k := range hash {
				list = append(list, k)
			}
			h = list[rand.Intn(len(list))]
			res, _ = hash[h].(map[string]interface{})
		}
	default:
		// 单个查询
		if res, ok = hash[h].(map[string]interface{}); !ok {
			switch kit.Format(kit.Value(meta, kit.MDB_SHORT)) {
			case "", "uniq":
			default:
				hh := kit.Hashs(h)
				if res, ok = hash[hh].(map[string]interface{}); ok {
					h = hh
					break
				}

				dir := path.Join(kit.Select(miss.store, kit.Format(meta[kit.MDB_STORE])), prefix)
				for _, k := range []string{h, hh} {
					p := path.Join(dir, kit.Keys(k, "json"))
					if f, e := os.Open(p); e == nil {
						defer f.Close()
						if b, e := ioutil.ReadAll(f); e == nil {
							if json.Unmarshal(b, &res) == e {
								log.Show("miss", "import", p)
								h = k
								break
							}
						}
					}
				}
			}
		}
	}

	// 返回数据
	if res != nil {
		switch cb := cb.(type) {
		case func(map[string]interface{}):
			cb(res)
		case func(string, map[string]interface{}):
			cb(h, res)
		}
	}
	return res
}
func (miss *Miss) Rich(prefix string, cache map[string]interface{}, data interface{}) string {
	// 数据结构
	meta, ok := cache[kit.MDB_META].(map[string]interface{})
	if !ok {
		meta = map[string]interface{}{}
		cache[kit.MDB_META] = meta
	}
	hash, ok := cache[kit.MDB_HASH].(map[string]interface{})
	if !ok {
		hash = map[string]interface{}{}
		cache[kit.MDB_HASH] = hash
	}

	// 通用数据
	nest := kit.Select("", "meta.", kit.Value(data, "meta") != nil)
	if kit.Value(data, nest+kit.MDB_TIME) == nil {
		kit.Value(data, nest+kit.MDB_TIME, time.Now().Format("2006-01-02 15:03:04"))
	}

	// 生成键值
	h := ""
	switch short := kit.Format(kit.Value(meta, kit.MDB_SHORT)); short {
	case "":
		h = kit.ShortKey(hash, 6)
	case "uniq":
		h = kit.Hashs("uniq")
	case "data":
		h = kit.Hashs(kit.Format(data))
	default:
		if kit.Value(data, "meta") != nil {
			h = kit.Hashs(kit.Format(kit.Value(data, "meta."+short)))
		} else {
			h = kit.Hashs(kit.Format(kit.Value(data, short)))
		}
	}

	// 添加数据
	if hash[h] = data; len(hash) >= kit.Int(kit.Select(miss.limit, kit.Format(meta[kit.MDB_LIMIT]))) {
		least := kit.Int(kit.Select(miss.least, kit.Format(meta[kit.MDB_LEAST])))
		store := kit.Select(miss.store, kit.Format(meta[kit.MDB_STORE]))

		// 时间淘汰
		list := []int{}
		for _, v := range hash {
			list = append(list, int(kit.Time(kit.Format(kit.Value(v, "time")))))
		}
		sort.Ints(list)
		dead := 0
		if len(list) > 0 {
			dead = list[len(list)-1-least]
		}

		dir := path.Join(store, prefix)
		for k, v := range hash {
			if int(kit.Time(kit.Format(kit.Value(v, "time")))) > dead {
				break
			}

			name := path.Join(dir, kit.Keys(k, "json"))
			if f, p, e := kit.Create(name); e == nil {
				defer f.Close()
				// 保存数据
				if n, e := f.WriteString(kit.Format(v)); e == nil {
					log.Show("miss", "export", p, kit.MDB_SIZE, n)
					delete(hash, k)
				}
			}
		}
	}

	return h
}
