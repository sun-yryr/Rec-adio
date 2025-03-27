package recording

import (
        "context"
        "os"
        "path/filepath"
        "testing"
        "time"

        "github.com/sun-yryr/Rec-adio/internal/db/models"
        "github.com/sun-yryr/Rec-adio/internal/event"
        "gorm.io/driver/sqlite"
        "gorm.io/gorm"
)

// テスト用のデータベース接続を作成
func setupTestDB(t *testing.T) *gorm.DB {
        // 一時ファイルの作成
        tempFile := filepath.Join(os.TempDir(), "recadio_test.db")
        
        // 既存のファイルがあれば削除
        os.Remove(tempFile)
        
        // データベース接続
        db, err := gorm.Open(sqlite.Open(tempFile), &gorm.Config{})
        if err != nil {
                t.Fatalf("Failed to connect to database: %v", err)
        }
        
        // マイグレーション
        err = db.AutoMigrate(
                &models.Schedule{},
                &models.Record{},
                &models.Performer{},
                &models.ProgramInfo{},
                &models.PlatformAccount{},
                &models.SchedulePerformer{},
        )
        if err != nil {
                t.Fatalf("Failed to migrate database: %v", err)
        }
        
        return db
}

// テスト用のイベントリスナー
type TestEventListener struct {
        events []event.Event
}

func (l *TestEventListener) OnEvent(e event.Event) error {
        l.events = append(l.events, e)
        return nil
}

func (l *TestEventListener) GetSupportedEvents() []string {
        return []string{"*"} // すべてのイベントをサポート
}

func TestRecordingManager(t *testing.T) {
        // テスト用のデータベース接続
        db := setupTestDB(t)
        
        // テスト用のイベントマネージャー
        eventManager := event.NewEventManager()
        testListener := &TestEventListener{events: make([]event.Event, 0)}
        eventManager.AddListener(testListener)
        
        // 録音マネージャーの作成
        manager := NewRecordingManager(db, eventManager)
        
        // テスト用の出力ディレクトリ
        outputDir := filepath.Join(os.TempDir(), "recadio_test_output")
        os.MkdirAll(outputDir, 0755)
        defer os.RemoveAll(outputDir)
        
        // テスト用のレコーダープラグインの登録
        recorder := NewAGQRRecorder(outputDir, "https://example.com/stream")
        err := manager.RegisterPlugin(recorder)
        if err != nil {
                t.Fatalf("Failed to register plugin: %v", err)
        }
        
        // テスト用のスケジュールの作成
        schedule := models.NewSchedule(
                "テスト番組",
                time.Now(),
                1, // 1秒間の録音（テスト用に短く）
                "agqr",
        )
        err = db.Create(schedule).Error
        if err != nil {
                t.Fatalf("Failed to create schedule: %v", err)
        }
        
        // 録音の開始
        ctx := context.Background()
        recordID, err := manager.StartRecording(ctx, schedule.ID)
        if err != nil {
                t.Fatalf("Failed to start recording: %v", err)
        }
        
        // 録音IDが返されることを確認
        if recordID == "" {
                t.Fatal("Record ID is empty")
        }
        
        // 録音開始イベントが発行されることを確認
        if len(testListener.events) == 0 {
                t.Fatal("No events emitted")
        }
        
        // 最初のイベントが録音開始イベントであることを確認
        if testListener.events[0].Type != "recording.started" {
                t.Fatalf("Expected recording.started event, got %s", testListener.events[0].Type)
        }
        
        // 録音の状態を確認
        status, err := manager.GetRecordStatus(recordID)
        if err != nil {
                t.Fatalf("Failed to get record status: %v", err)
        }
        
        // 録音中または録音待ちの状態であることを確認
        if status != models.RecordStatusRecording && status != models.RecordStatusPending {
                t.Fatalf("Expected recording or pending status, got %s", status)
        }
        
        // 録音を停止
        err = manager.StopRecording(ctx, recordID)
        if err != nil {
                t.Logf("Warning: Failed to stop recording: %v", err)
        }
        
        // 録音レコードの取得
        var record models.Record
        err = db.First(&record, "id = ?", recordID).Error
        if err != nil {
                t.Fatalf("Failed to get record: %v", err)
        }
        
        // 録音マネージャーのシャットダウン
        err = manager.Shutdown(ctx)
        if err != nil {
                t.Logf("Warning: Failed to shutdown manager: %v", err)
        }
        
        t.Logf("Record status: %s", record.Status)
        t.Logf("Events received: %d", len(testListener.events))
        for i, e := range testListener.events {
                t.Logf("Event %d: %s", i, e.Type)
        }
}