---
name: repo
type: repo
agent: CodeActAgent
---

リポジトリ：Rec-adio
説明：Go言語で実装された自動録音アプリケーション

ディレクトリ構造：
- cmd/：エントリポイント
- pkg/：再利用可能なパッケージ
- internal/：内部専用パッケージ
- api/：gRPCのprotobuf定義
- scripts/：補助スクリプト

セットアップ手順：
- 依存関係のインストール： `go mod tidy`
- 静的解析・フォーマット： `golangci-lint run`
- テスト実行： `go test ./...`

開発ガイドライン：
- すべてのコードは `golangci-lint` でエラー・警告がないこと
- 新規パッケージ・ファイルはGoの命名規則（小文字、アンダースコアなし、意味のある名前）に従う
- 外部公開APIは `pkg/` 、内部専用は `internal/` に配置
- すべての新機能・修正にはテストを追加
- コードレビュー時は必ず `golangci-lint` と `go test` を通すこと
- 依存パッケージの追加・更新・削除時は `go mod tidy` を必ず実行し、 `go.mod` と `go.sum` を最新の状態に保つこと
- protobufを追加・編集した際は `buf lint` でエラー・警告がないことを確認し、 `buf generate` でコード生成を行うこと

テスト要件：
- すべての新規機能・修正には対応するテストコード（`*_test.go`）を追加
- テストは`go test ./...`で全てパスすること
