package fileutil

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/spf13/afero"
)

// FileSystem は基本的なファイル操作を抽象化します。
type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Remove(name string) error
	Rename(oldpath, newpath string) error
	UserHomeDir() (string, error)
}

// AferoFS は基盤となるafero.Fsを使用してFileSystemを実装します。
type AferoFS struct {
	base afero.Fs
}

// NewOsFS はOSファイルシステムで実装されたAferoFSを返します。
func NewOsFS() *AferoFS {
	return &AferoFS{base: afero.NewOsFs()}
}

// NewMemFS はインメモリファイルシステムで実装されたAferoFSを返します。
func NewMemFS() *AferoFS {
	return &AferoFS{base: afero.NewMemMapFs()}
}

// Stat は指定された名前のファイル情報を返します。
func (fs *AferoFS) Stat(name string) (os.FileInfo, error) {
	info, err := fs.base.Stat(name)
	if err != nil {
		return nil, errors.Wrap(err, "stat failed")
	}

	return info, nil
}

// ReadFile はファイル名で指定されたファイルを読み取り、その内容を返します。
func (fs *AferoFS) ReadFile(name string) ([]byte, error) {
	data, err := afero.ReadFile(fs.base, name)
	if err != nil {
		return nil, errors.Wrap(err, "read file failed")
	}

	return data, nil
}

// WriteFile は指定されたファイルにデータを書き込み、必要に応じてファイルを作成します。
func (fs *AferoFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	if err := afero.WriteFile(fs.base, name, data, perm); err != nil {
		return errors.Wrap(err, "write file failed")
	}

	return nil
}

// MkdirAll は指定されたパスのディレクトリを、必要な親ディレクトリと共に作成します。
func (fs *AferoFS) MkdirAll(path string, perm os.FileMode) error {
	if err := fs.base.MkdirAll(path, perm); err != nil {
		return errors.Wrap(err, "mkdir failed")
	}

	return nil
}

// Remove は指定されたファイルまたはディレクトリを削除します。
func (fs *AferoFS) Remove(name string) error {
	if err := fs.base.Remove(name); err != nil {
		return errors.Wrap(err, "remove file failed")
	}

	return nil
}

// Rename はoldpathをnewpathに名前変更（移動）します。
func (fs *AferoFS) Rename(oldpath, newpath string) error {
	if err := fs.base.Rename(oldpath, newpath); err != nil {
		return errors.Wrap(err, "rename file failed")
	}

	return nil
}

// UserHomeDir は現在のユーザーのホームディレクトリを返します。
func (fs *AferoFS) UserHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "user home dir")
	}

	return home, nil
}
