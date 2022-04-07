import Vapor
import Queues

struct RecordingScheduleJob: AsyncScheduledJob {
    func run(context: QueueContext) async throws {
        context.logger.debug("RecordingScheduleJob - start. \(Date())")
        let db = context.application.db
        // DBを検索して録音するべき番組を取得
        // 現時刻から2分(1分マージン)以内に始まる番組 → between now() and now()+2min
        let calendar = Calendar(identifier: .gregorian)
        let now = Date()
        let searchEndAt = calendar.date(byAdding: .minute, value: 2, to: now)!
        let schedules = try await Schedule.query(on: db)
            .filter(\.$isProcessing, .equal, false)
            .filter(\.$startDatetime, .greaterThanOrEqual, now)
            .filter(\.$startDatetime, .lessThan, searchEndAt)
            .all()
        for schedule in schedules {
            try await context.queue.dispatch(RecordingJob.self, schedule, delayUntil: schedule.startDatetime)
            schedule.isProcessing = true
            try await schedule.save(on: context.application.db)
        }
    }
}
