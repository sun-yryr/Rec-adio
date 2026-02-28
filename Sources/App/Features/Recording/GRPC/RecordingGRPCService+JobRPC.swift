import GRPCCore
import Logging

extension RecordingGRPCService {
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
            throw mapJobUseCaseError(error, operation: .list, logger: logger)
        } catch {
            throw mapUnexpectedJobError(error, operation: .list, logger: logger)
        }
    }

    func updateJob(request: Recoto_Recording_V1_UpdateJobRequest, context _: ServerContext)
        async throws
        -> Recoto_Recording_V1_UpdateJobResponse
    {
        let logger = self.logger.rpc()
        let updatable = makeUpdatable(from: request)

        try validateUpdateJobRequest(updatable, jobId: request.jobID, logger: logger)

        do {
            let job = try await jobUseCase.update(jobId: request.jobID, updatable: updatable)

            var response = Recoto_Recording_V1_UpdateJobResponse()
            response.job = job.toGrpcJob()
            return response
        } catch let error as JobUseCaseError {
            throw mapJobUseCaseError(
                error,
                operation: .update(jobId: request.jobID),
                logger: logger
            )
        } catch {
            throw mapUnexpectedJobError(
                error,
                operation: .update(jobId: request.jobID),
                logger: logger
            )
        }
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
}

private extension RecordingGRPCService {
    func create(_ job: Job, logger: Logger) async throws {
        do {
            try await jobUseCase.create(job)
        } catch let error as JobUseCaseError {
            throw mapJobUseCaseError(
                error,
                operation: .create(jobId: job.jobId),
                logger: logger
            )
        } catch {
            throw mapUnexpectedJobError(
                error,
                operation: .create(jobId: job.jobId),
                logger: logger
            )
        }
    }

    func validateUpdateJobRequest(_ updatable: Job.Updatable, jobId: String, logger: Logger) throws {
        guard !updatable.isEmpty else {
            logJobOperationFailure(
                operation: .update(jobId: jobId),
                grpcStatus: "invalidArgument",
                reason: "no updatable fields provided",
                logger: logger
            )
            throw RPCError(
                code: .invalidArgument,
                message: "at least one updatable field is required"
            )
        }
    }

    func mapJobUseCaseError(
        _ error: JobUseCaseError,
        operation: JobRPCOperation,
        logger: Logger
    ) -> RPCError {
        switch (operation, error) {
        case (.create, let .duplicateJob(reason)):
            logJobOperationFailure(
                operation: operation,
                grpcStatus: "alreadyExists",
                reason: reason,
                logger: logger
            )
            return RPCError(
                code: .alreadyExists,
                message: "recording job already exists",
                cause: error
            )
        case (.update, .jobNotFound):
            logJobOperationFailure(
                operation: operation,
                grpcStatus: "notFound",
                reason: error.reason,
                logger: logger
            )
            return RPCError(
                code: .notFound,
                message: "recording job not found",
                cause: error
            )
        default:
            logJobOperationFailure(
                operation: operation,
                grpcStatus: "internalError",
                reason: error.reason,
                logger: logger
            )
            return RPCError(
                code: .internalError,
                message: operation.failureMessage,
                cause: error
            )
        }
    }

    func mapUnexpectedJobError(
        _ error: Error,
        operation: JobRPCOperation,
        logger: Logger
    ) -> RPCError {
        logJobOperationFailure(
            operation: operation,
            grpcStatus: "internalError",
            reason: String(describing: error),
            logger: logger
        )
        return RPCError(
            code: .internalError,
            message: operation.failureMessage,
            cause: error
        )
    }

    func logJobOperationFailure(
        operation: JobRPCOperation,
        grpcStatus: String,
        reason: String,
        logger: Logger
    ) {
        var metadata: Logger.Metadata = [
            "grpc_status": .string(grpcStatus),
            "reason": .string(reason),
        ]
        if let jobId = operation.jobId {
            metadata["job_id"] = .string(jobId)
        }
        logger.error("\(operation.failureLogEvent)", metadata: metadata)
    }
}

private enum JobRPCOperation {
    case create(jobId: String)
    case list
    case update(jobId: String)

    var jobId: String? {
        switch self {
        case let .create(jobId):
            return jobId
        case .list:
            return nil
        case let .update(jobId):
            return jobId
        }
    }

    var failureLogEvent: String {
        switch self {
        case .create:
            return "recording.create_job.failed"
        case .list:
            return "recording.list_jobs.failed"
        case .update:
            return "recording.update_job.failed"
        }
    }

    var failureMessage: String {
        switch self {
        case .create:
            return "failed to create recording job"
        case .list:
            return "failed to list jobs"
        case .update:
            return "failed to update recording job"
        }
    }
}
