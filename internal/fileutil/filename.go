// Package fileutil は、ファイルに関するユーティリティを提供する。
package fileutil

import (
	"regexp"
	"strings"
)

var (
	/* ファイル名に使用できない文字を正規表現で定義。
	* Windows: \ / : * ? " < > |
	* macOS: / :
	* Linux: /
	 */
	invalidCharsPattern = regexp.MustCompile(`[\\/:*?"<>|]`)
	// 連続するアンダースコアを1つにまとめる。
	multipleUnderscoresPattern = regexp.MustCompile(`_+`)
)

// SanitizeFilename は、ファイル名に使用できない文字を安全な文字に置換する。
// Windows、macOS、Linuxで共通して使用できない文字を処理する。
func SanitizeFilename(filename string) string {
	// 不正な文字をアンダースコアに置換
	sanitized := invalidCharsPattern.ReplaceAllString(filename, "_")
	// 連続するアンダースコアを1つにまとめる
	sanitized = multipleUnderscoresPattern.ReplaceAllString(sanitized, "_")
	sanitized = strings.TrimSpace(sanitized)

	// 空文字列になった場合はデフォルト名を設定
	if sanitized == "" {
		sanitized = "recording"
	}

	return sanitized
}
