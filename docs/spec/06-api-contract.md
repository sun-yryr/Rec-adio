# 06. API 契約（gRPC）

更新日: 2026-02-10

## 1. 方針

- 既存互換として `StartFromURL` は維持する。
- ジョブ管理 API を追加し、CLI から同APIを利用する。

## 2. 追加予定 RPC

- `CreateJob(CreateJobRequest) returns (CreateJobResponse)`
- `UpdateJob(UpdateJobRequest) returns (UpdateJobResponse)`
- `PauseJob(PauseJobRequest) returns (PauseJobResponse)`
- `ResumeJob(ResumeJobRequest) returns (ResumeJobResponse)`
- `DeleteJob(DeleteJobRequest) returns (DeleteJobResponse)`
- `GetJob(GetJobRequest) returns (GetJobResponse)`
- `ListJobs(ListJobsRequest) returns (ListJobsResponse)`
- `CancelRun(CancelRunRequest) returns (CancelRunResponse)`
- `GetRun(GetRunRequest) returns (GetRunResponse)`
- `ListRuns(ListRunsRequest) returns (ListRunsResponse)`

## 3. 既存 API の再定義

`StartFromURL` は内部で「即時1回実行ジョブ」を生成し、以下を返す:

- `recording_id`（`run_id` として扱う）

## 4. ステータス表現

- Job: `ACTIVE | PAUSED | DELETED`
- Run: `QUEUED | RUNNING | SUCCEEDED | FAILED | CANCELED`

## 5. エラーハンドリング方針

- 入力不正: `InvalidArgument`
- 未存在: `NotFound`
- 状態遷移不正: `FailedPrecondition`
- 実行中競合: `Aborted`
- 内部障害: `Internal`

## 6. 互換性方針

- `recoto.recording.v1` を維持し、後方互換でフィールド追加する
- 破壊的変更が必要な場合は `v2` パッケージを新設する

