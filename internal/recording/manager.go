package recording

import (
        "context"
        "errors"
        "fmt"
        "sync"
        "time"

        "github.com/sun-yryr/Rec-adio/internal/db"
        "github.com/sun-yryr/Rec-adio/internal/db/models"
        "github.com/sun-yryr/Rec-adio/internal/event"
        "github.com/sun-yryr/Rec-adio/internal/plugin"
        "gorm.io/gorm"
)

// RecordingManager は録音の管理を行うマネージャーです
type RecordingManager struct {
        db           *gorm.DB
        eventEmitter event.EventEmitter
        plugins      map[string]plugin.RecorderPlugin
        sessions     map[string]plugin.RecordingSession
        mu           sync.RWMutex
}

// NewRecordingManager は新しいRecordingManagerインスタンスを作成します
func NewRecordingManager(db *gorm.DB, eventEmitter event.EventEmitter) *RecordingManager {
        return &RecordingManager{
                db:           db,
                eventEmitter: eventEmitter,
                plugins:      make(map[string]plugin.RecorderPlugin),
                sessions:     make(map[string]plugin.RecordingSession),
        }
}

// RegisterPlugin はレコーダープラグインを登録します
func (m *RecordingManager) RegisterPlugin(p plugin.RecorderPlugin) error {
        m.mu.Lock()
        defer m.mu.Unlock()

        info := p.GetPlatformInfo()
        if _, exists := m.plugins[info.ID]; exists {
                return fmt.Errorf("plugin with ID %s already registered", info.ID)
        }

        m.plugins[info.ID] = p
        return nil
}

// GetPlugin はプラットフォームIDに対応するプラグインを取得します
func (m *RecordingManager) GetPlugin(platformID string) (plugin.RecorderPlugin, error) {
        m.mu.RLock()
        defer m.mu.RUnlock()

        p, exists := m.plugins[platformID]
        if !exists {
                return nil, fmt.Errorf("plugin for platform %s not found", platformID)
        }

        return p, nil
}

// StartRecording はスケジュールに基づいて録音を開始します
func (m *RecordingManager) StartRecording(ctx context.Context, scheduleID string) (string, error) {
        // スケジュールの取得
        var schedule models.Schedule
        if err := m.db.First(&schedule, "id = ?", scheduleID).Error; err != nil {
                return "", fmt.Errorf("failed to get schedule: %w", err)
        }

        // プラグインの取得
        p, err := m.GetPlugin(schedule.Platform)
        if err != nil {
                return "", err
        }

        // 録音レコードの作成
        record := models.NewRecord(
                schedule.ID,
                schedule.Title,
                time.Now(),
                schedule.Platform,
        ).SetStatus(models.RecordStatusRecording)

        if err := m.db.Create(record).Error; err != nil {
                return "", fmt.Errorf("failed to create record: %w", err)
        }

        // 録音の開始
        session, err := p.Record(ctx, &schedule)
        if err != nil {
                // 録音開始に失敗した場合はレコードを更新
                record.SetStatus(models.RecordStatusFailed).SetErrorMessage(err.Error())
                m.db.Save(record)
                return "", fmt.Errorf("failed to start recording: %w", err)
        }

        // セッションの保存
        m.mu.Lock()
        m.sessions[record.ID] = session
        m.mu.Unlock()

        // ステータス変更時のコールバックを設定
        session.OnStatusChange(func(status plugin.RecordingStatus) {
                m.handleStatusChange(record.ID, status)
        })

        // 録音開始イベントの発行
        m.emitRecordingEvent("recording.started", map[string]interface{}{
                "record_id":   record.ID,
                "schedule_id": schedule.ID,
                "title":       schedule.Title,
                "platform":    schedule.Platform,
        })

        return record.ID, nil
}

// StopRecording は録音を停止します
func (m *RecordingManager) StopRecording(ctx context.Context, recordID string) error {
        m.mu.RLock()
        session, exists := m.sessions[recordID]
        m.mu.RUnlock()

        if !exists {
                return fmt.Errorf("recording session for record %s not found", recordID)
        }

        // 録音の停止
        if err := session.Stop(ctx); err != nil {
                return fmt.Errorf("failed to stop recording: %w", err)
        }

        return nil
}

// GetRecordStatus は録音のステータスを取得します
func (m *RecordingManager) GetRecordStatus(recordID string) (models.RecordStatus, error) {
        var record models.Record
        if err := m.db.First(&record, "id = ?", recordID).Error; err != nil {
                return "", fmt.Errorf("failed to get record: %w", err)
        }

        return record.Status, nil
}

// ListRecords は録音レコードの一覧を取得します
func (m *RecordingManager) ListRecords(filter db.RecordFilter) ([]*models.Record, error) {
        var records []*models.Record
        query := m.db.Model(&models.Record{})

        if filter.ScheduleID != "" {
                query = query.Where("schedule_id = ?", filter.ScheduleID)
        }

        if filter.Platform != "" {
                query = query.Where("platform = ?", filter.Platform)
        }

        if filter.Status != "" {
                query = query.Where("status = ?", filter.Status)
        }

        if !filter.StartTimeFrom.IsZero() {
                query = query.Where("record_time >= ?", filter.StartTimeFrom)
        }

        if !filter.StartTimeTo.IsZero() {
                query = query.Where("record_time <= ?", filter.StartTimeTo)
        }

        if filter.Title != "" {
                query = query.Where("title LIKE ?", "%"+filter.Title+"%")
        }

        if err := query.Find(&records).Error; err != nil {
                return nil, fmt.Errorf("failed to list records: %w", err)
        }

        return records, nil
}

