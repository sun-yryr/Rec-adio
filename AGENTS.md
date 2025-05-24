# Recoto開発者ガイド

## プロジェクト概要

RecotoはGo言語で実装された自動録音アプリケーションです。gRPCサーバーとして動作し、NATSメッセージブローカーを使用するイベント駆動アーキテクチャを採用しています。

## アーキテクチャ

- **cmd/**: アプリケーションのエントリポイント
  - `recotod/`: gRPCサーバーや録音のメインプロセス
  - `recoto/`: クライアント
- **pkg/**: 外部に公開される再利用可能なパッケージ
  - `api/`: protobufで生成されたgRPC API定義
- **internal/**: 内部専用パッケージ
  - `adapter/`: gRPCサーバー実装など外部とのインターフェース
  - `broker/`: NATSブローカー関連
  - `config/`: 設定管理
  - `domain/`: ドメインモデル
  - `event/`: イベント定義
  - `eventutil/`: イベント関連ユーティリティ
  - `fileutil/`: ファイル操作ユーティリティ
  - `logger/`: ログ機能
  - `recorder/`: 録音機能の実装
- **api/**: protobuf定義ファイル

## 開発環境セットアップ

### 依存関係のインストール

```bash
go mod tidy
```

### 主要ツール

- golangci-lint: 静的解析とフォーマット
- buf: protobufのlintとコード生成
- testify: テストフレームワーク

## 開発ワークフロー

### コード品質チェック

1. **フォーマット**: `golangci-lint fmt`
2. **静的解析**: `golangci-lint run`
3. **テスト**: `go test ./... -v -coverprofile=coverage.out`
4. **ビルド**: `go build ./...`

### protobuf関連

- **lint**: `buf lint`
- **コード生成**: `buf generate`
- protobufを編集した際は必ずlintとgenerateを実行してください

## コーディング規約

### 基本ルール
- すべてのコードは `golangci-lint` でエラー・警告がないこと
- 新規パッケージ・ファイルはGoの命名規則（小文字、アンダースコアなし、意味のある名前）に従う
- 外部公開APIは `pkg/` 、内部専用は `internal/` に配置

### エラー処理

- `github.com/cockroachdb/errors` パッケージを使用してエラーをラップ
- パッケージレベルのエラー変数を再利用
- コンテキスト付きエラーハンドリングを実装

### ログ出力

- `log/slog` を使用（ `go.uber.org/zap` から移行中）
- 構造化ログを心がける
- ロガーは `slog.SetDefault` でグローバルに設定し、使い回す
- TraceIDなど、全てのログに含めたい値はcontextにセットし、slogのHandlerで回収する

## テスト要件

### テストフレームワーク

- `github.com/stretchr/testify` を使用
- `assert`: 一般的なアサーション（テスト継続）
- `require`: 重要な条件（失敗時テスト中断）

### テストファイル命名規則

- `*_test.go` ファイルにテストを記述
- テスト関数は `TestXxx(t *testing.T)` 形式

### テストの書き方

```go
func TestXxx(t *testing.T) {
    t.Parallel() // 並列実行を有効化

    // テーブル駆動テストを推奨
    tests := []struct {
        name     string
        input    string
        expected string
        isValid  bool
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            
            // 早期リターンでネストを減らす
            if !tt.isValid {
                assert.Error(t, err)
                return
            }
            
            require.NoError(t, err)
            assert.Equal(t, tt.expected, actual)
        })
    }
}
```

### 必須テスト

- すべての新規機能・修正には対応するテストコードを追加
- テストは `go test ./...` で全てパスすること
- カバレッジを意識したテスト作成

## 検証手順

### 開発時の必須チェック

1. **依存関係**: `go mod tidy` （Codexはネットワークアクセスがないため実行しない）
2. **ビルド**: `go build ./...`
3. **テスト**: `go test ./... -v`
4. **静的解析**: `golangci-lint run`
5. **protobuf（該当時）**: `buf lint && buf generate`

## PR作成ガイドライン

### PRタイトル

```
[component] 簡潔な変更内容の説明
```

### PR作成前チェックリスト

- [ ] すべてのテストがパス（`go test ./...`）
- [ ] golangci-lintでエラー・警告なし
- [ ] 新機能にはテストを追加
- [ ] protobuf変更時はbuf lintとgenerateを実行
- [ ] 依存関係の変更時は`go mod tidy`実行
