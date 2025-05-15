// Package fileutil は、ファイルに関するユーティリティを提供する。
package fileutil

import (
	"regexp"
	"strings"
)

// SanitizeFilename は、ファイル名に使用できない文字を安全な文字に置換する。
// Windows、macOS、Linuxで共通して使用できない文字を処理する。
func SanitizeFilename(filename string) string {
	// ファイル名に使用できない文字を正規表現で定義
	// Windows: \ / : * ? " < > |
	// macOS: / :
	// Linux: /
	pattern := regexp.MustCompile(`[\\/:*?"<>|]`)

	// 不正な文字をアンダースコアに置換
	sanitized := pattern.ReplaceAllString(filename, "_")

	// 連続するアンダースコアを1つにまとめる
	pattern = regexp.MustCompile(`_+`)
	sanitized = pattern.ReplaceAllString(sanitized, "_")

	sanitized = strings.TrimSpace(sanitized)

	// 空文字列になった場合はデフォルト名を設定
	if sanitized == "" {
		sanitized = "recording"
	}

	return sanitized
}
