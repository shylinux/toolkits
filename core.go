package kit

// Split 字符串切分
func Split(str string) (res []string) {
	list := []rune(str)
	left, space, begin := '\000', true, 0
	for i := 0; i < len(list); i++ {
		switch list[i] {
		case '\'', '"', '`':
			if left == '\000' {
				left, space, begin = list[i], false, i+1
			} else if left == list[i] {
				res = append(res, string(list[begin:i]))
				left, space, begin = '\000', true, i+1
			}
		case '\t', ' ', '\n':
			if left != '\000' {
				break
			}
			if !space {
				res = append(res, string(list[begin:i]))
			}
			space, begin = true, i+1
		case '\\':
			space = false
		default:
			space = false
		}
	}

	if begin < len(list) {
		res = append(res, string(list[begin:]))
	}
	return res
}

// Value 字符串解析
func Value(data interface{}, key string, val ...string) interface{} {
	root := &data
	list := []*interface{}{root}
	for _, v := range val {

		switch last := data.(type) {
		case map[string]interface{}:
			switch v {
			case "{":
				node := map[string]interface{}{}
				last[key], data, list = node, node, append(list, &data)
			case "[":
				node := []interface{}{}
				last[key], data, list = node, node, append(list, &data)
			case "]":
				data, list = *list[len(list)-1], list[:len(list)-1]
			case "}":
				data, list = *list[len(list)-1], list[:len(list)-1]
			default:
				if key == "" {
					key = v
				} else {
					last[key], key = v, ""
				}
			}
		case []interface{}:
			switch v {
			case "{":
			case "[":
			case "]":
			case "}":
			default:
				last = append(last, v)
				data = last
			}
		case nil:
			switch v {
			case "{":
				node := map[string]interface{}{}
				data, *root = node, node
				list[0] = root
			case "[":
				node := []interface{}{}
				data, *root = node, node
				list[0] = root
			default:
				return v
			}
		}
	}
	return *root
}

// 数据读写
func Chain(data interface{}, key interface{}, val interface{}) interface{} {
	return nil
}
