# API経由での即時録音開始の実装方法

録音チェックワーカーからレコーダープラグインへの録音開始通知と、API経由での即時録音開始を実現するには、以下のような設計が考えられます。

## 1. レコーダープラグインインターフェース

まず、すべてのレコーダープラグインが実装すべきインターフェースを定義します：

```go
// RecorderPlugin はレコーダープラグインのインターフェースです
type RecorderPlugin interface {
    // GetPlatformInfo はサポートするプラットフォームの情報を返します
    GetPlatformInfo() PlatformInfo
    
    // Record は指定されたスケジュールに基づいて録音を開始します
    Record(ctx context.Context, schedule Schedule) (RecordingSession, error)
    
    // StopRecording は録音を停止します
    StopRecording(ctx context.Context, recordID string) error
    
    // GetSchedule はプラットフォームから番組スケジュールを取得します
    GetSchedule(ctx context.Context, filter ScheduleFilter) ([]Schedule, error)
    
    // Initialize はプラグインの初期化を行います
    Initialize(ctx context.Context, config Config) error
    
    // Shutdown はプラグインのシャットダウン処理を行います
    Shutdown(ctx context.Context) error
}

// RecordingSession は録音セッションを表すインターフェースです
type RecordingSession interface {
    // GetID は録音セッションのIDを返します
    GetID() string
    
    // GetStatus は録音の現在のステータスを返します
    GetStatus() RecordingStatus
    
    // Stop は録音を停止します
    Stop(ctx context.Context) error
    
    // GetOutputPath は録音ファイルのパスを返します（完了後）
    GetOutputPath() string
    
    // GetMetadata は録音のメタデータを返します
    GetMetadata() map[string]interface{}
}
```

## 2. 録音マネージャーの実装

録音マネージャーは、レコーダープラグインを管理し、録音リクエストを適切なプラグインに転送する役割を担います：

