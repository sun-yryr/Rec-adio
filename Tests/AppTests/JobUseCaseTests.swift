@testable import App
import Foundation
import GRDB
import Testing

@Suite
struct JobUseCaseTests {
    @Test
    func createStoresJobThroughRepository() async throws {
        let repository = JobRepositoryDouble()
        let useCase = DefaultJobUseCase(jobRepository: repository)
        let job = makeJob(jobId: "job-1")

        try await useCase.create(job)
        let stored = try await repository.find(jobId: job.jobId)

        #expect(stored?.jobId == job.jobId)
    }

    @Test
    func createTranslatesConstraintViolationToDuplicateJobError() async {
        let repository = JobRepositoryDouble(createBehavior: .throwConstraintViolation)
        let useCase = DefaultJobUseCase(jobRepository: repository)

        do {
            try await useCase.create(makeJob(jobId: "job-1"))
            Issue.record("Expected create to throw")
        } catch let error as JobUseCaseError {
            switch error {
            case let .duplicateJob(reason):
                #expect(reason.contains("UNIQUE constraint failed"))
            default:
                Issue.record("Expected duplicateJob but got: \(error)")
            }
        } catch {
            Issue.record("Expected JobUseCaseError but got: \(error)")
        }
    }

    @Test
    func getReturnsStoredJob() async throws {
        let repository = JobRepositoryDouble(initialJobs: [makeJob(jobId: "job-1")])
        let useCase = DefaultJobUseCase(jobRepository: repository)

        let job = try await useCase.get(jobId: "job-1")

        #expect(job.jobId == "job-1")
    }

    @Test
    func getThrowsNotFoundForUnknownJob() async {
        let repository = JobRepositoryDouble()
        let useCase = DefaultJobUseCase(jobRepository: repository)

        do {
            _ = try await useCase.get(jobId: "missing")
            Issue.record("Expected get to throw")
        } catch let error as JobUseCaseError {
            switch error {
            case let .jobNotFound(jobId):
                #expect(jobId == "missing")
            default:
                Issue.record("Expected jobNotFound but got: \(error)")
            }
        } catch {
            Issue.record("Expected JobUseCaseError but got: \(error)")
        }
    }

    @Test
    func listReturnsJobsFromRepository() async throws {
        let job1 = makeJob(jobId: "job-1")
        let job2 = makeJob(jobId: "job-2")
        let repository = JobRepositoryDouble(initialJobs: [job1, job2])
        let useCase = DefaultJobUseCase(jobRepository: repository)

        let jobs = try await useCase.list()

        #expect(jobs.count == 2)
        #expect(Set(jobs.map(\.jobId)) == Set(["job-1", "job-2"]))
    }

    @Test
    func updatePersistsChangesAndReturnsUpdatedJob() async throws {
        let initial = makeJob(jobId: "job-1")
        let repository = JobRepositoryDouble(initialJobs: [initial])
        let useCase = DefaultJobUseCase(jobRepository: repository)
        let updated = Job(
            jobId: initial.jobId,
            sourceType: initial.sourceType,
            sourceValue: initial.sourceValue,
            title: "Updated title",
            durationSec: 120,
            scheduledAt: Date(timeIntervalSince1970: 1_700_000_123),
            timezone: "UTC",
            state: .paused
        )

        let result = try await useCase.update(updated)
        let stored = try await repository.find(jobId: initial.jobId)

        #expect(result.jobId == updated.jobId)
        #expect(result.title == "Updated title")
        #expect(stored?.state == .paused)
    }

    @Test
    func updateThrowsNotFoundForUnknownJob() async {
        let repository = JobRepositoryDouble(updateBehavior: .returnFalse)
        let useCase = DefaultJobUseCase(jobRepository: repository)

        do {
            _ = try await useCase.update(makeJob(jobId: "missing"))
            Issue.record("Expected update to throw")
        } catch let error as JobUseCaseError {
            switch error {
            case let .jobNotFound(jobId):
                #expect(jobId == "missing")
            default:
                Issue.record("Expected jobNotFound but got: \(error)")
            }
        } catch {
            Issue.record("Expected JobUseCaseError but got: \(error)")
        }
    }

    @Test
    func deleteRemovesExistingJob() async throws {
        let initial = makeJob(jobId: "job-1")
        let repository = JobRepositoryDouble(initialJobs: [initial])
        let useCase = DefaultJobUseCase(jobRepository: repository)

        try await useCase.delete(jobId: initial.jobId)
        let found = try await repository.find(jobId: initial.jobId)

        #expect(found == nil)
    }

    @Test
    func deleteThrowsNotFoundForUnknownJob() async {
        let repository = JobRepositoryDouble(deleteBehavior: .returnFalse)
        let useCase = DefaultJobUseCase(jobRepository: repository)

        do {
            try await useCase.delete(jobId: "missing")
            Issue.record("Expected delete to throw")
        } catch let error as JobUseCaseError {
            switch error {
            case let .jobNotFound(jobId):
                #expect(jobId == "missing")
            default:
                Issue.record("Expected jobNotFound but got: \(error)")
            }
        } catch {
            Issue.record("Expected JobUseCaseError but got: \(error)")
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

    enum UpdateBehavior {
        case useStorage
        case returnFalse
        case throwUnexpectedError
    }

    enum DeleteBehavior {
        case useStorage
        case returnFalse
        case throwUnexpectedError
    }

    private let createBehavior: CreateBehavior
    private let updateBehavior: UpdateBehavior
    private let deleteBehavior: DeleteBehavior
    private var jobs: [String: Job]

    init(
        createBehavior: CreateBehavior = .succeed,
        updateBehavior: UpdateBehavior = .useStorage,
        deleteBehavior: DeleteBehavior = .useStorage,
        initialJobs: [Job] = []
    ) {
        self.createBehavior = createBehavior
        self.updateBehavior = updateBehavior
        self.deleteBehavior = deleteBehavior
        self.jobs = Dictionary(uniqueKeysWithValues: initialJobs.map { ($0.jobId, $0) })
    }

    func create(_ job: Job) async throws {
        switch createBehavior {
        case .succeed:
            jobs[job.jobId] = job
        case .throwConstraintViolation:
            throw DatabaseError(
                resultCode: .SQLITE_CONSTRAINT,
                message: "UNIQUE constraint failed: jobs.job_id"
            )
        case .throwUnexpectedError:
            throw JobUseCaseTestError.unexpected
        }
    }

    func find(jobId: String) async throws -> Job? {
        jobs[jobId]
    }

    func findAll() async throws -> [Job] {
        jobs.values.sorted { $0.scheduledAt < $1.scheduledAt }
    }

    func update(_ job: Job) async throws -> Bool {
        switch updateBehavior {
        case .useStorage:
            guard jobs[job.jobId] != nil else {
                return false
            }
            jobs[job.jobId] = job
            return true
        case .returnFalse:
            return false
        case .throwUnexpectedError:
            throw JobUseCaseTestError.unexpected
        }
    }

    func delete(jobId: String) async throws -> Bool {
        switch deleteBehavior {
        case .useStorage:
            guard jobs[jobId] != nil else {
                return false
            }
            jobs[jobId] = nil
            return true
        case .returnFalse:
            return false
        case .throwUnexpectedError:
            throw JobUseCaseTestError.unexpected
        }
    }
}

private enum JobUseCaseTestError: Error {
    case unexpected
}
