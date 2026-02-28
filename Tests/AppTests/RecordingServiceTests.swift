@testable import App
import Dependencies
import Foundation
import GRDB
import GRPCCore
import SwiftProtobuf
import Testing

@Suite
struct RecordingServiceTests {
    @Test
    func createJobPersistsMappedJobAndReturnsMappedResponse() async throws {
        let fixedUUID = try #require(UUID(uuidString: "0AF76519-16CD-43DD-8448-EB211C80319C"))
        let scheduledAt = Date(timeIntervalSince1970: 1_700_123_456)
        let request = makeCreateJobRequest(
            url: "https://example.com/live.m3u8",
            title: "Morning Show",
            durationSec: 180,
            scheduledAt: scheduledAt,
            timezone: "Asia/Tokyo"
        )
        let repository = JobRepositoryStub()
        let service = RecordingService(jobRepo: repository)

        let response = try await withDependencies {
            $0.uuid = .constant(fixedUUID)
            $0.date = .constant(Date(timeIntervalSince1970: 1_700_000_000))
        } operation: {
            try await service.createJob(
                request: request,
                context: makeRecordingServerContext(method: "CreateJob")
            )
        }

        let createdJobs = await repository.capturedCreatedJobs()
        #expect(createdJobs.count == 1)

        guard let createdJob = createdJobs.first else {
            Issue.record("Expected one created job")
            return
        }

        #expect(createdJob.jobId == fixedUUID.uuidString)
        #expect(createdJob.sourceType == "url")
        #expect(createdJob.sourceValue == request.url)
        #expect(createdJob.title == request.title)
        #expect(createdJob.durationSec == request.duration.seconds)
        #expect(createdJob.scheduledAt == scheduledAt)
        #expect(createdJob.timezone == request.timezone)
        #expect(createdJob.state == .active)

        #expect(response.job.jobID == createdJob.jobId)
        #expect(response.job.sourceType == .url)
        #expect(response.job.sourceValue == request.url)
        #expect(response.job.title == request.title)
        #expect(response.job.duration.seconds == request.duration.seconds)
        #expect(response.job.scheduledAt.date == scheduledAt)
        #expect(response.job.timezone == request.timezone)
        #expect(response.job.state == .active)
    }

    @Test
    func createJobTranslatesRepositoryRPCErrorToInternalError() async throws {
        let expected = RPCError(code: .invalidArgument, message: "invalid create request")
        let repository = JobRepositoryStub(createBehavior: .throwRPCError(expected))
        let service = RecordingService(jobRepo: repository)
        let fixedUUID = try #require(UUID(uuidString: "00000000-0000-0000-0000-000000000001"))

        do {
            _ = try await withDependencies {
                $0.uuid = .constant(fixedUUID)
                $0.date = .constant(Date(timeIntervalSince1970: 1_700_000_000))
            } operation: {
                try await service.createJob(
                    request: makeCreateJobRequest(),
                    context: makeRecordingServerContext(method: "CreateJob")
                )
            }
            Issue.record("Expected createJob to throw")
        } catch let error as RPCError {
            #expect(error.code == .internalError)
            #expect(error.message == "failed to create recording job")
            #expect(error.cause is JobUseCaseError)
            let cause = error.cause as? JobUseCaseError
            switch cause {
            case let .createFailed(reason):
                #expect(reason.contains(expected.message))
            default:
                Issue.record("Expected createFailed cause but got: \(String(describing: cause))")
            }
        } catch {
            Issue.record("Expected RPCError but got: \(error)")
        }
    }

