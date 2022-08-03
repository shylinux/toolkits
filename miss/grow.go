package miss

import (
	kit "shylinux.com/x/toolkits"

	"encoding/csv"
	"sort"
	"strings"
)

const (
	RECORDS = "records"
	OFFSET  = "offset"

	COUNT = kit.MDB_COUNT
	EXTRA = kit.MDB_EXTRA
	ID    = kit.MDB_ID
)

const CSV = "csv"

func (miss *Miss) Grow(prefix string, cache Map, data Any) int {
	meta, _, list := miss.cache(cache, true)

	// 添加数据
	id := kit.Int(meta[COUNT]) + 1
	miss.data(data, ID, id)
	meta[COUNT] = id
	list = append(list, data)
	cache[LIST] = list

	if len(list) < kit.Int(miss.meta(meta, LIMIT)) {
		return id // 直接返回
	}

	// 打开文件
	p := miss.filename(meta, prefix, LIST, CSV)
	f, _ := miss.file.AppendFile(p)
	miss.Logger("open file", p)
	defer f.Close()

	// 保存表头
	head := []string{}
	w := csv.NewWriter(f)
	s, _ := miss.file.StatFile(p)
	if s.Size() > 0 {
		head, _ = csv.NewReader(f).Read()
		miss.Logger("read head", head)
	}
	if len(head) == 0 {
		if field := kit.Split(kit.Format(meta[FIELD])); len(field) > 0 {
			head = append(head, field...)
		} else {
			for k := range list[0].(Map) {
				head = append(head, k)
			}
			sort.Strings(head)
		}

		if s.Size() == 0 {
			w.Write(head)
			w.Flush()
			miss.Logger("write head", head)
			s, _ = miss.file.StatFile(p)
		}
	}

	// 创建索引
	least := kit.Int(miss.meta(meta, LEAST))
	offset, count := kit.Int(meta[OFFSET]), len(list)-least
	meta[RECORDS] = append(kit.List(meta[RECORDS]), Map{
		TIME: miss.now(), FILE: p, "position": s.Size(),
		OFFSET: offset, COUNT: count,
	})

	// 保存数据
	defer w.Flush()
	for i, v := range list {
		if i >= count {
			break
		}

		// 写出数据
		val := v.(Map)
		values := []string{}
		for _, k := range head {
			values = append(values, kit.Format(val[k]))
		}
		w.Write(values)

		// 移动数据
		if i < least {
			list[i] = list[count+i]
		}
	}

	miss.Logger("write data", p, OFFSET, offset, COUNT, count)
	meta[OFFSET], list = offset+count, list[count:]
	cache[LIST] = list
	return id
}
func (miss *Miss) Grows(prefix string, cache Map, offend, limit int, field string, value string, cb Any) (res Map) {
	meta, _, list := miss.cache(cache, false)
	if list == nil || len(list) == 0 {
		return nil
	}

	// 数据范围
	offset := kit.Int(meta[OFFSET])
	end := offset + len(list) - offend
	begin := end - limit
	switch limit {
	case -1:
		begin = offset
	case -2:
		begin = 0
	}
	if field == ID && value != "" {
		begin, end = kit.Int(value)-1, kit.Int(value)
		field, value = "", ""
	}

	// 读取数据
	index, done := 0, false
	if begin < offset && miss._grows_record(meta, begin, end, func(item Map) bool { // 读取文件
		res, index, done = _grow_match(item, field, value, index, cb)
		return done
	}) {
		return
	}
	if begin < offset {
		begin = offset
	}
	for i := begin - offset; i < len(list) && i < end-offset; i++ { // 读取缓存
		if res, index, done = _grow_match(kit.Dict(list[i]), field, value, index, cb); done {
			break
		}
	}
	return
}

func (miss *Miss) _grows_record(meta Map, begin, end int, cb func(Map) bool) bool {
	records := kit.List(meta[RECORDS])
	for i := len(records) - 1; i > -1; i-- {
		record := kit.Dict(records[i])
		offset := kit.Int(record[OFFSET])
		miss.Logger(RECORDS, i, OFFSET, offset, COUNT, record[COUNT])
		if begin < offset && i > 0 {
			if kit.Int(record[COUNT]) != 0 {
				i -= (offset - begin) / kit.Int(record[COUNT])
			}
			continue // 向后查找
		}

		for ; begin < end && i < len(records); i++ {
			record := kit.Dict(records[i])
			offset := kit.Int(record[OFFSET])
			miss.Logger(RECORDS, i, OFFSET, offset, COUNT, record[COUNT])
			if offset+kit.Int(record[COUNT]) <= begin {
				continue // 向前查找
			}

			// 打开文件
			done := false
			miss._grows_file(kit.Format(record[FILE]), func(data, head []string) bool {
				defer func() { offset++ }()
				if offset < begin {
					return false
				}
				if offset < end {
					item := Map{}
					for i := range head {
						if head[i] == EXTRA {
							item[head[i]] = kit.UnMarshal(data[i])
						} else {
							item[head[i]] = data[i]
						}
					}
					miss.Logger(OFFSET, offset, "type", item["type"], "name", item["name"], "text", item["text"])
					if done = cb(item); done {
						return true
					}
					return false
				}
				done = true
				return true
			})
			if done {
				return true
			}
		}
		break
	}
	return false
}
func (miss *Miss) _grows_file(p string, cb func(data []string, head []string) bool) {
	miss.Logger("open file", p)
	if f, e := miss.file.OpenFile(p); e == nil {
		defer f.Close()

		r := csv.NewReader(f)
		head, _ := r.Read()
		miss.Logger("read head", head)

		for {
			data, e := r.Read()
			if e != nil {
				miss.Logger("read data", e)
				break
			}
			if cb(data, head) {
				break
			}
		}
	}
}

func _grow_match(item Map, field, value string, index int, cb Any) (Map, int, bool) {
	if field == "" || strings.Contains(kit.Format(kit.Value(item, field)), value) {
		switch cb := cb.(type) {
		case func(Map):
			cb(item)
		case func(int, Map):
			cb(index, item)
		case func(int, Map) bool:
			if cb(index, item) {
				return item, index, true
			}
		}
		index++
		return item, index, false
	}
	return nil, index, false
}
