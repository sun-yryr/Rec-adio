import Foundation
import GRDB
import GRPCCore
import Logging
import UUIDV7

struct RecordingService: Recoto_Recording_V1_RecordingService.SimpleServiceProtocol {
    private let logger: Logger
    private let jobRepo: any JobRepository

    init(jobRepo: any JobRepository, logger: Logger = Logger(label: "Recoto.RecordingService")) {
        self.jobRepo = jobRepo
        self.logger = logger
    }

    func createJob(request: Recoto_Recording_V1_CreateJobRequest, context _: ServerContext) async throws
        -> Recoto_Recording_V1_CreateJobResponse
    {
        let logger = self.logger.rpc()
        logger.debug(
            "recording.create_job.started",
            metadata: [
                "source_type": .string("url"),
                "source_value": .string(request.url),
                "title": .string(request.title),
            ]
        )

        let job = makeJob(from: request)
        try await persistJob(job, logger: logger)

        var response = Recoto_Recording_V1_CreateJobResponse()
        response.job = convertGrpcJob(job: job)

        logger.info(
            "recording.create_job.finished",
            metadata: ["job_id": .string(job.jobId)]
        )

        return response
    }

    private func makeJob(from request: Recoto_Recording_V1_CreateJobRequest) -> Job {
        Job(
            jobId: UUIDV7().uuidString, sourceType: "url", sourceValue: request.url, title: request.title,
            durationSec: request.duration.seconds, scheduledAt: request.scheduledAt.date,
            timezone: request.timezone, state: .active
        )
    }

    private func persistJob(_ job: Job, logger: Logger) async throws {
        do {
            try await jobRepo.create(job)
        } catch let error as RPCError {
            throw error
        } catch let error as DatabaseError where error.resultCode == .SQLITE_CONSTRAINT {
            logger.error(
                "recording.create_job.failed",
                metadata: [
                    "job_id": .string(job.jobId),
                    "grpc_status": .string("alreadyExists"),
                    "reason": .string(error.message ?? "constraint violation"),
                ]
            )
            throw RPCError(
                code: .alreadyExists,
                message: "recording job already exists",
                cause: error
            )
        } catch {
            logger.error(
                "recording.create_job.failed",
                metadata: [
                    "job_id": .string(job.jobId),
                    "grpc_status": .string("internalError"),
                    "reason": .string(String(describing: error)),
                ]
            )
            throw RPCError(
                code: .internalError,
                message: "failed to create recording job",
                cause: error
            )
        }
    }

    func listJobs(request _: Recoto_Recording_V1_ListJobsRequest, context _: ServerContext) async throws
        -> Recoto_Recording_V1_ListJobsResponse
    {
        let logger = self.logger.rpc()

        do {
            let jobs = try await jobRepo.findAll()

            var response = Recoto_Recording_V1_ListJobsResponse()
            response.jobs = jobs.map(convertGrpcJob)

            return response
        } catch {
            logger.error(
                "failed list jobs",
                metadata: [
                    "error": .string(String(describing: error)),
                ]
            )
            throw RPCError(
                code: .internalError, message: "failed to list jobs", cause: error
            )
        }
    }
}
