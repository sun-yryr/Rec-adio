import Foundation
import GRPCCore
import Logging

struct RecordingService: Recoto_Recording_V1_RecordingService.SimpleServiceProtocol {
    private let logger: Logger
    private let jobUseCase: any JobUseCase

    init(
        jobUseCase: any JobUseCase,
        logger: Logger = Logger(label: "Recoto.RecordingService")
    ) {
        self.jobUseCase = jobUseCase
        self.logger = logger
    }

    init(jobRepo: any JobRepository, logger: Logger = Logger(label: "Recoto.RecordingService")) {
        self.init(jobUseCase: DefaultJobUseCase(jobRepository: jobRepo), logger: logger)
    }

    func createJob(request: Recoto_Recording_V1_CreateJobRequest, context _: ServerContext)
        async throws
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
        try await create(job, logger: logger)

        var response = Recoto_Recording_V1_CreateJobResponse()
        response.job = job.toGrpcJob()

        logger.info(
            "recording.create_job.finished",
            metadata: ["job_id": .string(job.jobId)]
        )

        return response
    }

    private func create(_ job: Job, logger: Logger) async throws {
        do {
            try await jobUseCase.create(job)
        } catch let error as JobUseCaseError {
            switch error {
            case let .duplicateJob(reason):
                logger.error(
                    "recording.create_job.failed",
                    metadata: [
                        "job_id": .string(job.jobId),
                        "grpc_status": .string("alreadyExists"),
                        "reason": .string(reason),
                    ]
                )
                throw RPCError(
                    code: .alreadyExists,
                    message: "recording job already exists",
                    cause: error
                )
            default:
                logger.error(
                    "recording.create_job.failed",
                    metadata: [
                        "job_id": .string(job.jobId),
                        "grpc_status": .string("internalError"),
                        "reason": .string(error.reason),
                    ]
                )
                throw RPCError(
                    code: .internalError,
                    message: "failed to create recording job",
                    cause: error
                )
            }
        }
    }

    func listJobs(request _: Recoto_Recording_V1_ListJobsRequest, context _: ServerContext)
        async throws
        -> Recoto_Recording_V1_ListJobsResponse
    {
        let logger = self.logger.rpc()

        do {
            let jobs = try await jobUseCase.list()

            var response = Recoto_Recording_V1_ListJobsResponse()
            response.jobs = jobs.map { job in job.toGrpcJob() }

            return response
        } catch let error as JobUseCaseError {
            logger.error(
                "recording.list_jobs.failed",
                metadata: [
                    "grpc_status": .string("internalError"),
                    "reason": .string(error.reason),
                ]
            )
            throw RPCError(
                code: .internalError, message: "failed to list jobs", cause: error
            )
        } catch {
            logger.error(
                "recording.list_jobs.failed",
                metadata: [
                    "grpc_status": .string("internalError"),
                    "reason": .string(String(describing: error)),
                ]
            )
            throw RPCError(
                code: .internalError, message: "failed to list jobs", cause: error
            )
        }
    }

    func updateJob(request _: Recoto_Recording_V1_UpdateJobRequest, context _: ServerContext)
        async throws
        -> Recoto_Recording_V1_UpdateJobResponse
    {
        try unimplementedRPC("updateJob")
    }

    func pauseJob(request _: Recoto_Recording_V1_PauseJobRequest, context _: ServerContext)
        async throws
        -> Recoto_Recording_V1_PauseJobResponse
    {
        try unimplementedRPC("pauseJob")
    }

    func resumeJob(request _: Recoto_Recording_V1_ResumeJobRequest, context _: ServerContext)
        async throws
        -> Recoto_Recording_V1_ResumeJobResponse
    {
        try unimplementedRPC("resumeJob")
    }

    func deleteJob(request _: Recoto_Recording_V1_DeleteJobRequest, context _: ServerContext)
        async throws
        -> Recoto_Recording_V1_DeleteJobResponse
    {
        try unimplementedRPC("deleteJob")
    }

    func getJob(request _: Recoto_Recording_V1_GetJobRequest, context _: ServerContext)
        async throws
        -> Recoto_Recording_V1_GetJobResponse
    {
        try unimplementedRPC("getJob")
    }

    func cancelRun(request _: Recoto_Recording_V1_CancelRunRequest, context _: ServerContext)
        async throws
        -> Recoto_Recording_V1_CancelRunResponse
    {
        try unimplementedRPC("cancelRun")
    }

    func getRun(request _: Recoto_Recording_V1_GetRunRequest, context _: ServerContext)
        async throws
        -> Recoto_Recording_V1_GetRunResponse
    {
        try unimplementedRPC("getRun")
    }

    func listRuns(request _: Recoto_Recording_V1_ListRunsRequest, context _: ServerContext)
        async throws
        -> Recoto_Recording_V1_ListRunsResponse
    {
        try unimplementedRPC("listRuns")
    }

    private func unimplementedRPC<Response>(_ method: String) throws -> Response {
        throw RPCError(
            code: .unimplemented,
            message: "\(method) is not implemented yet"
        )
    }
}
