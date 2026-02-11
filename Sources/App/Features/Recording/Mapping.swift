import SwiftProtobuf

func convertGrpcJob(job: Job) -> Recoto_Recording_V1_Job {
  var grpcJob = Recoto_Recording_V1_Job()
  grpcJob.jobID = job.id
  grpcJob.duration = Google_Protobuf_Duration(seconds: job.durationSec, nanos: 0)
  grpcJob.title = job.title
  grpcJob.sourceType = {
    switch job.sourceType {
    case "url":
      return .url
    default:
      return .unspecified
    }
  }()
  grpcJob.sourceValue = job.sourceValue
  grpcJob.scheduledAt = Google_Protobuf_Timestamp(date: job.scheduledAt)
  grpcJob.timezone = job.timezone
  grpcJob.state = {
    switch job.state {
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
