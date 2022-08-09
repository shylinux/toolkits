package file

import (
	"io"
	"os"
	"path"
	"strings"
	"time"

	kit "shylinux.com/x/toolkits"
)

type VoidFile struct {
	list map[string]FileInfo
	lock Lock
}

func NewVoidFile() File {
	return &VoidFile{list: map[string]FileInfo{}}
}
func (s *VoidFile) get(p string) (FileInfo, bool) {
	f, ok := s.list[p]
	return f, ok
}
func (s *VoidFile) put(p string, f FileInfo) FileInfo {
	if _, ok := s.list[p]; !ok {
		if b, ok := s.list[path.Dir(p)]; ok {
			switch b := b.(type) {
			case *voidFileInfo:
				b.n++
			case *packFileInfo:
				b.n++
			}
		}
	}
	s.list[p] = f
	return f
}
func (s *VoidFile) StatFile(p string) (os.FileInfo, error) {
	defer s.lock.RLock()()
	if f, ok := s.get(p); ok {
		return f, nil
	}
	return nil, os.ErrNotExist
}
func (s *VoidFile) OpenFile(p string) (io.ReadCloser, error) {
	defer s.lock.RLock()()
	if f, ok := s.get(p); ok {
		return f, nil
	}
	return nil, os.ErrNotExist
}
func (s *VoidFile) CreateFile(p string) (io.WriteCloser, string, error) {
	if f, p, e := createFile(s, p); f != nil {
		return f, p, e
	}
	defer s.lock.Lock()()
	return s.put(p, NewVoidFileInfo(p, 0)), p, nil
}
func (s *VoidFile) AppendFile(p string) (io.ReadWriteCloser, error) {
	defer s.lock.Lock()()
	if f, ok := s.get(p); ok {
		return f, nil
	}
	return s.put(p, NewVoidFileInfo(p, 0)), nil
}
func (s *VoidFile) WriteFile(p string, b []byte) error {
	w, p, e := s.CreateFile(p)
	if e != nil {
		return e
	}
	defer w.Close()
	_, e = w.Write(b)
	return e
}
func (s *VoidFile) ReadDir(p string) ([]os.FileInfo, error) {
	defer s.lock.RLock()()
	list := []os.FileInfo{}
	for k, s := range s.list {
		if strings.HasPrefix(k, p+kit.Select("", "/", !strings.HasSuffix(p, "/"))) {
			if len(kit.Split(strings.TrimPrefix(k, p), PS)) == 1 {
				list = append(list, s)
			}
		}
	}
	return list, nil
}
func (s *VoidFile) MkdirAll(p string, m os.FileMode) error {
	defer s.lock.Lock()()
	ls := kit.Split(p, PS)
	for i := 0; i < len(ls); i++ {
		p := kit.Join(ls[:i+1], PS)
		if _, ok := s.list[p]; !ok {
			s.put(p, NewVoidPathInfo(p))
		}
	}
	return nil
}
func (s *VoidFile) RemoveAll(p string) error {
	defer s.lock.Lock()()
	list := []string{}
	for k := range s.list {
		if strings.HasPrefix(k, p+PS) || k == p {
			list = append(list, k)
		}
	}
	for _, k := range list {
		delete(s.list, k)
	}
	return nil
}
func (s *VoidFile) Remove(p string) error {
	defer s.lock.Lock()()
	for k := range s.list {
		if strings.HasPrefix(k, p+PS) {
			return os.ErrInvalid
		}
	}
	delete(s.list, p)
	return nil
}
func (s *VoidFile) Rename(oldname, newname string) error {
	defer s.lock.Lock()()
	if _, ok := s.list[newname]; !ok {
		if f, ok := s.list[oldname]; ok {
			delete(s.list, oldname)
			s.list[newname] = f
			return nil
		}
		return os.ErrNotExist
	}
	return os.ErrExist
}
func (s *VoidFile) Symlink(oldname, newname string) error {
	defer s.lock.Lock()()
	if _, ok := s.list[newname]; !ok {
		if f, ok := s.list[oldname]; ok {
			s.list[newname] = f
			return nil
		}
		return os.ErrNotExist
	}
	return os.ErrExist
}
func (s *VoidFile) Link(oldname, newname string) error {
	defer s.lock.Lock()()
	if _, ok := s.list[newname]; !ok {
		if f, ok := s.list[oldname]; ok {
			s.list[newname] = f
			return nil
		}
		return os.ErrNotExist
	}
	return os.ErrExist
}
func (s *VoidFile) Close() {}

type voidFileInfo struct {
	p string
	m os.FileMode
	t time.Time
	d bool
	n int
}

func NewVoidFileInfo(p string, n int) *voidFileInfo {
	return &voidFileInfo{p: p, m: FILE_MODE, t: time.Now(), n: n}
}
func NewVoidPathInfo(p string) *voidFileInfo {
	return &voidFileInfo{p: p, m: PATH_MODE, t: time.Now(), d: true}
}
func (s *voidFileInfo) Name() string       { return path.Base(s.p) }
func (s *voidFileInfo) Size() int64        { return int64(s.n) }
func (s *voidFileInfo) Mode() os.FileMode  { return s.m }
func (s *voidFileInfo) ModTime() time.Time { return s.t }
func (s *voidFileInfo) IsDir() bool        { return s.d }
func (s *voidFileInfo) Sys() interface{}   { return nil }

func (s *voidFileInfo) Read(p []byte) (n int, err error) { return 0, io.EOF }
func (s *voidFileInfo) Write(p []byte) (n int, err error) {
	s.t = time.Now()
	s.n += len(p)
	return len(p), nil
}
func (s *voidFileInfo) Close() error { return nil }
