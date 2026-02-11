import Foundation
import GRPCCore
import Logging

struct RecordingService: Recoto_Recording_V1_RecordingService.SimpleServiceProtocol {
  private let logger: Logger

  init(logger: Logger = Logger(label: "Recoto.RecordingService")) {
    self.logger = logger
  }

  func createJob(request: Recoto_Recording_V1_CreateJobRequest, context: ServerContext) async throws
    -> Recoto_Recording_V1_CreateJobResponse
  {
    let logger = self.logger.rpc()
    logger.info(
      "recording.create_job.started",
      metadata: [
        "source_type": .string("url"),
        "source_value": .string(request.url),
        "title": .string(request.title),
      ]
    )

    let job = Job(
      jobId: "testid", sourceType: "url", sourceValue: request.url, title: request.title,
      durationSec: 60, scheduledAt: .now, timezone: "Asia/Tokyo", state: .active)

    var response = Recoto_Recording_V1_CreateJobResponse()
    response.job = convertGrpcJob(job: job)

    logger.info(
      "recording.create_job.finished",
      metadata: ["job_id": .string(job.jobId)]
    )

    return response
  }
}
