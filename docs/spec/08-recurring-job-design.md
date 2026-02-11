# 08. 繰り返しJob設計（追加設計資料）

更新日: 2026-02-10

## 1. 目的

単発Jobのみを対象にした現行設計へ、後方互換で繰り返しJobを追加するための設計方針を定義する。

## 2. 設計決定（確定）

- 繰り返し表現: `RFC5545 RRULE` 文字列
- 初期サポート範囲: `FREQ`, `INTERVAL`, `BYDAY`, `BYHOUR`, `BYMINUTE`
- misfire（停止中取りこぼし）: `catch_up_window` 内の最新1回のみ補填
- overlap（同一Job未完了時の次回発火）: `skip`
- タイムゾーン: サーバー設定で管理（`server.timezone`）
  - デフォルト値: `Asia/Tokyo`

## 3. ドメインモデル追加

## Job

- 追加属性:
  - `schedule_type`: `once | recurring`
  - `rrule`: string nullable（`recurring` の場合は必須）
  - `next_run_at`: timestamp nullable（次回発火時刻、UTC）

## Run

- 既存属性を維持
- 繰り返し実行でも 1発火 = 1Run として履歴を残す

## 4. 主要ルール

- `schedule_type=once`:
  - 既存 `scheduled_at` を使用
- `schedule_type=recurring`:
  - `next_run_at` 到達でRunを生成
  - 発火処理の最後に次回 `next_run_at` を再計算して保存
- 同一Jobで `queued/running` Run が存在する場合:
  - 新規Runを作らず `skip`

## 5. misfireポリシー

- デーモン起動時、`next_run_at < now` の recurring Job を評価
- `catch_up_window` 内に複数候補がある場合でも **最新1回のみ** 実行
- `catch_up_window` を超える候補は補填しない（スキップ）

## 6. API設計（gRPC）

- `CreateJob` / `UpdateJob` に以下を追加:
  - `schedule_type`
  - `rrule`（recurring時必須）
  - `scheduled_at`（once時利用）
- バリデーション:
  - once: `scheduled_at` 必須, `rrule` 禁止
  - recurring: `rrule` 必須
- `GetJob` / `ListJobs` に以下を追加:
  - `schedule_type`
  - `rrule`
  - `next_run_at`
- `StartFromURL` は既存互換を維持（内部では即時 once Job 扱い）

## 7. SQLite設計

`jobs` テーブルに追加:

- `schedule_type TEXT NOT NULL CHECK(schedule_type IN ('once','recurring')) DEFAULT 'once'`
- `rrule TEXT NULL`
- `next_run_at TEXT NULL`

制約:

- `schedule_type='recurring'` の場合は `rrule` 必須
- `schedule_type='once'` の場合は `scheduled_at` 必須

インデックス追加:

- `idx_jobs_state_next_run_at (state, next_run_at)`

## 8. 設定設計

`daemon.toml` 追加項目:

- `[server] timezone = "Asia/Tokyo"`
- `[recording] catch_up_window = "15m"`
- `[worker] count = 2`

## 9. テスト観点

- RRULEバリデーション（許可項目/非対応項目）
- 次回発火計算（daily/weekly + interval + byday/time）
- misfire補填（最新1回のみ）
- overlap時skip確認
- API入力バリデーション（once/recurringの組み合わせ）
- 再起動リカバリで `next_run_at` とRun生成整合

## 10. 将来拡張（初期非対応）

- RFC5545の追加項目（`BYMONTHDAY`, `COUNT`, `UNTIL` など）
- Job単位タイムゾーン指定
- overlapポリシーの切替（queue one / parallel）
- misfire全補填モード

