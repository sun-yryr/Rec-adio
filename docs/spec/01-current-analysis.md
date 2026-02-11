# 01. 現状分析（As-Is）

更新日: 2026-02-10  
対象コミット: `8dac181`

## 1. 既存実装でできること

- `recotod` が gRPC サーバーとして起動できる。
- `RecordingService/StartFromURL` で即時録音要求イベントを発行できる。
- `RecordingManager` が `RequestedEvent` を購読して `URLRecorder` を実行できる。
- `URLRecorder` が ffmpeg を使った録音を実行できる。
- `HealthService` でブローカー疎通チェックができる。

## 2. 既存実装の制約

- CLI (`cmd/recoto`) は未実装。
- 録音ジョブの永続化（再起動耐性）がない。
- 録音予約の CRUD がない。
- 録音状態の参照 API（一覧・詳細）がない。
- イベント基盤が NATS 前提で、単体運用に対して構成が重い。

## 3. 課題（今回の再設計対象）

1. NATS 依存がリリース・運用コストを増やしている。  
2. 「録音中の予約変更」「ジョブ方式」が未実現。  
3. プロセス再起動でメモリ状態（キュー/実行中情報）が失われる。  
4. 仕様が分散しており、実装との差分が可視化しづらい。

## 4. 再設計方針（To-Be）

- メッセージブローカー（NATS）を廃止する。
- デーモン内で以下を内製する。
  - メモリベーススケジューラー（実行時刻管理）
  - メモリベースキュー（実行待ちジョブ）
  - ワーカープール（録音実行）
- 永続化は SQLite（単一ファイル）を採用する。
- 起動時リカバリでジョブと実行履歴を復元する。

## 5. 影響範囲（現コードベース）

- 置き換え・縮小対象:
  - `internal/broker/*`
  - `internal/event/*`（NATSサブジェクト中心の責務）
  - `internal/eventutil/*`
- 継続利用候補:
  - `internal/domain/*`（一部拡張）
  - `internal/recorder/url.go`
  - `internal/fileutil/*`
  - `internal/config/server/*`（設定項目を拡張）
  - `internal/adapter/grpcserver/*`（APIを拡張）

