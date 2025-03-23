package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config は、アプリケーション全体の設定を表します
type Config struct {
	// 基本設定
	App AppConfig `toml:"app"`

	// データベース設定
	Database DatabaseConfig `toml:"database"`

	// 録音設定
	Recording RecordingConfig `toml:"recording"`

	// 通知設定
	Notification NotificationConfig `toml:"notification"`

	// プラットフォーム設定
	Platforms PlatformsConfig `toml:"platforms"`
}

// AppConfig は、アプリケーションの基本設定を表します
type AppConfig struct {
	// アプリケーション名
	Name string `toml:"name"`

	// 環境（development, production）
	Environment string `toml:"environment"`

	// ログレベル（debug, info, warn, error）
	LogLevel string `toml:"log_level"`

	// データディレクトリ
	DataDir string `toml:"data_dir"`

	// 一時ディレクトリ
	TempDir string `toml:"temp_dir"`

	// 録音ファイル保存ディレクトリ
	RecordingsDir string `toml:"recordings_dir"`

	// APIサーバー設定
	API APIConfig `toml:"api"`
}

// APIConfig は、APIサーバーの設定を表します
type APIConfig struct {
	// ホスト
	Host string `toml:"host"`

	// ポート
	Port int `toml:"port"`

	// 基本パス
	BasePath string `toml:"base_path"`

	// CORSの許可オリジン
	AllowOrigins []string `toml:"allow_origins"`

	// JWT認証シークレット
	JWTSecret string `toml:"jwt_secret"`
}

// DatabaseConfig は、データベースの設定を表します
type DatabaseConfig struct {
	// データベースの種類（sqlite, mysql, postgres）
	Type string `toml:"type"`

	// データベースファイルのパス（SQLite用）
	Path string `toml:"path"`

	// ホスト（MySQL, PostgreSQL用）
	Host string `toml:"host"`

	// ポート（MySQL, PostgreSQL用）
	Port int `toml:"port"`

	// データベース名（MySQL, PostgreSQL用）
	Name string `toml:"name"`

	// ユーザー名（MySQL, PostgreSQL用）
	User string `toml:"user"`

	// パスワード（MySQL, PostgreSQL用）
	Password string `toml:"password"`

	// 最大接続数
	MaxOpenConns int `toml:"max_open_conns"`

	// 最大アイドル接続数
	MaxIdleConns int `toml:"max_idle_conns"`

	// 接続の最大生存時間（秒）
	ConnMaxLifetime int `toml:"conn_max_lifetime"`
}

// RecordingConfig は、録音の設定を表します
type RecordingConfig struct {
	// 録音前の準備時間（秒）
	PrepareSeconds int `toml:"prepare_seconds"`

	// 録音後の余裕時間（秒）
	ExtraSeconds int `toml:"extra_seconds"`

	// 同時録音の最大数
	MaxConcurrentRecordings int `toml:"max_concurrent_recordings"`

	// 録音ファイルの形式（mp3, aac, wav）
	Format string `toml:"format"`

	// 録音ファイルのビットレート
	Bitrate int `toml:"bitrate"`

	// 録音ファイルのサンプルレート
	SampleRate int `toml:"sample_rate"`

	// 録音ファイルのチャンネル数
	Channels int `toml:"channels"`

	// 録音ファイルの命名パターン
	FileNamePattern string `toml:"file_name_pattern"`

	// 録音ファイルの保存期間（日）
	RetentionDays int `toml:"retention_days"`

	// 録音スケジュールの更新間隔（分）
	ScheduleUpdateIntervalMinutes int `toml:"schedule_update_interval_minutes"`
}

// NotificationConfig は、通知の設定を表します
type NotificationConfig struct {
	// 通知を有効にするかどうか
	Enabled bool `toml:"enabled"`

	// 録音開始時に通知するかどうか
	NotifyOnStart bool `toml:"notify_on_start"`

	// 録音完了時に通知するかどうか
	NotifyOnComplete bool `toml:"notify_on_complete"`

	// 録音失敗時に通知するかどうか
	NotifyOnError bool `toml:"notify_on_error"`

	// LINE通知設定
	Line LineNotificationConfig `toml:"line"`

	// Webhook通知設定
	Webhook WebhookNotificationConfig `toml:"webhook"`
}

