# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 言語設定

このリポジトリで作業する際は、**日本語で回答してください**。コメントやドキュメントも日本語で記述してください。

## プロジェクト概要

Recoto v4は、Radiko、超A&G+、音泉、響などの日本の音声配信サービスから録音を行うアプリケーションのGo言語による書き直し版です。このバージョンでは、ジョブベースの録音管理とアクティブな録音中の予約変更をサポートするイベント駆動アーキテクチャを導入しています。

## 開発コマンド

### ビルドと実行

```bash
# 全パッケージのビルド
go build ./...

# デーモンの実行
go run ./cmd/recotod/main.go

# CLI実行（プレースホルダー実装）
go run ./cmd/recoto/main.go

# カバレッジ付きテスト
go test ./... -v -coverprofile=coverage.out

# 特定パッケージのテスト
go test ./internal/recorder -v
```

### Protocol Buffers
```bash
# protobufコード生成（bufが必要）
buf generate

# スキーマ検証
buf lint
buf breaking --against '.git#branch=main'
```

### 設定

- デフォルト設定場所: `~/.config/recoto/daemon.toml`
- 設定例: `config_example/daemon.toml`
- 設定ファイルが存在しない場合は自動生成

## アーキテクチャ

### イベント駆動設計

アプリケーションはNATS組み込みブローカーを使用し、3つのコアイベントタイプがあります：

- `RequestedEvent` - 録音リクエスト開始
- `StartedEvent` - 録音実際開始
- `FinishedEvent` - 録音完了/失敗

イベントフロー: CLI/gRPC → EventService → NATS → RecordingManager → URLRecorder

### クリーンアーキテクチャレイヤー

- `cmd/` - エントリーポイント（recoto CLI、recotod デーモン）
- `internal/domain/` - コアエンティティ（Recording、Source）
- `internal/adapter/grpcserver/` - ミドルウェア付きgRPCトランスポート
- `internal/recorder/` - 録音実装（URLRecorderはffmpegを使用）
- `internal/broker/` - メッセージブローカー抽象化（NATS実装）
- `internal/event/` - イベント定義とEventService
- `pkg/api/` - 生成されたprotobufコード

### 主要コンポーネント

- `RecordingManager` - キャンセルサポート付き録音ライフサイクル管理
- `URLRecorder` - FFmpegベースの録音実装
- `EventService` - メタデータ付きNATSイベント発行
- 異なるソースタイプ用のプラガブル録音アーキテクチャ

### 開発ワークフロー

- ブランチは `claude/<slug>` の形式で作成すること。なお、 `slug` は英語です
- 開発中には定期的にフォーマッターとリンターを実行して問題がないことを確認すること
  - フォーマッター: `golangci-lint fmt`
  - リンター: `golangci-lint run`
- 実装が完了したタイミングでテストを実行すること
  - テスト: `go test ./...`
- もしprotoファイルを更新した場合は以下のコマンドを全て実行すること
  - コード生成: `buf generate`
  - スキーマ検証: `buf lint`


### テスト戦略

- テーブル駆動テストを使用した全パッケージのユニットテスト
- RecordingManagerシナリオの統合テスト
- モックブローカーを使用したイベントハンドリングテスト
- CIでoctocovを使用したカバレッジレポート
- 詳細なテストの書き方は @.openhands/microagents/test-style.md を参照してください
