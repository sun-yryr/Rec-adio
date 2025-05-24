package fileutil

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/spf13/afero"
)

// FileSystem abstracts basic file operations.
type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Remove(name string) error
	Rename(oldpath, newpath string) error
	UserHomeDir() (string, error)
}

// aferoFS implements FileSystem using an underlying afero.Fs.
type aferoFS struct {
	base afero.Fs
}

// NewOsFS returns a FileSystem implemented with the os filesystem.
//
//nolint:ireturn
func NewOsFS() FileSystem {
	return &aferoFS{base: afero.NewOsFs()}
}

// NewMemFS returns a FileSystem implemented with an in-memory filesystem.
//
//nolint:ireturn
func NewMemFS() FileSystem {
	return &aferoFS{base: afero.NewMemMapFs()}
}

func (fs *aferoFS) Stat(name string) (os.FileInfo, error) {
	info, err := fs.base.Stat(name)
	if err != nil {
		return nil, errors.Wrap(err, "stat failed")
	}

	return info, nil
}

func (fs *aferoFS) ReadFile(name string) ([]byte, error) {
	data, err := afero.ReadFile(fs.base, name)
	if err != nil {
		return nil, errors.Wrap(err, "read file failed")
	}

	return data, nil
}

func (fs *aferoFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	if err := afero.WriteFile(fs.base, name, data, perm); err != nil {
		return errors.Wrap(err, "write file failed")
	}

	return nil
}

func (fs *aferoFS) MkdirAll(path string, perm os.FileMode) error {
	if err := fs.base.MkdirAll(path, perm); err != nil {
		return errors.Wrap(err, "mkdir failed")
	}

	return nil
}

func (fs *aferoFS) Remove(name string) error {
	if err := fs.base.Remove(name); err != nil {
		return errors.Wrap(err, "remove file failed")
	}

	return nil
}

func (fs *aferoFS) Rename(oldpath, newpath string) error {
	if err := fs.base.Rename(oldpath, newpath); err != nil {
		return errors.Wrap(err, "rename file failed")
	}

	return nil
}

func (fs *aferoFS) UserHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "user home dir")
	}

	return home, nil
}