```go
// RecordingManager は録音を管理するコンポーネントです
type RecordingManager struct {
    db             *gorm.DB
    plugins        map[string]RecorderPlugin
    eventEmitter   EventEmitter
    config         *config.Config
    recordSessions map[string]RecordingSession
    mu             sync.Mutex
}

// NewRecordingManager は新しいRecordingManagerを作成します
func NewRecordingManager(db *gorm.DB, eventEmitter EventEmitter, config *config.Config) *RecordingManager {
    return &RecordingManager{
        db:             db,
        plugins:        make(map[string]RecorderPlugin),
        eventEmitter:   eventEmitter,
        config:         config,
        recordSessions: make(map[string]RecordingSession),
    }
}

// RegisterPlugin はレコーダープラグインを登録します
func (rm *RecordingManager) RegisterPlugin(platform string, plugin RecorderPlugin) {
    rm.plugins[platform] = plugin
}

// StartRecording はスケジュールに基づいて録音を開始します
func (rm *RecordingManager) StartRecording(ctx context.Context, scheduleID string) (string, error) {
    // スケジュールを取得
    var schedule models.Schedule
    if err := rm.db.First(&schedule, "id = ?", scheduleID).Error; err != nil {
        return "", fmt.Errorf("スケジュールが見つかりません: %w", err)
    }
    
    // プラグインを取得
    plugin, ok := rm.plugins[schedule.Platform]
    if !ok {
        return "", fmt.Errorf("プラットフォーム %s のプラグインが見つかりません", schedule.Platform)
    }
    
    // 録音レコードを作成
    record := models.NewRecord(
        schedule.ID,
        schedule.Title,
        time.Now(),
        schedule.Platform,
    )
    record.SetStatus(models.RecordStatusRecording)
    
    if err := rm.db.Create(record).Error; err != nil {
        return "", fmt.Errorf("録音レコードの作成に失敗しました: %w", err)
    }
    
    // 録音開始イベントを発行
    rm.emitEvent("recording.started", map[string]interface{}{
        "record_id":   record.ID,
        "schedule_id": schedule.ID,
        "title":       schedule.Title,
        "platform":    schedule.Platform,
    })
    
    // 録音を開始
    go func() {
        session, err := plugin.Record(context.Background(), schedule)
        if err != nil {
            // 録音失敗
            rm.handleRecordingFailure(record.ID, err)
            return
        }
        
        // セッションを保存
        rm.mu.Lock()
        rm.recordSessions[record.ID] = session
        rm.mu.Unlock()
        
        // 録音完了を監視
        go rm.monitorRecordingSession(record.ID, session)
    }()
    
    return record.ID, nil
}

// StopRecording は録音を停止します
func (rm *RecordingManager) StopRecording(ctx context.Context, recordID string) error {
    rm.mu.Lock()
    session, ok := rm.recordSessions[recordID]
    rm.mu.Unlock()
    
    if !ok {
        return fmt.Errorf("録音セッションが見つかりません: %s", recordID)
    }
    
    return session.Stop(ctx)
}

// handleRecordingFailure は録音失敗時の処理を行います
func (rm *RecordingManager) handleRecordingFailure(recordID string, err error) {
    // レコードを更新
    var record models.Record
    if dbErr := rm.db.First(&record, "id = ?", recordID).Error; dbErr != nil {
        log.Printf("録音レコードの取得に失敗しました: %v", dbErr)
        return
    }
    
    record.SetStatus(models.RecordStatusFailed)
    record.SetErrorMessage(err.Error())
    
    if dbErr := rm.db.Save(&record).Error; dbErr != nil {
        log.Printf("録音レコードの更新に失敗しました: %v", dbErr)
    }
    
    // 録音失敗イベントを発行
    rm.emitEvent("recording.failed", map[string]interface{}{
        "record_id": record.ID,
        "title":     record.Title,
        "platform":  record.Platform,
        "error":     err.Error(),
    })
}

// monitorRecordingSession は録音セッションを監視します
func (rm *RecordingManager) monitorRecordingSession(recordID string, session RecordingSession) {
    // セッションのステータスが変わるまで待機
    for session.GetStatus() == models.RecordStatusRecording {
        time.Sleep(5 * time.Second)
    }
    
    // セッションが完了した場合
    if session.GetStatus() == models.RecordStatusCompleted {
        // レコードを更新
        var record models.Record
        if err := rm.db.First(&record, "id = ?", recordID).Error; err != nil {
            log.Printf("録音レコードの取得に失敗しました: %v", err)
            return
        }
        
        outputPath := session.GetOutputPath()
        
        // ファイルサイズを取得
        fileInfo, err := os.Stat(outputPath)
        var fileSize int64
        if err == nil {
            fileSize = fileInfo.Size()
        }
        
        record.SetStatus(models.RecordStatusCompleted)
        record.SetFilePath(outputPath)
        record.SetFileSize(fileSize)
        
        // メタデータから期間を取得
        metadata := session.GetMetadata()
        if duration, ok := metadata["duration"].(int); ok {
            record.SetDuration(duration)
        }
        
        if err := rm.db.Save(&record).Error; err != nil {
            log.Printf("録音レコードの更新に失敗しました: %v", err)
        }
        
        // 録音完了イベントを発行
        rm.emitEvent("recording.completed", map[string]interface{}{
            "record_id":  record.ID,
            "title":      record.Title,
            "platform":   record.Platform,
            "file_path":  outputPath,
            "file_size":  fileSize,
            "duration":   record.Duration,
        })
    }
    
    // セッションを削除
    rm.mu.Lock()
    delete(rm.recordSessions, recordID)
    rm.mu.Unlock()
}

// emitEvent はイベントを発行します
func (rm *RecordingManager) emitEvent(eventType string, data map[string]interface{}) {
    if rm.eventEmitter != nil {
        if err := rm.eventEmitter.EmitEvent(eventType, data); err != nil {
            log.Printf("イベントの発行に失敗しました: %v", err)
        }
    }
}
```

## 3. API経由での即時録音開始

API経由で即時録音を開始するためのエンドポイントを実装します：

```go
// APIサーバーの設定
func setupRoutes(app *fiber.App, recordingManager *RecordingManager) {
    api := app.Group("/api")
    
    // 録音関連のエンドポイント
    recordings := api.Group("/recordings")
    
    // 録音開始エンドポイント
    recordings.Post("/", func(c *fiber.Ctx) error {
        // リクエストボディをパース
        var req struct {
            ScheduleID string `json:"schedule_id"`
        }
        
        if err := c.BodyParser(&req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "無効なリクエスト形式です",
            })
        }
        
        // 録音を開始
        recordID, err := recordingManager.StartRecording(c.Context(), req.ScheduleID)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        
        return c.Status(fiber.StatusCreated).JSON(fiber.Map{
            "record_id": recordID,
            "message":   "録音を開始しました",
        })
    })
    
    // 録音停止エンドポイント
    recordings.Delete("/:id", func(c *fiber.Ctx) error {
        recordID := c.Params("id")
        
        // 録音を停止
        if err := recordingManager.StopRecording(c.Context(), recordID); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        
        return c.JSON(fiber.Map{
            "message": "録音を停止しました",
        })
    })
    
    // 録音状態取得エンドポイント
    recordings.Get("/:id", func(c *fiber.Ctx) error {
        recordID := c.Params("id")
        
        // 録音レコードを取得
        var record models.Record
        if err := recordingManager.db.First(&record, "id = ?", recordID).Error; err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "録音レコードが見つかりません",
            })
        }
        
        return c.JSON(record)
    })
    
    // 録音リスト取得エンドポイント
    recordings.Get("/", func(c *fiber.Ctx) error {
        var records []models.Record
        query := recordingManager.db.Order("created_at DESC")
        
        // フィルタリング
        if platform := c.Query("platform"); platform != "" {
            query = query.Where("platform = ?", platform)
        }
        
        if status := c.Query("status"); status != "" {
            query = query.Where("status = ?", status)
        }
        
        // ページネーション
        page, _ := strconv.Atoi(c.Query("page", "1"))
        pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
        offset := (page - 1) * pageSize
        
        if err := query.Limit(pageSize).Offset(offset).Find(&records).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "録音レコードの取得に失敗しました",
            })
        }
        
        return c.JSON(records)
    })
}
```

