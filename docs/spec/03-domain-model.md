# 03. ドメインモデル

更新日: 2026-02-10

## 1. エンティティ

## Job

- 録音予約の定義
- 主な属性:
  - `job_id`
  - `source_type` (`url`)
  - `source_value` (録音元URL)
  - `title`
  - `duration_sec`
  - `scheduled_at` (UTC)
  - `timezone`
  - `state`
  - `created_at`, `updated_at`

## Run

- ジョブの実行インスタンス（1回分）
- 主な属性:
  - `run_id`
  - `job_id`
  - `planned_at`
  - `started_at`, `finished_at`
  - `state`
  - `output_path`
  - `error_message`
  - `created_at`, `updated_at`

## 2. 状態

## JobState

- `active`: スケジュール対象
- `paused`: 一時停止（スケジュール対象外）
- `deleted`: 削除済み（論理削除）

## RunState

- `queued`: 実行待ちキュー投入済み
- `running`: 録音実行中
- `succeeded`: 正常終了
- `failed`: 異常終了
- `canceled`: ユーザーまたはシステムで中断

## 3. 状態遷移ルール

- Job:
  - `active -> paused`
  - `paused -> active`
  - `active|paused -> deleted`
- Run:
  - `queued -> running`
  - `running -> succeeded|failed|canceled`
  - `queued -> canceled`（実行前キャンセル）

不正遷移は永続化層で拒否する。

## 4. 集約境界

- `Job` が集約ルート
- `Run` は `Job` 配下の履歴として扱う
- `Job` 更新時に次回実行計算（必要なら）を再評価する

## 5. ドメインサービス

- `SchedulerService`
  - メモリ内で `scheduled_at` を監視し、時刻到達で `Run` を作成してキューへ投入
- `QueueService`
  - `Run` の投入/取り出しを管理
- `ExecutionService`
  - `Run` を `running` に遷移し、Recorder実行後に終端状態へ遷移
- `RecoveryService`
  - 起動時に `queued/running` 残骸を整合回復

## 6. 不変条件

- `duration_sec > 0`
- `source_value` は有効URL
- `deleted` ジョブは新規 `Run` を生成しない
- `running` な `Run` は同時に同一 `run_id` で複数存在しない