    @Test
    func createJobTranslatesConstraintViolationToAlreadyExists() async throws {
        let repository = JobRepositoryStub(createBehavior: .throwConstraintViolation)
        let service = RecordingService(jobRepo: repository)
        let fixedUUID = try #require(UUID(uuidString: "00000000-0000-0000-0000-000000000002"))

        do {
            _ = try await withDependencies {
                $0.uuid = .constant(fixedUUID)
                $0.date = .constant(Date(timeIntervalSince1970: 1_700_000_000))
            } operation: {
                try await service.createJob(
                    request: makeCreateJobRequest(),
                    context: makeRecordingServerContext(method: "CreateJob")
                )
            }
            Issue.record("Expected createJob to throw")
        } catch let error as RPCError {
            #expect(error.code == .alreadyExists)
            #expect(error.message == "recording job already exists")
            #expect(error.cause is JobUseCaseError)
            let cause = error.cause as? JobUseCaseError
            switch cause {
            case let .duplicateJob(reason):
                #expect(reason.contains("UNIQUE constraint failed"))
            default:
                Issue.record("Expected duplicateJob cause but got: \(String(describing: cause))")
            }
        } catch {
            Issue.record("Expected RPCError but got: \(error)")
        }
    }

    @Test
    func createJobTranslatesUnexpectedErrorToInternalError() async throws {
        let repository = JobRepositoryStub(createBehavior: .throwUnexpectedError)
        let service = RecordingService(jobRepo: repository)
        let fixedUUID = try #require(UUID(uuidString: "00000000-0000-0000-0000-000000000003"))

        do {
            _ = try await withDependencies {
                $0.uuid = .constant(fixedUUID)
                $0.date = .constant(Date(timeIntervalSince1970: 1_700_000_000))
            } operation: {
                try await service.createJob(
                    request: makeCreateJobRequest(),
                    context: makeRecordingServerContext(method: "CreateJob")
                )
            }
            Issue.record("Expected createJob to throw")
        } catch let error as RPCError {
            #expect(error.code == .internalError)
            #expect(error.message == "failed to create recording job")
            #expect(error.cause is JobUseCaseError)
            let cause = error.cause as? JobUseCaseError
            switch cause {
            case let .createFailed(reason):
                #expect(reason.contains("unexpected"))
            default:
                Issue.record("Expected createFailed cause but got: \(String(describing: cause))")
            }
        } catch {
            Issue.record("Expected RPCError but got: \(error)")
        }
    }

    @Test
    func listJobsReturnsMappedJobs() async throws {
        let job1 = makeJob(
            jobId: "job-1",
            sourceType: "url",
            state: .active,
            scheduledAt: Date(timeIntervalSince1970: 1_700_000_000)
        )
        let job2 = makeJob(
            jobId: "job-2",
            sourceType: "unsupported",
            state: .paused,
            scheduledAt: Date(timeIntervalSince1970: 1_700_000_060)
        )
        let job3 = makeJob(
            jobId: "job-3",
            sourceType: "url",
            state: .deleted,
            scheduledAt: Date(timeIntervalSince1970: 1_700_000_120)
        )
        let repository = JobRepositoryStub(findAllBehavior: .succeed([job1, job2, job3]))
        let service = RecordingService(jobRepo: repository)

        let response = try await service.listJobs(
            request: Recoto_Recording_V1_ListJobsRequest(),
            context: makeRecordingServerContext(method: "ListJobs")
        )

        #expect(response.jobs.count == 3)
        #expect(response.jobs[0].jobID == job1.jobId)
        #expect(response.jobs[0].sourceType == .url)
        #expect(response.jobs[0].state == .active)
        #expect(response.jobs[1].jobID == job2.jobId)
        #expect(response.jobs[1].sourceType == .unspecified)
        #expect(response.jobs[1].state == .paused)
        #expect(response.jobs[2].jobID == job3.jobId)
        #expect(response.jobs[2].sourceType == .url)
        #expect(response.jobs[2].state == .deleted)
    }

