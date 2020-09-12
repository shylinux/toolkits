package miss

import (
	kit "github.com/shylinux/toolkits"
	log "github.com/shylinux/toolkits/logs"

	"encoding/csv"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

func (miss *Miss) Grow(prefix string, cache map[string]interface{}, data interface{}) int {
	// 数据结构
	meta, ok := cache[kit.MDB_META].(map[string]interface{})
	if !ok {
		meta = map[string]interface{}{}
		cache[kit.MDB_META] = meta
	}
	list, _ := cache[kit.MDB_LIST].([]interface{})

	// 通用数据
	id := kit.Int(meta[kit.MDB_COUNT]) + 1
	nest := kit.Select("", "meta.", kit.Value(data, "meta") != nil)
	if kit.Value(data, nest+kit.MDB_ID, id); kit.Value(data, nest+kit.MDB_TIME) == nil {
		kit.Value(data, nest+kit.MDB_TIME, kit.Select(time.Now().Format("2006-01-02 15:03:04")))
	}

	// 添加数据
	list = append(list, data)
	cache[kit.MDB_LIST] = list
	meta[kit.MDB_COUNT] = id

	// 保存数据
	if len(list) >= kit.Int(kit.Select(miss.limit, kit.Format(meta[kit.MDB_LIMIT]))) {
		least := kit.Int(kit.Select(miss.least, kit.Format(meta[kit.MDB_LEAST])))
		store := kit.Select(miss.store, kit.Format(meta[kit.MDB_STORE]))

		record, _ := meta["record"].([]interface{})

		// 文件命名
		dir := path.Join(store, prefix)
		name := path.Join(dir, kit.Keys(prefix, "csv"))
		if len(record) > 0 {
			name = kit.Format(kit.Value(record, kit.Keys(len(record)-1, "file")))
			if s, e := os.Stat(name); e == nil {
				if s.Size() > kit.Int64(kit.Select(kit.Format(kit.MDB_FSIZE), kit.Format(meta[kit.MDB_FSIZE]))) {
					name = path.Join(dir, kit.Keys(prefix+"_"+kit.Format(meta["offset"]), "csv"))
				}
			}
		}

		// 打开文件
		f, e := os.OpenFile(name, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if e != nil {
			f, _, e = kit.Create(name)
			log.Show("miss", "create", name)
		} else {
			log.Show("miss", "append", name)
		}
		defer f.Close()
		s, _ := f.Stat()

		// 保存表头
		keys := []string{}
		w := csv.NewWriter(f)
		if s.Size() == 0 {
			for k := range list[0].(map[string]interface{}) {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			w.Write(keys)
			log.Show("miss", "write head", keys)
			w.Flush()
			s, e = f.Stat()
		} else {
			r := csv.NewReader(f)
			keys, e = r.Read()
			log.Show("miss", "read head", keys)
		}

		// 创建索引
		count := len(list) - least
		offset := kit.Int(meta["offset"])
		meta["record"] = append(record, map[string]interface{}{
			"time": miss.now(), "offset": offset, "count": count,
			"file": name, "position": s.Size(),
		})

		// 保存数据
		for i, v := range list {
			if i >= count {
				break
			}

			val := v.(map[string]interface{})

			values := []string{}
			for _, k := range keys {
				values = append(values, kit.Format(val[k]))
			}
			w.Write(values)

			if i < least {
				list[i] = list[count+i]
			}
		}

		log.Show("miss", "save", name, "offset", offset, "count", count)
		meta["offset"] = offset + count
		list = list[count:]
		cache[kit.MDB_LIST] = list
		w.Flush()
	}
	return id
}
func (miss *Miss) Grows(prefix string, cache map[string]interface{}, offend, limit int, match string, value string, cb interface{}) map[string]interface{} {
	// 数据结构
	meta, ok := cache[kit.MDB_META].(map[string]interface{})
	list, ok := cache[kit.MDB_LIST].([]interface{})
	if !ok {
		return nil
	}

	// 数据范围
	current := kit.Int(meta["offset"])
	end := current + len(list) - offend
	begin := end - limit
	switch limit {
	case -1:
		begin = current
	case -2:
		begin = 0
	}

	if match == kit.MDB_ID && value != "" {
		begin, end = kit.Int(value)-1, kit.Int(value)
		match, value = "", ""
	}

	order := 0
	if begin < current {
		// 读取文件
		// m.Log(LOG_INFO, "%s.%v read %v-%v from %v-%v", key, chain, begin, end, current, current+len(list))
		store, _ := meta["record"].([]interface{})
		for s := len(store) - 1; s > -1; s-- {
			item, _ := store[s].(map[string]interface{})
			line := kit.Int(item["offset"])
			log.Show("miss", "action", "check", "record", s, "offset", line, "count", item["count"])
			if begin < line && s > 0 {
				if kit.Int(item["count"]) != 0 {
					s -= (line - begin) / kit.Int(item["count"])
				}
				// 向后查找
				continue
			}

			for ; begin < end && s < len(store); s++ {
				item, _ := store[s].(map[string]interface{})
				name := kit.Format(item["file"])
				pos := kit.Int(item["position"])
				offset := kit.Int(item["offset"])
				if offset+kit.Int(item["count"]) <= begin {
					log.Show("miss", "action", "check", "record", s, "offset", line, "count", item["count"])
					// 向前查找
					continue
				}

				if f, e := os.Open(name); e == nil {
					defer f.Close()
					// 打开文件
					r := csv.NewReader(f)
					heads, _ := r.Read()
					log.Show("miss", "import head", heads)

					f.Seek(int64(pos), os.SEEK_SET)
					r = csv.NewReader(f)
					for i := offset; i < end; i++ {
						lines, e := r.Read()
						if e != nil {
							log.Show("miss", "import line", e)
							break
						}
						if i < begin {
							continue
						}

						// 读取数据
						item := map[string]interface{}{}
						for i := range heads {
							if heads[i] == "extra" {
								item[heads[i]] = kit.UnMarshal(lines[i])
							} else {
								item[heads[i]] = lines[i]
							}
						}
						log.Show("miss", "offset", i, "type", item["type"], "name", item["name"], "text", item["text"])

						if match == "" || strings.Contains(kit.Format(item[match]), value) {
							// 匹配成功
							switch cb := cb.(type) {
							case func(int, map[string]interface{}):
								cb(order, item)
							case func(int, map[string]interface{}) bool:
								if cb(order, item) {
									return meta
								}
							}
							order++
						}
						begin = i + 1
					}
				}
			}
			break
		}
	}

	if begin < current {
		begin = current
	}
	for i := begin - current; i < end-current; i++ {
		// 读取缓存
		if match == "" || i < len(list) && strings.Contains(kit.Format(kit.Value(list[i], match)), value) {
			switch cb := cb.(type) {
			case func(int, map[string]interface{}):
				cb(order, list[i].(map[string]interface{}))
			case func(int, map[string]interface{}) bool:
				if cb(order, list[i].(map[string]interface{})) {
					return meta
				}
			}
			order++
		}
	}
	return meta
}