// handleStatusChange は録音ステータスの変更を処理します
func (m *RecordingManager) handleStatusChange(recordID string, status plugin.RecordingStatus) {
        // レコードの取得
        var record models.Record
        if err := m.db.First(&record, "id = ?", recordID).Error; err != nil {
                // エラーログを出力
                fmt.Printf("Failed to get record %s: %v\n", recordID, err)
                return
        }

        // セッションの取得
        m.mu.RLock()
        session, exists := m.sessions[recordID]
        m.mu.RUnlock()

        if !exists {
                fmt.Printf("Recording session for record %s not found\n", recordID)
                return
        }

        // ステータスの更新
        var recordStatus models.RecordStatus
        switch status {
        case plugin.RecordingStatusPending:
                recordStatus = models.RecordStatusPending
        case plugin.RecordingStatusRecording:
                recordStatus = models.RecordStatusRecording
        case plugin.RecordingStatusCompleted:
                recordStatus = models.RecordStatusCompleted
        case plugin.RecordingStatusFailed:
                recordStatus = models.RecordStatusFailed
        case plugin.RecordingStatusCanceled:
                recordStatus = models.RecordStatusCanceled
        default:
                fmt.Printf("Unknown recording status: %s\n", status)
                return
        }

        // レコードの更新
        record.SetStatus(recordStatus)

        // 録音が完了した場合は追加情報を設定
        if status == plugin.RecordingStatusCompleted {
                // 出力パスの設定
                record.SetFilePath(session.GetOutputPath())

                // メタデータの取得
                metadata := session.GetMetadata()

                // 期間の設定
                if duration, ok := metadata["duration"].(int); ok {
                        record.SetDuration(duration)
                }

                // ファイルサイズの設定
                if fileSize, ok := metadata["file_size"].(int64); ok {
                        record.SetFileSize(fileSize)
                }

                // セッションの削除
                m.mu.Lock()
                delete(m.sessions, recordID)
                m.mu.Unlock()
        }

        // 録音が失敗した場合はエラーメッセージを設定
        if status == plugin.RecordingStatusFailed {
                if metadata := session.GetMetadata(); metadata != nil {
                        if errorMsg, ok := metadata["error_message"].(string); ok {
                                record.SetErrorMessage(errorMsg)
                        }
                }

                // セッションの削除
                m.mu.Lock()
                delete(m.sessions, recordID)
                m.mu.Unlock()
        }

        // レコードの保存
        if err := m.db.Save(&record).Error; err != nil {
                fmt.Printf("Failed to update record %s: %v\n", recordID, err)
                return
        }

        // イベントの発行
        var eventType string
        switch status {
        case plugin.RecordingStatusCompleted:
                eventType = "recording.completed"
        case plugin.RecordingStatusFailed:
                eventType = "recording.failed"
        case plugin.RecordingStatusCanceled:
                eventType = "recording.canceled"
        default:
                // その他のステータスはイベントを発行しない
                return
        }

        // イベントデータの作成
        eventData := map[string]interface{}{
                "record_id":   record.ID,
                "schedule_id": record.ScheduleID,
                "title":       record.Title,
                "platform":    record.Platform,
                "status":      string(record.Status),
        }

        // エラーメッセージがある場合は追加
        if record.ErrorMessage != "" {
                eventData["error_message"] = record.ErrorMessage
        }

        // ファイルパスがある場合は追加
        if record.FilePath != "" {
                eventData["file_path"] = record.FilePath
        }

        // イベントの発行
        m.emitRecordingEvent(eventType, eventData)
}

// emitRecordingEvent は録音関連のイベントを発行します
func (m *RecordingManager) emitRecordingEvent(eventType string, data map[string]interface{}) {
        if m.eventEmitter == nil {
                return
        }

        event := event.Event{
                Type:      eventType,
                Timestamp: time.Now(),
                Data:      data,
        }

        if err := m.eventEmitter.EmitEvent(event); err != nil {
                fmt.Printf("Failed to emit event %s: %v\n", eventType, err)
        }
}

// Shutdown はマネージャーのシャットダウン処理を行います
func (m *RecordingManager) Shutdown(ctx context.Context) error {
        m.mu.Lock()
        defer m.mu.Unlock()

        // 実行中の録音セッションをすべて停止
        var errs []error
        for id, session := range m.sessions {
                if err := session.Stop(ctx); err != nil {
                        errs = append(errs, fmt.Errorf("failed to stop session %s: %w", id, err))
                }
        }

        // プラグインのシャットダウン
        for id, p := range m.plugins {
                if err := p.Shutdown(ctx); err != nil {
                        errs = append(errs, fmt.Errorf("failed to shutdown plugin %s: %w", id, err))
                }
        }

        // エラーがあれば結合して返す
        if len(errs) > 0 {
                var errMsg string
                for _, err := range errs {
                        errMsg += err.Error() + "; "
                }
                return errors.New(errMsg)
        }

        return nil
}