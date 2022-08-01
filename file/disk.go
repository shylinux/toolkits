package file

import (
	"io"
	"io/ioutil"
	"os"
)

type DiskFile struct{}

func NewDiskFile() File {
	return &DiskFile{}
}
func (s *DiskFile) StatFile(p string) (os.FileInfo, error) {
	return os.Stat(p)
}
func (s *DiskFile) OpenFile(p string) (io.ReadCloser, error) {
	return os.Open(p)
}
func (s *DiskFile) CreateFile(p string) (io.WriteCloser, string, error) {
	if f, p, e := createFile(s, p); f != nil {
		return f, p, e
	}
	f, e := os.Create(p)
	return f, p, e
}
func (s *DiskFile) AppendFile(p string) (io.ReadWriteCloser, error) {
	if f, _, e := createFile(s, p); f != nil {
		return f, e
	}
	return os.OpenFile(p, os.O_RDWR|os.O_APPEND|os.O_CREATE, FILE_MODE)
}
func (s *DiskFile) WriteFile(p string, b []byte) error {
	return ioutil.WriteFile(p, b, FILE_MODE)
}
func (s *DiskFile) ReadDir(p string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(p)
}
func (s *DiskFile) MkdirAll(p string, m os.FileMode) error {
	return os.MkdirAll(p, m)
}
func (s *DiskFile) RemoveAll(p string) error {
	return os.RemoveAll(p)
}
func (s *DiskFile) Remove(p string) error {
	return os.Remove(p)
}
func (s *DiskFile) Rename(oldname, newname string) error {
	return os.Rename(oldname, newname)
}
func (s *DiskFile) Symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}
func (s *DiskFile) Link(oldname, newname string) error {
	return os.Link(oldname, newname)
}
func (s *DiskFile) Close() {}
