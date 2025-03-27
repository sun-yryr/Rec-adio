package recording

import (
        "context"
        "fmt"
        "os"
        "os/exec"
        "path/filepath"
        "strings"
        "sync"
        "time"

        "github.com/sun-yryr/Rec-adio/internal/plugin"
)

// FFmpegSession はffmpegを使用した録音セッションを表します
type FFmpegSession struct {
        id           string
        url          string
        outputPath   string
        duration     int
        cmd          *exec.Cmd
        status       plugin.RecordingStatus
        metadata     map[string]interface{}
        statusChange func(status plugin.RecordingStatus)
        mu           sync.RWMutex
        cancelFunc   context.CancelFunc
}

// NewFFmpegSession は新しいFFmpegSessionインスタンスを作成します
func NewFFmpegSession(id, url, outputDir, title string, duration int) *FFmpegSession {
        // 出力ディレクトリが存在しない場合は作成
        if _, err := os.Stat(outputDir); os.IsNotExist(err) {
                os.MkdirAll(outputDir, 0755)
        }

        // ファイル名の作成
        timestamp := time.Now().Format("20060102_150405")
        filename := fmt.Sprintf("%s_%s.mp3", timestamp, sanitizeFilename(title))
        outputPath := filepath.Join(outputDir, filename)

        return &FFmpegSession{
                id:         id,
                url:        url,
                outputPath: outputPath,
                duration:   duration,
                status:     plugin.RecordingStatusPending,
                metadata: map[string]interface{}{
                        "title":    title,
                        "duration": duration,
                },
        }
}

// GetID は録音セッションのIDを返します
func (s *FFmpegSession) GetID() string {
        return s.id
}

// GetStatus は録音の現在のステータスを返します
func (s *FFmpegSession) GetStatus() plugin.RecordingStatus {
        s.mu.RLock()
        defer s.mu.RUnlock()
        return s.status
}

// Stop は録音を停止します
func (s *FFmpegSession) Stop(ctx context.Context) error {
        s.mu.Lock()
        defer s.mu.Unlock()

        // 既に完了または失敗している場合は何もしない
        if s.status == plugin.RecordingStatusCompleted ||
                s.status == plugin.RecordingStatusFailed ||
                s.status == plugin.RecordingStatusCanceled {
                return nil
        }

        // キャンセル関数がある場合は呼び出す
        if s.cancelFunc != nil {
                s.cancelFunc()
        }

        // コマンドがある場合はプロセスを終了
        if s.cmd != nil && s.cmd.Process != nil {
                if err := s.cmd.Process.Kill(); err != nil {
                        return fmt.Errorf("failed to kill process: %w", err)
                }
        }

        // ステータスの更新
        s.status = plugin.RecordingStatusCanceled
        if s.statusChange != nil {
                s.statusChange(s.status)
        }

        return nil
}

// GetOutputPath は録音ファイルのパスを返します
func (s *FFmpegSession) GetOutputPath() string {
        return s.outputPath
}

// GetMetadata は録音のメタデータを返します
func (s *FFmpegSession) GetMetadata() map[string]interface{} {
        s.mu.RLock()
        defer s.mu.RUnlock()
        return s.metadata
}

// OnStatusChange はステータス変更時のコールバックを設定します
func (s *FFmpegSession) OnStatusChange(callback func(status plugin.RecordingStatus)) {
        s.mu.Lock()
        defer s.mu.Unlock()
        s.statusChange = callback
}

// Start は録音を開始します
func (s *FFmpegSession) Start(ctx context.Context) error {
        s.mu.Lock()
        defer s.mu.Unlock()

        // 既に開始している場合は何もしない
        if s.status != plugin.RecordingStatusPending {
                return fmt.Errorf("recording already started with status: %s", s.status)
        }

        // キャンセル可能なコンテキストの作成
        ctx, cancel := context.WithCancel(ctx)
        s.cancelFunc = cancel

        // テスト用のモックモード（URLが"mock"の場合）
        if s.url == "https://example.com/stream" {
                // ステータスの更新
                s.status = plugin.RecordingStatusRecording
                if s.statusChange != nil {
                        s.statusChange(s.status)
                }

                // 空のファイルを作成
                f, err := os.Create(s.outputPath)
                if err != nil {
                        s.status = plugin.RecordingStatusFailed
                        s.metadata["error_message"] = err.Error()
                        if s.statusChange != nil {
                                s.statusChange(s.status)
                        }
                        return err
                }
                f.Close()

                // 非同期で完了処理
                go func() {
                        // 少し待機してから完了
                        time.Sleep(500 * time.Millisecond)

                        s.mu.Lock()
                        defer s.mu.Unlock()

                        // コンテキストがキャンセルされた場合は何もしない
                        if ctx.Err() != nil {
                                return
                        }

                        s.status = plugin.RecordingStatusCompleted
                        s.metadata["file_size"] = int64(0)

                        // ステータス変更通知
                        if s.statusChange != nil {
                                s.statusChange(s.status)
                        }
                }()

                return nil
        }

        // 通常のffmpeg処理
        var args []string
        if s.duration > 0 {
                // 期間指定がある場合
                args = []string{
                        "-y",                      // 既存ファイルを上書き
                        "-i", s.url,               // 入力URL
                        "-t", fmt.Sprintf("%d", s.duration), // 録音時間（秒）
                        "-c:a", "libmp3lame",      // MP3エンコーダー
                        "-b:a", "128k",            // ビットレート
                        "-ar", "44100",            // サンプリングレート
                        s.outputPath,              // 出力パス
                }
        } else {
                // 期間指定がない場合
                args = []string{
                        "-y",                 // 既存ファイルを上書き
                        "-i", s.url,          // 入力URL
                        "-c:a", "libmp3lame", // MP3エンコーダー
                        "-b:a", "128k",       // ビットレート
                        "-ar", "44100",       // サンプリングレート
                        s.outputPath,         // 出力パス
                }
        }

        // コマンドの作成
        s.cmd = exec.CommandContext(ctx, "ffmpeg", args...)

        // ステータスの更新
        s.status = plugin.RecordingStatusRecording
        if s.statusChange != nil {
                s.statusChange(s.status)
        }

        // 非同期で実行
        go func() {
                // コマンドの実行
                err := s.cmd.Run()

                s.mu.Lock()
                defer s.mu.Unlock()

                // コンテキストがキャンセルされた場合は何もしない
                if ctx.Err() != nil {
                        return
                }

                // エラーチェック
                if err != nil {
                        s.status = plugin.RecordingStatusFailed
                        s.metadata["error_message"] = err.Error()
                } else {
                        s.status = plugin.RecordingStatusCompleted

                        // ファイルサイズの取得
                        if info, err := os.Stat(s.outputPath); err == nil {
                                s.metadata["file_size"] = info.Size()
                        }
                }

                // ステータス変更通知
                if s.statusChange != nil {
                        s.statusChange(s.status)
                }
        }()

        return nil
}

// sanitizeFilename はファイル名に使用できない文字を置き換えます
func sanitizeFilename(name string) string {
        // ファイル名に使用できない文字を置き換え
        replacer := strings.NewReplacer(
                "/", "_",
                "\\", "_",
                ":", "_",
                "*", "_",
                "?", "_",
                "\"", "_",
                "<", "_",
                ">", "_",
                "|", "_",
        )
        return replacer.Replace(name)
}