## 4. 録音チェックワーカーの実装

録音チェックワーカーは、定期的に実行され、録音予定のレコードを検索し、録音マネージャーを通じて録音を開始します：

```go
// runRecordingCheckWorker は録音チェックワーカーを実行します
func runRecordingCheckWorker(ctx context.Context, db *gorm.DB, recordingManager *RecordingManager, config *config.Config) {
    log.Println("録音チェックワーカーを起動しました")
    
    // 録音チェック間隔（秒）
    interval := 60
    if config != nil && config.Recording.ScheduleCheckIntervalSeconds > 0 {
        interval = config.Recording.ScheduleCheckIntervalSeconds
    }
    
    ticker := time.NewTicker(time.Duration(interval) * time.Second)
    defer ticker.Stop()
    
    // 初回実行
    checkRecordings(ctx, db, recordingManager)
    
    for {
        select {
        case <-ctx.Done():
            log.Println("録音チェックワーカーを終了します")
            return
        case <-ticker.C:
            checkRecordings(ctx, db, recordingManager)
        }
    }
}

// checkRecordings は録音をチェックします
func checkRecordings(ctx context.Context, db *gorm.DB, recordingManager *RecordingManager) {
    log.Println("録音をチェックしています...")
    
    // 現在時刻
    now := time.Now()
    
    // 録音予定のスケジュールを検索
    // 1. IsProcessingがfalseのもの
    // 2. 開始時間が現在時刻から前後2分以内のもの
    var schedules []models.Schedule
    if err := db.Where("is_processing = ? AND start_time BETWEEN ? AND ?",
        false,
        now.Add(-2*time.Minute),
        now.Add(2*time.Minute),
    ).Find(&schedules).Error; err != nil {
        log.Printf("スケジュールの検索に失敗しました: %v", err)
        return
    }
    
    for _, schedule := range schedules {
        // スケジュールを処理中に設定
        schedule.SetIsProcessing(true)
        if err := db.Save(&schedule).Error; err != nil {
            log.Printf("スケジュールの更新に失敗しました: %v", err)
            continue
        }
        
        // 録音を開始
        recordID, err := recordingManager.StartRecording(ctx, schedule.ID)
        if err != nil {
            log.Printf("録音の開始に失敗しました: %v", err)
            
            // スケジュールを未処理に戻す
            schedule.SetIsProcessing(false)
            if dbErr := db.Save(&schedule).Error; dbErr != nil {
                log.Printf("スケジュールの更新に失敗しました: %v", dbErr)
            }
            continue
        }
        
        log.Printf("録音を開始しました: スケジュールID=%s, レコードID=%s", schedule.ID, recordID)
    }
}
```

## 5. メインアプリケーションでの統合

最後に、これらのコンポーネントをメインアプリケーションで統合します：

