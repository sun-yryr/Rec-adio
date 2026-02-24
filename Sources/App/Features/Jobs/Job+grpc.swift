import Dependencies
import SwiftProtobuf

func makeJob(from request: Recoto_Recording_V1_CreateJobRequest) -> Job {
    @Dependency(\.uuid) var uuid

    return Job(
        jobId: uuid().uuidString, sourceType: "url", sourceValue: request.url,
        title: request.title,
        durationSec: request.duration.seconds, scheduledAt: request.scheduledAt.date,
        timezone: request.timezone, state: .active
    )
}

extension Job {
    func toGrpcJob() -> Recoto_Recording_V1_Job {
        var grpcJob = Recoto_Recording_V1_Job()
        grpcJob.jobID = id
        grpcJob.duration = Google_Protobuf_Duration(seconds: durationSec, nanos: 0)
        grpcJob.title = title
        grpcJob.sourceType = {
            switch self.sourceType {
            case "url":
                return .url
            default:
                return .unspecified
            }
        }()
        grpcJob.sourceValue = sourceValue
        grpcJob.scheduledAt = Google_Protobuf_Timestamp(date: scheduledAt)
        grpcJob.timezone = timezone
        grpcJob.state = {
            switch self.state {
            case .active:
                return .active
            case .deleted:
                return .deleted
            case .paused:
                return .paused
            }
        }()

        return grpcJob
    }
}
