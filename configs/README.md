# configs - Rec-adio 設定ファイル

このディレクトリには、Rec-adioの設定ファイルテンプレートが含まれています。

## 設定ファイル

- **recadiod.toml**: コアサービス（デーモン）の設定ファイル
- **recadio.toml**: コマンドラインツールの設定ファイル

## 設定ファイル形式

Rec-adioの設定ファイルは、TOML形式で記述されています。TOML（Tom's Obvious, Minimal Language）は、設定ファイル用の言語で、読みやすく、書きやすく、パースしやすいという特徴があります。

## recadiod.toml

```toml
# グローバル設定
[global]
data_dir = "/path/to/data"
log_level = "info"
api_port = 8080

# 録音設定
[recording]
default_format = "mp3"
max_concurrent_recordings = 3
retry_count = 3
retry_interval = 5

# スケジュール設定
[schedule]
check_interval = 60
advance_notice = 120

# プラグイン設定
[plugins]
enabled = ["agqr", "twitter_space", "bilibili", "notification"]

# 超A&G+プラグイン設定
[plugin.agqr]
timeout = 30
quality = "high"

# Twitter Spaceプラグイン設定
[plugin.twitter_space]
auth_token = "your-auth-token"
check_interval = 300

# BiliBiliプラグイン設定
[plugin.bilibili]
cookie = "your-cookie"
check_interval = 300

# 通知プラグイン設定
[plugin.notification]
enabled_notifications = ["recording.completed", "recording.failed"]

# LINE通知設定
[plugin.notification.line]
enabled = true
token = "your-line-token"

# メール通知設定
[plugin.notification.email]
enabled = false
smtp_server = "smtp.example.com"
from = "recorder@example.com"
to = "user@example.com"
```

## recadio.toml

```toml
# APIクライアント設定
[api]
url = "http://localhost:8080"
timeout = 10

# 出力設定
[output]
format = "table"
color = true

# ログ設定
[log]
level = "info"
file = ""
```

## 設定ファイルの配置場所

設定ファイルは、以下の場所に配置されます（予定）：

### Linux

- `/etc/rec-adio/recadiod.toml`
- `~/.config/rec-adio/recadiod.toml`
- `./recadiod.toml`

### macOS

- `/usr/local/etc/rec-adio/recadiod.toml`
- `~/Library/Application Support/rec-adio/recadiod.toml`
- `./recadiod.toml`

### Windows

- `C:\ProgramData\rec-adio\recadiod.toml`
- `%APPDATA%\rec-adio\recadiod.toml`
- `.\recadiod.toml`

設定ファイルは、上記の順序で検索され、最初に見つかったファイルが使用されます。また、コマンドライン引数で設定ファイルのパスを指定することもできます。
