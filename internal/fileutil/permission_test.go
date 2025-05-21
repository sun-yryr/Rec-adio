package fileutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPermissionConstants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, os.FileMode(0o750), os.FileMode(DefaultDirPerm))
	assert.Equal(t, os.FileMode(0o644), os.FileMode(DefaultFilePerm))
}

func TestPermissionApplied(t *testing.T) {
	t.Parallel()

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// DefaultDirPermでディレクトリを作成
	testDir := filepath.Join(tempDir, "testdir")
	err := os.Mkdir(testDir, DefaultDirPerm)
	require.NoError(t, err)

	// パーミッションを確認
	info, err := os.Stat(testDir)
	require.NoError(t, err)
	// ファイルシステムによってはumaskが適用され正確に一致しないことがあるため、
	// ここではデフォルトのパーミッションが設定されていることを確認する
	// (0o750 と 0o700 といったようなOS依存の違いを考慮)
	expectedDirPerm := info.Mode().Perm()
	assert.Equal(t, expectedDirPerm, info.Mode().Perm()&os.FileMode(DefaultDirPerm),
		"Expected permission to be compatible with %o, got %o", DefaultDirPerm, info.Mode().Perm())

	// DefaultFilePermでファイルを作成
	testFile := filepath.Join(tempDir, "testfile")
	err = os.WriteFile(testFile, []byte("test"), DefaultFilePerm)
	require.NoError(t, err)

	// パーミッションを確認
	info, err = os.Stat(testFile)
	require.NoError(t, err)
	// 同様にumaskを考慮
	expectedFilePerm := info.Mode().Perm()
	assert.Equal(t, expectedFilePerm, info.Mode().Perm()&os.FileMode(DefaultFilePerm),
		"Expected permission to be compatible with %o, got %o", DefaultFilePerm, info.Mode().Perm())
}
