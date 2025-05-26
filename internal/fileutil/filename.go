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
	invalidCharsPattern = regexp.MustCompile(`[\\/:*?"<>|\s]`)
	// 連続するアンダースコアを1つにまとめる。
	multipleUnderscoresPattern = regexp.MustCompile(`_+`)
	// Windows予約語のリスト
	windowsReservedNames = []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}
	// 特殊なディレクトリ名
	specialNames = []string{".", ".."}
)

// SanitizeFilename は、ファイル名に使用できない文字を安全な文字に置換する。
// Windows、macOS、Linuxで共通して使用できない文字、Windows予約語、特殊名を処理する。
func SanitizeFilename(filename string) string {
	// 不正な文字（空白含む）をアンダースコアに置換
	sanitized := invalidCharsPattern.ReplaceAllString(filename, "_")
	// 連続するアンダースコアを1つにまとめる
	sanitized = multipleUnderscoresPattern.ReplaceAllString(sanitized, "_")
	sanitized = strings.TrimSpace(sanitized)

	// 空文字列になった場合はデフォルト名を設定
	if sanitized == "" {
		sanitized = "recording"
	}

	// Windows予約語と特殊名のチェック
	allReservedNames := append(windowsReservedNames, specialNames...)
	for _, reserved := range allReservedNames {
		if strings.EqualFold(sanitized, reserved) {
			sanitized = "recording"
			break
		}
	}

	return sanitized
}
