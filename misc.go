package kit

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

func Hash(arg ...interface{}) (string, []string) {
	if len(arg) == 0 {
		arg = append(arg, "uniq")
	}
	args := []string{}
	for _, v := range Simple(arg...) {
		switch v {
		case "time":
			args = append(args, Format(time.Now()))
		case "rand":
			args = append(args, Format(rand.Int()))
		case "uniq":
			args = append(args, Format(time.Now()))
			args = append(args, Format(rand.Int()))
		default:
			args = append(args, v)
		}
	}

	h := md5.Sum([]byte(strings.Join(args, "")))
	return hex.EncodeToString(h[:]), args
}
func Hashs(arg ...interface{}) string {
	h, _ := Hash(arg...)
	return h
}
func ShortKey(list map[string]interface{}, min int, arg ...interface{}) string {
	h := Hashs(arg...)
	for i := min; i < len(h); i++ {
		if _, ok := list[h[:i]]; !ok {
			return h[:i]
		}
	}
	return h
}

func Revert(str []string) []string {
	for i := 0; i < len(str)/2; i++ {
		str[i], str[len(str)-1-i] = str[len(str)-1-i], str[i]
	}
	return str
}
func IndexOf(str []string, sub string) int {
	for i, v := range str {
		if v == sub {
			return i
		}
	}
	return -1
}

func Create(p string) (*os.File, string, error) {
	if dir, _ := path.Split(p); dir != "" {
		if e := os.MkdirAll(dir, 0777); e != nil {
			return nil, p, e
		}
	}
	f, e := os.Create(p)
	return f, p, e
}
func Duration(str interface{}) time.Duration {
	switch str := str.(type) {
	case string:
		d, _ := time.ParseDuration(str)
		return d
	}
	return time.Millisecond
}
