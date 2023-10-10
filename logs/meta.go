package logs

import (
	"fmt"
	"strings"

	kit "shylinux.com/x/toolkits"
)

const (
	PREFIX   = "prefix"
	SUFFIX   = "suffix"
	FILELINE = "fileline"
	TRACEID  = "id"
)

type Meta struct{ Key, Value string }

func (s Meta) String() string {
	switch s.Key {
	case PREFIX, "", FILELINE, SUFFIX:
		return " " + s.Value
	}
	return " " + kit.FormatShow(s.Key, s.Value)
}
func ValueMeta(p string) Meta   { return Meta{"", p} }
func PrefixMeta(p string) Meta  { return Meta{PREFIX, p} }
func SuffixMeta(p string) Meta  { return Meta{SUFFIX, p} }
func TraceidMeta(p string) Meta { return Meta{TRACEID, p + " "} }
func FileLineMeta(p Any) Meta {
	switch p := p.(type) {
	case int:
		if p > 0 {
			return Meta{FILELINE, FileLine(p + 1)}
		} else {
			return Meta{FILELINE, FileLine(p - 1)}
		}
	case string:
		return Meta{FILELINE, p}
	}
	return Meta{}
}

func Format(str string, arg ...Any) string {
	prefix, args, has, suffix := []Any{}, []Any{}, false, []Any{}
	traceid := []Any{}
	for _, v := range arg {
		switch v := v.(type) {
		case Meta:
			switch v.Key {
			case PREFIX:
				prefix = append(prefix, v)
			case TRACEID:
				if strings.TrimSpace(v.Value) != "" {
					traceid = append(traceid, v)
				}
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
	if len(traceid) > 1 {
		traceid = traceid[:1]
	}
	return fmt.Sprint(append(append(traceid, prefix...), suffix...)...)
}
