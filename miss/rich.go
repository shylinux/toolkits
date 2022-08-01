package miss

import (
	"math/rand"
	"sort"

	kit "shylinux.com/x/toolkits"
)

const (
	UNIQ    = kit.MDB_UNIQ
	DATA    = kit.MDB_DATA
	FOREACH = kit.MDB_FOREACH
	RANDOMS = kit.MDB_RANDOMS
)

const JSON = "json"

func (miss *Miss) Richs(prefix string, cache Map, raw Any, cb Any) (res Map) {
	meta, hash, _ := miss.cache(cache, false)
	if hash == nil || len(hash) == 0 { // 没有数据
		return nil
	}

	h := kit.Format(raw)
	switch h {
	case FOREACH: // 全部遍历
		switch cb := cb.(type) {
		case func(string, Map):
			for k, v := range hash {
				cb(k, v.(Map))
			}
		case func(string, string):
			for k, v := range hash {
				cb(k, kit.Format(v))
			}
		}
		return res
	case RANDOMS: // 随机选取
		if len(hash) > 0 {
			list := []string{}
			for k := range hash {
				list = append(list, k)
			}
			h = list[rand.Intn(len(list))]
			res, _ = hash[h].(Map)
		}
	default:
		var ok bool
		if res, ok = hash[h].(Map); ok {
			break // 键值查询
		}

		switch miss.meta(meta, SHORT) {
		case "", UNIQ: // 查询失败
		default:
			hh := miss.Hashs(h)
			if res, ok = hash[hh].(Map); ok {
				h = hh
				break // 哈希查询
			}

			for _, k := range []string{h, hh} {
				if b, e := miss.readfile(miss.filename(meta, prefix, k, JSON)); e == nil {
					res, h = kit.Dict(b), k
					break // 磁盘查询
				}
			}
		}
	}

	if res != nil { // 同步回调
		switch cb := cb.(type) {
		case func(Map):
			cb(res)
		case func(string, Map):
			cb(h, res)
		}
	}
	return res // 返回数据
}
func (miss *Miss) Rich(prefix string, cache Map, data Any) string {
	meta, hash, _ := miss.cache(cache, true)

	// 生成键值
	h := ""
	switch short := miss.meta(meta, SHORT); short {
	case "":
		h = kit.ShortKey(hash, miss.short)
	case UNIQ:
		h = miss.Hashs(UNIQ)
	case DATA:
		h = miss.Hashs(kit.Format(data))
	default:
		list := []string{}
		for _, k := range kit.Split(short) {
			list = append(list, kit.Format(kit.Value(kit.GetMeta(kit.Dict(data)), k)))
		}
		h = miss.Hashs(kit.Join(list))
	}

	// 添加数据
	miss.data(data)
	if old, ok := hash[h]; ok {
		for k, v := range data.(Map) {
			switch k {
			case META, HASH:
				for k1, v1 := range v.(Map) {
					kit.Value(old, kit.Keys(k, k1), v1)
				}
			case LIST:
			default:
				kit.Value(old, k, v)
			}
		}
	} else {
		hash[h] = data
	}

	if len(hash) < kit.Int(miss.meta(meta, LIMIT)) {
		return h // 直接返回
	}

	// 淘汰时间
	keys, list := map[string]int{}, []int{}
	kit.Fetch(hash, func(k string, v Map) {
		t := int(kit.Time(kit.Format(kit.Value(kit.GetMeta(v), TIME))))
		keys[k], list = t, append(list, t)
	})
	sort.Ints(list)

	dead := list[len(list)-1-kit.Int(miss.meta(meta, LEAST))]
	for k, t := range keys {
		if t < dead && miss.writefile(miss.filename(meta, prefix, k, JSON), kit.Format(hash[k])) == nil {
			delete(hash, k) // 淘汰数据
		}
	}
	return h
}
