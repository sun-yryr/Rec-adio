@testable import App
import Foundation
import GRDB
import Testing

@Suite
struct JobServiceTests {
    @Test
    func createStoresJobThroughRepository() async throws {
        let repository = JobRepositoryDouble()
        let service = JobService(jobRepo: repository)
        let job = makeJob(jobId: "job-1")

        try await service.create(job)
        let createdJobs = await repository.capturedCreatedJobs()

        #expect(createdJobs.count == 1)
        #expect(createdJobs.first?.jobId == job.jobId)
    }

    @Test
    func createTranslatesConstraintViolationToDuplicateJobError() async {
        let repository = JobRepositoryDouble(createBehavior: .throwConstraintViolation)
        let service = JobService(jobRepo: repository)

        do {
            try await service.create(makeJob(jobId: "job-1"))
            Issue.record("Expected create to throw")
        } catch let error as JobServiceError {
            switch error {
            case let .duplicateJob(reason):
                #expect(reason.contains("UNIQUE constraint failed"))
            default:
                Issue.record("Expected duplicateJob but got: \(error)")
            }
        } catch {
            Issue.record("Expected JobServiceError but got: \(error)")
        }
    }

    @Test
    func createTranslatesUnexpectedErrorToCreateFailed() async {
        let repository = JobRepositoryDouble(createBehavior: .throwUnexpectedError)
        let service = JobService(jobRepo: repository)

        do {
            try await service.create(makeJob(jobId: "job-1"))
            Issue.record("Expected create to throw")
        } catch let error as JobServiceError {
            switch error {
            case let .createFailed(reason):
                #expect(reason.contains("unexpected"))
            default:
                Issue.record("Expected createFailed but got: \(error)")
            }
        } catch {
            Issue.record("Expected JobServiceError but got: \(error)")
        }
    }

    @Test
    func listReturnsJobsFromRepository() async throws {
        let job1 = makeJob(jobId: "job-1")
        let job2 = makeJob(jobId: "job-2")
        let repository = JobRepositoryDouble(findAllBehavior: .succeed([job1, job2]))
        let service = JobService(jobRepo: repository)

        let jobs = try await service.list()

        #expect(jobs.count == 2)
        #expect(jobs.map(\.jobId) == ["job-1", "job-2"])
    }

    @Test
    func listTranslatesUnexpectedErrorToListFailed() async {
        let repository = JobRepositoryDouble(findAllBehavior: .throwUnexpectedError)
        let service = JobService(jobRepo: repository)

        do {
            _ = try await service.list()
            Issue.record("Expected list to throw")
        } catch let error as JobServiceError {
            switch error {
            case let .listFailed(reason):
                #expect(reason.contains("unexpected"))
            default:
                Issue.record("Expected listFailed but got: \(error)")
            }
        } catch {
            Issue.record("Expected JobServiceError but got: \(error)")
        }
    }
}

private func makeJob(jobId: String) -> Job {
    Job(
        jobId: jobId,
        sourceType: "url",
        sourceValue: "https://example.com/\(jobId)",
        title: "Title \(jobId)",
        durationSec: 60,
        scheduledAt: Date(timeIntervalSince1970: 1_700_000_000),
        timezone: "Asia/Tokyo",
        state: .active
    )
}

private actor JobRepositoryDouble: JobRepository {
    enum CreateBehavior {
        case succeed
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
        case .throwConstraintViolation:
            throw DatabaseError(
                resultCode: .SQLITE_CONSTRAINT,
                message: "UNIQUE constraint failed: jobs.job_id"
            )
        case .throwUnexpectedError:
            throw JobServiceTestError.unexpected
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
            throw JobServiceTestError.unexpected
        }
    }

    func delete(jobId _: String) async throws -> Bool {
        false
    }

    func capturedCreatedJobs() -> [Job] {
        createdJobs
    }
}

private enum JobServiceTestError: Error {
    case unexpected
}
