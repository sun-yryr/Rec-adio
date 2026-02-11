# 07. 実装ロードマップ

更新日: 2026-02-10

## Phase 0: 仕様固定

- 本 `docs/spec` のレビュー完了
- 未決事項を `open-questions.md` で合意

## Phase 1: 永続化基盤

- SQLite ストア実装（jobs/runs）
- マイグレーション管理
- Job/Run の CRUD テスト

完了条件:
- ストア単体テストが通る
- 不正状態遷移が DB レベルで拒否される

## Phase 2: スケジューラー/キュー

- メモリスケジューラー実装
- メモリキュー実装（最大長、バックプレッシャ）
- ワーカープール実装

完了条件:
- 時刻到達で Run が `queued -> running` へ進む
- 同時実行上限が守られる

## Phase 3: 録音実行統合

- `URLRecorder` 統合
- Run の終端状態反映
- キャンセル動作統合

完了条件:
- 成功/失敗/キャンセルの3パターンがE2Eで再現

## Phase 4: API/CLI 拡張

- gRPC の Job/Run 管理 API 追加
- CLI 実装（job create/list/update/pause/resume/delete, run list/cancel）

完了条件:
- CLI からの一連操作で予約録音が完結

## Phase 5: NATS 依存撤去

- `internal/broker`, `internal/event`, `internal/eventutil` の整理/削除
- 不要設定の削除

完了条件:
- go.mod から NATS 依存が除去される
- 既存テストを新設計向けに置換完了