// LineNotificationConfig は、LINE通知の設定を表します
type LineNotificationConfig struct {
	// 有効かどうか
	Enabled bool `toml:"enabled"`

	// アクセストークン
	AccessToken string `toml:"access_token"`
}

// WebhookNotificationConfig は、Webhook通知の設定を表します
type WebhookNotificationConfig struct {
	// 有効かどうか
	Enabled bool `toml:"enabled"`

	// WebhookのURL
	URL string `toml:"url"`

	// 認証トークン
	AuthToken string `toml:"auth_token"`
}

// PlatformsConfig は、各プラットフォームの設定を表します
type PlatformsConfig struct {
	// A&G+ (AGQR) 設定
	AGQR AGQRConfig `toml:"agqr"`

	// Radiko設定
	Radiko RadikoConfig `toml:"radiko"`

	// 音泉設定
	Onsen OnsenConfig `toml:"onsen"`

	// 響設定
	Hibiki HibikiConfig `toml:"hibiki"`

	// Twitter Space設定
	TwitterSpace TwitterSpaceConfig `toml:"twitter_space"`

	// BiliBili設定
	BiliBili BiliBiliConfig `toml:"bilibili"`
}

// AGQRConfig は、A&G+ (AGQR) の設定を表します
type AGQRConfig struct {
	// 有効かどうか
	Enabled bool `toml:"enabled"`

	// ストリームURL
	StreamURL string `toml:"stream_url"`
}

// RadikoConfig は、Radikoの設定を表します
type RadikoConfig struct {
	// 有効かどうか
	Enabled bool `toml:"enabled"`

	// ユーザーエージェント
	UserAgent string `toml:"user_agent"`

	// プレミアム会員かどうか
	IsPremium bool `toml:"is_premium"`

	// メールアドレス（プレミアム会員用）
	Email string `toml:"email"`

	// パスワード（プレミアム会員用）
	Password string `toml:"password"`
}

// OnsenConfig は、音泉の設定を表します
type OnsenConfig struct {
	// 有効かどうか
	Enabled bool `toml:"enabled"`
}

// HibikiConfig は、響の設定を表します
type HibikiConfig struct {
	// 有効かどうか
	Enabled bool `toml:"enabled"`
}

// TwitterSpaceConfig は、Twitter Spaceの設定を表します
type TwitterSpaceConfig struct {
	// 有効かどうか
	Enabled bool `toml:"enabled"`

	// APIキー
	APIKey string `toml:"api_key"`

	// APIシークレット
	APISecret string `toml:"api_secret"`

	// アクセストークン
	AccessToken string `toml:"access_token"`

	// アクセストークンシークレット
	AccessTokenSecret string `toml:"access_token_secret"`

	// ベアラートークン
	BearerToken string `toml:"bearer_token"`
}

// BiliBiliConfig は、BiliBiliの設定を表します
type BiliBiliConfig struct {
	// 有効かどうか
	Enabled bool `toml:"enabled"`

	// セッションID（SESSDATA）
	SessionData string `toml:"session_data"`

	// ビリジェクト（bili_jct）
	BiliJCT string `toml:"bili_jct"`

	// ユーザーID（DedeUserID）
	DedeUserID string `toml:"dede_user_id"`
}

// LoadConfig は、指定されたパスから設定を読み込みます
func LoadConfig(path string) (*Config, error) {
	// ファイルが存在するか確認
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", path)
	}

	// 設定を読み込む
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// 設定を検証
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// パスを絶対パスに変換
	if err := normalizePaths(&config); err != nil {
		return nil, fmt.Errorf("failed to normalize paths: %w", err)
	}

	return &config, nil
}

// validateConfig は、設定を検証します
func validateConfig(config *Config) error {
	// アプリケーション名が設定されているか確認
	if config.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}

	// 環境が設定されているか確認
	if config.App.Environment == "" {
		config.App.Environment = "development"
	}

	// ログレベルが設定されているか確認
	if config.App.LogLevel == "" {
		config.App.LogLevel = "info"
	}

	// データベースの種類が設定されているか確認
	if config.Database.Type == "" {
		config.Database.Type = "sqlite"
	}

	// SQLiteの場合、パスが設定されているか確認
	if config.Database.Type == "sqlite" && config.Database.Path == "" {
		return fmt.Errorf("database.path is required for sqlite")
	}

	// MySQL/PostgreSQLの場合、必要な情報が設定されているか確認
	if (config.Database.Type == "mysql" || config.Database.Type == "postgres") &&
		(config.Database.Host == "" || config.Database.Name == "" || config.Database.User == "") {
		return fmt.Errorf("database.host, database.name, and database.user are required for mysql/postgres")
	}

	return nil
}

