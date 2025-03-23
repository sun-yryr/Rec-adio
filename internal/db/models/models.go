package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateID は、ランダムなIDを生成します
func GenerateID(prefix string) string {
	// 16バイト（32文字のHEX）のランダムな値を生成
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// エラーが発生した場合は、現在時刻をベースにしたIDを生成
		return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	}
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(b))
}

// GeneratePerformerID は、演者用のIDを生成します
func GeneratePerformerID() string {
	return GenerateID("perf")
}

// GeneratePlatformAccountID は、プラットフォームアカウント用のIDを生成します
func GeneratePlatformAccountID() string {
	return GenerateID("acct")
}

// GenerateProgramInfoID は、番組情報用のIDを生成します
func GenerateProgramInfoID() string {
	return GenerateID("prog")
}

// GenerateScheduleID は、スケジュール用のIDを生成します
func GenerateScheduleID() string {
	return GenerateID("sched")
}

// GenerateRecordID は、録音記録用のIDを生成します
func GenerateRecordID() string {
	return GenerateID("rec")
}

// GeneratePerformerSubscriptionID は、演者購読用のIDを生成します
func GeneratePerformerSubscriptionID() string {
	return GenerateID("psub")
}

// GeneratePlatformAccountSubscriptionID は、プラットフォームアカウント購読用のIDを生成します
func GeneratePlatformAccountSubscriptionID() string {
	return GenerateID("asub")
}

// TimeToString は、時間を文字列に変換します
func TimeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}

// StringToTime は、文字列を時間に変換します
func StringToTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// FormatDuration は、秒単位の期間を「HH:MM:SS」形式の文字列に変換します
func FormatDuration(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

// ParseDuration は、「HH:MM:SS」形式の文字列を秒単位の期間に変換します
func ParseDuration(s string) (int, error) {
	var hours, minutes, seconds int
	_, err := fmt.Sscanf(s, "%d:%d:%d", &hours, &minutes, &seconds)
	if err != nil {
		return 0, err
	}
	return hours*3600 + minutes*60 + seconds, nil
}
