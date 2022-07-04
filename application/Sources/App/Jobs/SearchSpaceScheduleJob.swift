import Vapor
import Queues
import TwitterAPIKit

struct SearchSpaceScheduleJob: AsyncScheduledJob {
    let client: TwitterAPIClient
    let calendar = Calendar(identifier: .gregorian)

    func run(context: QueueContext) async throws {
        context.logger.debug("Search space job - start. \(Date())")
        let startDatetime = self.calendar.date(byAdding: .minute, value: 1, to: Date())!
        let twitterUsers = try await TwitterUser.query(on: context.application.db).all()
        let userIds = twitterUsers.map { $0.userId }
        let response = await self.client.v2.getSpacesByCreators(
            .init(userIDs: userIds, spaceFields: [.title, .creatorID])
        ).responseDecodable(type: GetSpacesByCreatorResponse.self)
        guard let response = response.success else {
            return
        }
        
        let activeSpaces = response.data.filter { $0.state == " live" }
        for activeSpace in activeSpaces {
            let isProcessingSchedule = try? await Schedule.query(on: context.application.db)
                .filter(\.$extraField, .equal, activeSpace.id)
                .filter(\.$platform, .equal, "space")
                .first()
            if (isProcessingSchedule != nil) {
                // 存在する場合はスキップ
                continue
            }

            guard let twitterUser = twitterUsers.filter({ $0.userId == activeSpace.creator_id }).first else {
                continue
            }

            let schedule = Schedule(
                title: activeSpace.title,
                startDatetime: startDatetime,
                duration: 0,
                platform: "space",
                programInfoId: twitterUser.$programInfo.id,
                extraField: activeSpace.id
            )
            do {
                try await schedule.create(on: context.application.db)
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

fileprivate struct GetSpacesByCreatorResponse: Codable {
    struct Data: Codable {
        let id: String
        let state: String // enumにしてもいいなぁ
        let title: String
        let creator_id: String
    }

    let data: [Data]
}
