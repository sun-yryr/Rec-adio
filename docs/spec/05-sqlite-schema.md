# 05. SQLite スキーマ設計

更新日: 2026-02-10

## 1. SQLite 運用設定

- `PRAGMA journal_mode=WAL;`
- `PRAGMA synchronous=NORMAL;`
- `PRAGMA foreign_keys=ON;`
- `PRAGMA busy_timeout=5000;`

## 2. テーブル定義（案）

```sql
CREATE TABLE IF NOT EXISTS jobs (
  job_id TEXT PRIMARY KEY,
  source_type TEXT NOT NULL,
  source_value TEXT NOT NULL,
  title TEXT NOT NULL,
  duration_sec INTEGER NOT NULL CHECK(duration_sec > 0),
  scheduled_at TEXT NOT NULL,     -- ISO8601 UTC
  timezone TEXT NOT NULL,
  state TEXT NOT NULL CHECK(state IN ('active','paused','deleted')),
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_jobs_state_scheduled_at
  ON jobs(state, scheduled_at);

CREATE TABLE IF NOT EXISTS runs (
  run_id TEXT PRIMARY KEY,
  job_id TEXT NOT NULL,
  planned_at TEXT NOT NULL,
  started_at TEXT,
  finished_at TEXT,
  state TEXT NOT NULL CHECK(state IN ('queued','running','succeeded','failed','canceled')),
  output_path TEXT NOT NULL,
  error_message TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  FOREIGN KEY(job_id) REFERENCES jobs(job_id)
);

CREATE INDEX IF NOT EXISTS idx_runs_job_id_created_at
  ON runs(job_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_runs_state_updated_at
  ON runs(state, updated_at DESC);

CREATE TABLE IF NOT EXISTS schema_migrations (
  version INTEGER PRIMARY KEY,
  applied_at TEXT NOT NULL
);
```

## 3. 原子的更新の例

```sql
UPDATE runs
SET state='running', started_at=?, updated_at=?
WHERE run_id=? AND state='queued';
```

影響行数が `1` のときのみ遷移成功とする。

## 4. リカバリクエリ例

- 起動時に再投入すべき Run:

```sql
SELECT run_id FROM runs WHERE state='queued';
```

- 異常終了として回復対象の Run:

```sql
SELECT run_id FROM runs WHERE state='running';
```

