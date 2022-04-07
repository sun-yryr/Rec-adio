import Vapor
import Queues

struct RecordingJob: AsyncJob {
    typealias Payload = Schedule
    
    func dequeue(_ context: QueueContext, _ payload: Payload) async throws {
        context.logger.debug("RecordingJob - dequeue [\(payload.title), \(Date())]")
        guard let recorder = context.application.radio.recorders.filter({ $0.supportPlatform == payload.platform }).first else {
            throw RecordingError.NotFoundPlatform(payload.platform)
        }

        switch recorder.run(application: context.application, schedule: payload) {
        case .success(let record):
            do {
                try await record.create(on: context.application.db)
                try await payload.delete(force: true, on: context.application.db)
                context.logger.info("recording complete. title[\(record.title)]")
            } catch {
                context.logger.error("録音後のDB書き込みに失敗しました。schedule_id[\(payload.id!)]")
            }
        case .failure(let error):
            switch error {
                case .processError(let message):
                    context.logger.error(.init(stringLiteral: message))
                case .defaultError:
                    context.logger.error("録音中に不明なエラーが発生しました。 schedule_id[\(payload.id!)]")
                case .NotFoundPlatform:
                    fatalError()
            }
            context.logger.error(.init(stringLiteral: error.localizedDescription))
        }
    }
    
    func error(_ context: QueueContext, _ error: Error, _ payload: Schedule) async throws {
        context.logger.error("Failure Recording. schedule_id = \(payload.id!)")
    }
}
