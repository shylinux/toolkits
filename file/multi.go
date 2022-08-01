package file

import (
	"io"
	"os"
)

type MultiFile struct {
	list []File
}

func NewMultiFile(list ...File) File {
	return &MultiFile{list: list}
}
func (s *MultiFile) StatFile(p string) (os.FileInfo, error) {
	for _, f := range s.list {
		if s, e := f.StatFile(p); e == nil {
			return s, e
		}
	}
	return nil, os.ErrExist
}
func (s *MultiFile) OpenFile(p string) (io.ReadCloser, error) {
	for _, f := range s.list {
		if s, e := f.OpenFile(p); e == nil {
			return s, e
		}
	}
	return nil, os.ErrExist
}
func (s *MultiFile) CreateFile(p string) (io.WriteCloser, string, error) {
	for _, f := range s.list {
		if f, p, e := f.CreateFile(p); e == nil {
			return f, p, e
		}
	}
	return nil, p, os.ErrInvalid
}
func (s *MultiFile) AppendFile(p string) (io.ReadWriteCloser, error) {
	for _, f := range s.list {
		if _, e := f.StatFile(p); e == nil {
			if f, e := f.AppendFile(p); e == nil {
				return f, e
			}
		}
	}
	for _, f := range s.list {
		if f, e := f.AppendFile(p); e == nil {
			return f, e
		}
	}
	return nil, os.ErrInvalid
}
func (s *MultiFile) WriteFile(p string, b []byte) error {
	for _, f := range s.list {
		if e := f.WriteFile(p, b); e == nil {
			return e
		}
	}
	return nil
}
func (s *MultiFile) ReadDir(p string) (list []os.FileInfo, err error) {
	has := map[string]bool{}
	for _, f := range s.list {
		if ls, e := f.ReadDir(p); e == nil {
			for _, f := range ls {
				if has[f.Name()] {
					continue
				}
				has[f.Name()] = true
				list = append(list, f)
			}
		} else {
			err = e
		}
	}
	return list, nil
}
func (s *MultiFile) MkdirAll(p string, m os.FileMode) error {
	for _, f := range s.list {
		if e := f.MkdirAll(p, m); e == nil {
			return e
		}
	}
	return os.ErrInvalid
}
func (s *MultiFile) RemoveAll(p string) error {
	for _, f := range s.list {
		f.RemoveAll(p)
	}
	return nil
}
func (s *MultiFile) Remove(p string) error {
	for _, f := range s.list {
		f.Remove(p)
	}
	return nil
}
func (s *MultiFile) Rename(oldname, newname string) error {
	for _, f := range s.list {
		f.Rename(oldname, newname)
	}
	return nil
}
func (s *MultiFile) Symlink(oldname, newname string) error {
	for _, f := range s.list {
		f.Symlink(oldname, newname)
	}
	return nil
}
func (s *MultiFile) Link(oldname, newname string) error {
	for _, f := range s.list {
		f.Link(oldname, newname)
	}
	return nil
}
func (s *MultiFile) Close() {}
