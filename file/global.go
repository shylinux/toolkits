package file

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	kit "shylinux.com/x/toolkits"
)

var file = NewDiskFile()

func Init(f File) { file = f }

func StatFile(p string) (os.FileInfo, error) {
	return file.StatFile(p)
}
func OpenFile(p string) (io.ReadCloser, error) {
	return file.OpenFile(p)
}
func CreateFile(p string) (io.WriteCloser, string, error) {
	return file.CreateFile(p)
}
func AppendFile(p string) (io.ReadWriteCloser, error) {
	return file.AppendFile(p)
}
func WriteFile(p string, b []byte) error {
	return file.WriteFile(p, b)
}

func ReadDir(p string) (list []os.FileInfo, err error) {
	return file.ReadDir(p)
}
func MkdirAll(p string, m os.FileMode) error {
	return file.MkdirAll(p, m)
}
func RemoveAll(p string) error {
	return file.RemoveAll(p)
}
func Remove(p string) error {
	return file.Remove(p)
}
func Rename(oldname, newname string) error {
	return file.Rename(oldname, newname)
}
func Symlink(oldname, newname string) error {
	return file.Symlink(oldname, newname)
}
func Link(oldname, newname string) error {
	return file.Link(oldname, newname)
}
func Close() { file.Close() }

func ExistsFile(p string) bool {
	if s, e := file.StatFile(p); s != nil && e == nil {
		return true
	}
	return false
}
func ReadFile(p string) ([]byte, error) {
	if f, e := file.OpenFile(p); e == nil {
		return ioutil.ReadAll(f)
	} else {
		return nil, e
	}
}
func createFile(s File, p string) (io.ReadWriteCloser, string, error) {
	switch p {
	case NULL, path.Base(NULL), "":
		return NewVoidFileInfo(p, FILE_MODE), NULL, nil
	case STDOUT, path.Base(STDOUT):
		return os.Stdout, STDOUT, nil
	case STDERR, path.Base(STDERR), "stderr.log", "stderr.err.log":
		return os.Stderr, STDERR, nil
	}
	if dir, _ := path.Split(p); strings.Contains(p, PS) && dir != "" {
		s.MkdirAll(dir, PATH_MODE)
	}
	return nil, "", nil
}

type writeCloser struct{ w, c kit.Any }

func NewWriteCloser(w kit.Any, c kit.Any) io.WriteCloser {
	return &writeCloser{w, c}
}
func (w *writeCloser) Write(buf []byte) (int, error) {
	switch cb := w.w.(type) {
	case func([]byte) (int, error):
		return cb(buf)
	case func([]byte):
		cb(buf)
		return len(buf), nil
	default:
		return len(buf), nil
	}
}
func (w *writeCloser) Close() error {
	switch cb := w.c.(type) {
	case func() error:
		return cb()
	case func():
		cb()
		return nil
	default:
		return nil
	}
}

type readCloser struct {
	r io.Reader
}

func NewReadCloser(r io.Reader) io.ReadCloser {
	return &readCloser{r: r}
}
func (r *readCloser) Read(buf []byte) (int, error) {
	return r.r.Read(buf)
}
func (r *readCloser) Close() error {
	if c, ok := r.r.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
