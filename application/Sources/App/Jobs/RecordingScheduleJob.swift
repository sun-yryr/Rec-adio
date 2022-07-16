import Vapor
import Queues

struct RecordingScheduleJob: AsyncScheduledJob {
    func run(context: QueueContext) async throws {
        sleep(1) // MEMO: https://github.com/vapor/queues/issues/94 の暫定対応
        context.logger.info("RecordingScheduleJob - start. \(Date())")
        let db = context.application.db
        // DBを検索して録音するべき番組を取得
        // 現時刻-2分から2分(1分マージン)以内に始まる番組 → between now()-2min and now()+2min
        let calendar = Calendar(identifier: .gregorian)
        let now = Date()
        let searchStartAt =  calendar.date(byAdding: .minute, value: -2, to: now)!
        let searchEndAt = calendar.date(byAdding: .minute, value: 2, to: now)!
        let schedules = try await Schedule.query(on: db)
            .filter(\.$isProcessing, .equal, false)
            .filter(\.$startDatetime, .greaterThanOrEqual, searchStartAt)
            .filter(\.$startDatetime, .lessThan, searchEndAt)
            .all()
        for schedule in schedules {
            do {
                try await context.queue.dispatch(RecordingJob.self, schedule, delayUntil: schedule.startDatetime)
                schedule.isProcessing = true
                try await schedule.save(on: context.application.db)
            } catch {
                // エラーが発生しても後続は通す
                context.logger.error("RecordingJobの登録に失敗しました。 schedule_id[\(schedule.id ?? -1)]")
            }
        }
    }
}
