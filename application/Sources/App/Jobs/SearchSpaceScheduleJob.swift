import Vapor
import Queues
import FluentKit
#if canImport(FoundationNetworking)
import FoundationNetworking
#endif

struct SearchSpaceScheduleJob: AsyncScheduledJob {
    let bearerToken: String 
    let calendar = Calendar(identifier: .gregorian)
    let decoder = JSONDecoder()

    func run(context: QueueContext) async throws {
        context.logger.info("Search space job - start. \(Date())")
        let startDatetime = self.calendar.date(byAdding: .minute, value: 1, to: Date())!
        // 録音キューに積まれてるユーザーを除く処理
        let recordingSpaces = try await Schedule.query(on: context.application.db)
            .filter(\.$platform, .equal, "space")
            .all()
            .map { $0.$programInfo.id }
        let twitterUsers = try await TwitterUser.query(on: context.application.db)
            .filter(\.$programInfo.$id !~ recordingSpaces)
            .all()
        let userIds = twitterUsers.map { $0.userId }

        let url = URL(string: "https://api.twitter.com/2/spaces/by/creator_ids?user_ids=\(userIds.joined(separator: ","))&space.fields=creator_id,title")!
        var request = URLRequest(url: url)
        request.addValue("Bearer \(self.bearerToken)", forHTTPHeaderField: "Authorization")
        let result = try await URLSession.shared.asyncData(from: request)
        guard let response = result.1 as? HTTPURLResponse,
            response.statusCode == 200 else {
            context.logger.notice("twitterAPIへのアクセスに失敗しました", metadata: ["response": .string(String(data: result.0, encoding: .utf8) ?? "empty")])
            return
        }
        let jsonResponse = try! self.decoder.decode(GetSpacesByCreatorResponse.self, from: result.0)
        let activeSpaces = jsonResponse.data.filter { $0.state == "live" }
        for activeSpace in activeSpaces {
            guard let twitterUser = twitterUsers.filter({ $0.userId == activeSpace.creator_id }).first else {
                context.logger.notice("userIdとcreatorIdの紐付けができませんでした", metadata: ["user_id": .string(activeSpace.creator_id), "title": .string(activeSpace.title)])
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
            } catch {
                // エラーが発生しても後続は通す
                context.logger.error("Scheduleの登録に失敗しました。", metadata: ["job": .string(self.name)])
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
