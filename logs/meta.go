package logs

import (
	"fmt"

	kit "shylinux.com/x/toolkits"
)

const (
	PREFIX   = "prefix"
	FILELINE = "fileline"
	SUFFIX   = "suffix"
)

type Meta struct{ Key, Value string }

func (s Meta) String() string {
	switch s.Key {
	case PREFIX, "", FILELINE, SUFFIX:
		return " " + s.Value
	}
	return " " + kit.FormatShow(s.Key, s.Value)
}
func ValueMeta(p string) Meta    { return Meta{"", p} }
func PrefixMeta(p string) Meta   { return Meta{PREFIX, p} }
func SuffixMeta(p string) Meta   { return Meta{SUFFIX, p} }
func FileLineMeta(p string) Meta { return Meta{FILELINE, p} }

func Format(str string, arg ...Any) string {
	prefix, args, has, suffix := []Any{}, []Any{}, false, []Any{}
	for _, v := range arg {
		switch v := v.(type) {
		case Meta:
			switch v.Key {
			case PREFIX:
				prefix = append(prefix, v)
			case FILELINE:
				if has {
					continue
				}
				has = true
				fallthrough
			case SUFFIX:
				fallthrough
			default:
				suffix = append(suffix, v)
			}
		default:
			args = append(args, v)
		}
	}
	if len(args) == 0 {
		prefix = append(prefix, fmt.Sprint(str))
	} else if str == "" {
		prefix = append(prefix, fmt.Sprint(args...))
	} else {
		prefix = append(prefix, fmt.Sprintf(str, args...))
	}
	return fmt.Sprint(append(prefix, suffix...)...)
}
