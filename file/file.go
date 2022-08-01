package file

import (
	"io"
	"os"
)

const (
	PS = "/"
	PT = "."

	FILE_MODE = 0640
	PATH_MODE = 0750

	NULL   = "/dev/null"
	STDIN  = "/dev/stdin"
	STDOUT = "/dev/stdout"
	STDERR = "/dev/stderr"
)

const FILE = "file"

type File interface {
	StatFile(string) (os.FileInfo, error)
	OpenFile(string) (io.ReadCloser, error)
	CreateFile(string) (io.WriteCloser, string, error)
	AppendFile(string) (io.ReadWriteCloser, error)
	WriteFile(string, []byte) error

	ReadDir(string) ([]os.FileInfo, error)
	MkdirAll(string, os.FileMode) error
	RemoveAll(string) error
	Remove(string) error
	Rename(string, string) error
	Symlink(string, string) error
	Link(string, string) error
	Close()
}
type FileInfo interface {
	os.FileInfo
	io.Reader
	io.Writer
	io.Closer
}