// normalizePaths は、設定内のパスを絶対パスに変換します
func normalizePaths(config *Config) error {
	// データディレクトリ
	if config.App.DataDir != "" {
		absPath, err := filepath.Abs(config.App.DataDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for data_dir: %w", err)
		}
		config.App.DataDir = absPath
	}

	// 一時ディレクトリ
	if config.App.TempDir != "" {
		absPath, err := filepath.Abs(config.App.TempDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for temp_dir: %w", err)
		}
		config.App.TempDir = absPath
	}

	// 録音ファイル保存ディレクトリ
	if config.App.RecordingsDir != "" {
		absPath, err := filepath.Abs(config.App.RecordingsDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for recordings_dir: %w", err)
		}
		config.App.RecordingsDir = absPath
	}

	// SQLiteデータベースファイルのパス
	if config.Database.Type == "sqlite" && config.Database.Path != "" {
		absPath, err := filepath.Abs(config.Database.Path)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for database.path: %w", err)
		}
		config.Database.Path = absPath
	}

	return nil
}

// DefaultConfig は、デフォルトの設定を返します
func DefaultConfig() *Config {
	// 現在のディレクトリを取得
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	// データディレクトリ
	dataDir := filepath.Join(cwd, "data")

	// 一時ディレクトリ
	tempDir := filepath.Join(dataDir, "temp")

	// 録音ファイル保存ディレクトリ
	recordingsDir := filepath.Join(dataDir, "recordings")

	// データベースファイルのパス
	dbPath := filepath.Join(dataDir, "recadio.db")

	return &Config{
		App: AppConfig{
			Name:          "Rec-adio",
			Environment:   "development",
			LogLevel:      "info",
			DataDir:       dataDir,
			TempDir:       tempDir,
			RecordingsDir: recordingsDir,
			API: APIConfig{
				Host:         "localhost",
				Port:         8080,
				BasePath:     "/api",
				AllowOrigins: []string{"*"},
				JWTSecret:    "recadio-secret",
			},
		},
		Database: DatabaseConfig{
			Type:            "sqlite",
			Path:            dbPath,
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 3600,
		},
		Recording: RecordingConfig{
			PrepareSeconds:                30,
			ExtraSeconds:                  30,
			MaxConcurrentRecordings:       3,
			Format:                        "mp3",
			Bitrate:                       192,
			SampleRate:                    44100,
			Channels:                      2,
			FileNamePattern:               "{date}_{time}_{title}_{platform}",
			RetentionDays:                 30,
			ScheduleUpdateIntervalMinutes: 60,
		},
		Notification: NotificationConfig{
			Enabled:          false,
			NotifyOnStart:    true,
			NotifyOnComplete: true,
			NotifyOnError:    true,
			Line: LineNotificationConfig{
				Enabled: false,
			},
			Webhook: WebhookNotificationConfig{
				Enabled: false,
			},
		},
		Platforms: PlatformsConfig{
			AGQR: AGQRConfig{
				Enabled:   false,
				StreamURL: "https://fms2.uniqueradio.jp/agqr10/aandg1.m3u8",
			},
			Radiko: RadikoConfig{
				Enabled:   false,
				UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36",
				IsPremium: false,
			},
			Onsen: OnsenConfig{
				Enabled: false,
			},
			Hibiki: HibikiConfig{
				Enabled: false,
			},
			TwitterSpace: TwitterSpaceConfig{
				Enabled: false,
			},
			BiliBili: BiliBiliConfig{
				Enabled: false,
			},
		},
	}
}

// SaveConfig は、設定をファイルに保存します
func SaveConfig(config *Config, path string) error {
	// ディレクトリが存在するか確認
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// ディレクトリが存在しない場合は作成
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// ファイルを作成
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	// 設定をエンコード
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}
