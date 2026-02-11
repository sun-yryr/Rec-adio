# 04. アーキテクチャ設計

更新日: 2026-02-10

## 1. 採用方針

- 廃止: NATS ベースの外部イベントブローカー依存
- 採用:
  - インプロセスのスケジューラー（メモリ）
  - インプロセスの実行キュー（メモリ）
  - 状態永続化として SQLite（ファイル）

## 2. 構成コンポーネント

1. `API Adapter`  
gRPC/CLI からの要求を UseCase へ渡す

2. `UseCase`  
ジョブ作成・更新・停止・削除・実行キャンセル・状態参照を調停する

3. `Scheduler`  
`active` ジョブを監視し、到達時刻で `Run` を生成してキューへ投入する  
初期リリースでは単発予約のみを扱う

4. `InMemory Queue`  
`Run` を FIFO で保持。最大長を設定可能にする

5. `Worker Pool`  
キューを消費して ffmpeg 実行。実行結果を SQLite へ反映する

6. `SQLite Repository`  
Job/Run の永続化、状態遷移の原子更新、起動時復旧クエリを担当する

## 3. 実行フロー

1. API で Job を作成
2. Job を SQLite に保存
3. Scheduler が `scheduled_at` 到達を検知
4. Run を `queued` で保存し、メモリキューへ投入
5. Worker が Run を `running` に更新して録音開始
6. 完了時に `succeeded/failed/canceled` を保存
7. 状態は API から参照可能

## 4. 起動時リカバリ

- `active` ジョブをロードし、スケジューラーへ再登録
- `queued` の Run は再投入
- `running` の Run は `failed`（理由: daemon interrupted）に遷移し、必要なら再実行ポリシーを適用
- 開始時刻超過ジョブは `catch_up_window=15m` を適用
  - 15分以内: `Run` を新規作成して即時キュー投入
  - 15分超過: 実行せずスキップログを残す

## 5. 並行実行制御

- `worker_count` で同時録音数を制限
- Run の `queued -> running` 遷移は条件付き更新で排他（楽観ロック）
- 初期値は `worker_count=2`

## 6. 録音形式

- 初期リリースは `m4a(aac)` 固定
- 形式選択は将来拡張項目とし、現時点では設計対象外

## 7. ディレクトリ案

- `internal/scheduler/`
- `internal/queue/`
- `internal/store/sqlite/`
- `internal/usecase/`
- `internal/adapter/grpcserver/`（既存拡張）
- `internal/recorder/`（既存流用）