    @Test
    func listJobsTranslatesUnexpectedErrorToInternalError() async {
        let repository = JobRepositoryStub(findAllBehavior: .throwUnexpectedError)
        let service = RecordingService(jobRepo: repository)

        do {
            _ = try await service.listJobs(
                request: Recoto_Recording_V1_ListJobsRequest(),
                context: makeRecordingServerContext(method: "ListJobs")
            )
            Issue.record("Expected listJobs to throw")
        } catch let error as RPCError {
            #expect(error.code == .internalError)
            #expect(error.message == "failed to list jobs")
            #expect(error.cause is JobUseCaseError)
            let cause = error.cause as? JobUseCaseError
            switch cause {
            case let .listFailed(reason):
                #expect(reason.contains("unexpected"))
            default:
                Issue.record("Expected listFailed cause but got: \(String(describing: cause))")
            }
        } catch {
            Issue.record("Expected RPCError but got: \(error)")
        }
    }
}

private func makeCreateJobRequest(
    url: String = "https://example.com/live.m3u8",
    title: String = "Morning Show",
    durationSec: Int64 = 180,
    scheduledAt: Date = Date(timeIntervalSince1970: 1_700_123_456),
    timezone: String = "Asia/Tokyo"
) -> Recoto_Recording_V1_CreateJobRequest {
    var request = Recoto_Recording_V1_CreateJobRequest()
    request.url = url
    request.title = title
    request.duration = Google_Protobuf_Duration(seconds: durationSec, nanos: 0)
    request.scheduledAt = Google_Protobuf_Timestamp(date: scheduledAt)
    request.timezone = timezone
    return request
}

private func makeJob(
    jobId: String,
    sourceType: String,
    state: JobState,
    scheduledAt: Date
) -> Job {
    Job(
        jobId: jobId,
        sourceType: sourceType,
        sourceValue: "https://example.com/\(jobId)",
        title: "Title \(jobId)",
        durationSec: 60,
        scheduledAt: scheduledAt,
        timezone: "Asia/Tokyo",
        state: state
    )
}

private func makeRecordingServerContext(method: String) -> ServerContext {
    ServerContext(
        descriptor: MethodDescriptor(
            fullyQualifiedService: "recoto.recording.v1.RecordingService",
            method: method
        ),
        remotePeer: "ipv4:127.0.0.1:50000",
        localPeer: "ipv4:127.0.0.1:50051",
        cancellation: .init()
    )
}

private actor JobRepositoryStub: JobRepository {
    enum CreateBehavior {
        case succeed
        case throwRPCError(RPCError)
        case throwConstraintViolation
        case throwUnexpectedError
    }

    enum FindAllBehavior {
        case succeed([Job])
        case throwUnexpectedError
    }

    private let createBehavior: CreateBehavior
    private let findAllBehavior: FindAllBehavior
    private var createdJobs: [Job] = []

    init(
        createBehavior: CreateBehavior = .succeed,
        findAllBehavior: FindAllBehavior = .succeed([])
    ) {
        self.createBehavior = createBehavior
        self.findAllBehavior = findAllBehavior
    }

    func create(_ job: Job) async throws {
        switch createBehavior {
        case .succeed:
            createdJobs.append(job)
        case let .throwRPCError(error):
            throw error
        case .throwConstraintViolation:
            throw DatabaseError(
                resultCode: .SQLITE_CONSTRAINT,
                message: "UNIQUE constraint failed: jobs.job_id"
            )
        case .throwUnexpectedError:
            throw RecordingRepositoryTestError.unexpected
        }
    }

    func find(jobId: String) async throws -> Job? {
        createdJobs.first(where: { $0.jobId == jobId })
    }

    func findAll() async throws -> [Job] {
        switch findAllBehavior {
        case let .succeed(jobs):
            return jobs
        case .throwUnexpectedError:
            throw RecordingRepositoryTestError.unexpected
        }
    }

    func update(_ job: Job) async throws -> Bool {
        guard let index = createdJobs.firstIndex(where: { $0.jobId == job.jobId }) else {
            return false
        }
        createdJobs[index] = job
        return true
    }

    func delete(jobId _: String) async throws -> Bool {
        false
    }

    func capturedCreatedJobs() -> [Job] {
        createdJobs
    }
}

private enum RecordingRepositoryTestError: Error {
    case unexpected
}
