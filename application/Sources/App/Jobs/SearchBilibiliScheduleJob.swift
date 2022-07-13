import Vapor
import Queues
import FluentKit
#if canImport(FoundationNetworking)
    import FoundationNetworking
#endif

struct SearchBilibiliScheduleJob: AsyncScheduledJob {
    let calendar = Calendar(identifier: .gregorian)

    func run(context: QueueContext) async throws {
        context.logger.debug("Search bilibili job - start. \(Date())")
        let startDatetime = self.calendar.date(byAdding: .minute, value: 1, to: Date())!
        let recordingBilis = try await Schedule.query(on: context.application.db)
            .filter(\.$platform, .equal, "bilibili")
            .all()
            .map { $0.$programInfo.id }
        let bilibiliUsers = try await BilibiliUser.query(on: context.application.db)
            .filter(\.$programInfo.$id !~ recordingBilis)
            .all()
        // urlを使って生放送をやっているかどうか確認
        let decoder = JSONDecoder()
        for bilibiliUser in bilibiliUsers {
            let url = URL(string: "https://api.bilibili.com/x/space/acc/info?mid=\(bilibiliUser.userId)")!
            let result = try await URLSession.shared.asyncData(from: url)
            guard let response = result.1 as? HTTPURLResponse,
                response.statusCode == 200,
                let jsonResponse = try? decoder.decode(GetBilibiliSpaceResponse.self, from: result.0),
                jsonResponse.data.live_room.isLiving == true else {
                    continue
            }
            let schedule = Schedule(
                title: jsonResponse.data.live_room.title,
                startDatetime: startDatetime,
                duration: 0,
                platform: "bilibili",
                programInfoId: bilibiliUser.$programInfo.id,
                extraField: jsonResponse.data.live_room.url
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

fileprivate struct GetBilibiliSpaceResponse: Codable {
    let code: Int
    let message: String
    let data: Self.Data

    struct Data: Codable {
        let live_room: LiveRoom
    }

    struct LiveRoom: Codable {
        let roomStatus: Int
        let liveStatus: Int
        let url: String
        let title: String
        let cover: String
        let roomid: Int
        var isLiving: Bool {
            liveStatus == 1
        }
    }
}