```go
// runDaemon はデーモンを実行します
func runDaemon() error {
    log.Println("Rec-adio Daemon を起動しています...")
    
    // データベース接続
    db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("データベースへの接続に失敗しました: %w", err)
    }
    
    // イベントエミッターを作成
    eventEmitter := event.NewEventEmitter()
    
    // 録音マネージャーを作成
    recordingManager := recording.NewRecordingManager(db, eventEmitter, cfg)
    
    // レコーダープラグインを登録
    if cfg.Platforms.AGQR.Enabled {
        agqrPlugin := recorders.NewAGQRRecorder(cfg.Platforms.AGQR)
        if err := agqrPlugin.Initialize(ctx, cfg); err != nil {
            log.Printf("AGQRプラグインの初期化に失敗しました: %v", err)
        } else {
            recordingManager.RegisterPlugin("agqr", agqrPlugin)
        }
    }
    
    if cfg.Platforms.TwitterSpace.Enabled {
        twitterSpacePlugin := recorders.NewTwitterSpaceRecorder(cfg.Platforms.TwitterSpace)
        if err := twitterSpacePlugin.Initialize(ctx, cfg); err != nil {
            log.Printf("Twitter Spaceプラグインの初期化に失敗しました: %v", err)
        } else {
            recordingManager.RegisterPlugin("twitter_space", twitterSpacePlugin)
        }
    }
    
    if cfg.Platforms.BiliBili.Enabled {
        bilibiliPlugin := recorders.NewBiliBiliRecorder(cfg.Platforms.BiliBili)
        if err := bilibiliPlugin.Initialize(ctx, cfg); err != nil {
            log.Printf("BiliBiliプラグインの初期化に失敗しました: %v", err)
        } else {
            recordingManager.RegisterPlugin("bilibili", bilibiliPlugin)
        }
    }
    
    // 通知プラグインを登録
    if cfg.Notification.Enabled {
        if cfg.Notification.Line.Enabled {
            linePlugin := notifications.NewLineNotificationPlugin(cfg.Notification.Line)
            eventEmitter.AddListener(linePlugin)
        }
        
        if cfg.Notification.Webhook.Enabled {
            webhookPlugin := notifications.NewWebhookNotificationPlugin(cfg.Notification.Webhook)
            eventEmitter.AddListener(webhookPlugin)
        }
    }
    
    // APIサーバーを起動
    app := fiber.New()
    setupRoutes(app, recordingManager)
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        log.Printf("APIサーバーを起動しています: %s:%d", cfg.App.API.Host, cfg.App.API.Port)
        addr := fmt.Sprintf("%s:%d", cfg.App.API.Host, cfg.App.API.Port)
        if err := app.Listen(addr); err != nil {
            log.Printf("APIサーバーの起動に失敗しました: %v", err)
        }
    }()
    
    // コンテキストのキャンセルを監視
    go func() {
        <-ctx.Done()
        log.Println("APIサーバーをシャットダウンしています...")
        
        // シャットダウンのタイムアウトを設定
        shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        if err := app.ShutdownWithContext(shutdownCtx); err != nil {
            log.Printf("APIサーバーのシャットダウンに失敗しました: %v", err)
        } else {
            log.Println("APIサーバーをシャットダウンしました")
        }
    }()
    
    // スケジュール更新ワーカーを起動
    wg.Add(1)
    go func() {
        defer wg.Done()
        runScheduleUpdateWorker(ctx, db, cfg)
    }()
    
    // 録音チェックワーカーを起動
    wg.Add(1)
    go func() {
        defer wg.Done()
        runRecordingCheckWorker(ctx, db, recordingManager, cfg)
    }()
    
    // クリーンアップワーカーを起動
    wg.Add(1)
    go func() {
        defer wg.Done()
        runCleanupWorker(ctx, db, cfg)
    }()
    
    // 全てのワーカーが終了するまで待機
    wg.Wait()
    
    log.Println("Rec-adio Daemon を終了しました")
    return nil
}
```

## 6. API経由での即時録音開始の使用例

API経由で即時録音を開始するには、以下のようなHTTPリクエストを送信します：

```http
POST /api/recordings HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "schedule_id": "sched_1234567890abcdef"
}
```

成功した場合、以下のようなレスポンスが返されます：

```json
{
  "record_id": "rec_abcdef1234567890",
  "message": "録音を開始しました"
}
```

このAPIを使用することで、スケジュールに基づいた録音だけでなく、ユーザーが手動で即時に録音を開始することも可能になります。

## 7. 実装上の注意点

1. **並行処理の管理**：
   - 複数の録音が同時に実行される可能性があるため、適切な並行処理の管理が必要です
   - ミューテックスやチャネルを使用して、共有リソースへのアクセスを制御します

2. **エラーハンドリング**：
   - 録音プロセスでは様々なエラーが発生する可能性があるため、適切なエラーハンドリングが重要です
   - エラーが発生した場合は、ログに記録し、適切なステータスに更新します

3. **リソース管理**：
   - 長時間の録音では、メモリリークやファイルディスクリプタのリークに注意が必要です
   - 録音プロセスが終了した際に、適切にリソースを解放します

4. **設定の柔軟性**：
   - 各プラグインの設定は、設定ファイルから読み込むようにします
   - ユーザーが必要に応じて設定を変更できるようにします

5. **テスト可能性**：
   - インターフェースを使用することで、モックを使用したテストが容易になります
   - 各コンポーネントを独立してテストできるようにします
