package file

import (
	"io"
	"os"
	"time"
)

type PackFile struct {
	*VoidFile
}

func NewPackFile() File {
	return &PackFile{VoidFile: NewVoidFile().(*VoidFile)}
}
func (s *PackFile) OpenFile(p string) (io.ReadCloser, error) {
	defer s.lock.RLock()()
	if f, ok := s.get(p); ok {
		if f, ok := f.(*packFileInfo); ok {
			return NewPackFileInfo(p, f.b), nil
		}
	}
	return nil, os.ErrNotExist
}
func (s *PackFile) CreateFile(p string) (io.WriteCloser, string, error) {
	if f, p, e := createFile(s, p); f != nil {
		return f, p, e
	}
	defer s.lock.Lock()()
	return s.put(p, NewPackFileInfo(p, nil)), p, nil
}
func (s *PackFile) AppendFile(p string) (io.ReadWriteCloser, error) {
	defer s.lock.Lock()()
	if f, ok := s.get(p); ok {
		if f, ok := f.(*packFileInfo); ok {
			return NewPackFileInfo(p, f.b), nil
		}
	}
	return s.put(p, NewPackFileInfo(p, nil)), nil
}
func (s *PackFile) WriteFile(p string, b []byte) error {
	w, p, e := s.CreateFile(p)
	if e != nil {
		return e
	}
	_, e = w.Write(b)
	return e
}

type packFileInfo struct {
	*voidFileInfo
	b []byte
	i int
}

func NewPackFileInfo(p string, b []byte) *packFileInfo {
	return &packFileInfo{voidFileInfo: NewVoidFileInfo(p, len(b)), b: b}
}
func NewPackPathInfo(p string) *packFileInfo {
	return &packFileInfo{voidFileInfo: NewVoidPathInfo(p)}
}
func (s *packFileInfo) Read(p []byte) (n int, e error) {
	if s.i >= len(s.b) {
		return 0, io.EOF
	}
	if len(s.b) == 0 {
		return 0, io.EOF
	}
	n = copy(p, s.b[s.i:])
	s.i += n
	return n, nil
}
func (s *packFileInfo) Write(p []byte) (n int, e error) {
	s.t = time.Now()
	s.b = append(s.b, p...)
	s.n += len(p)
	return len(p), nil
}
