@testable import App
import Foundation
import GRDB
import Testing

@Suite
struct JobRepositoryTests {
    @Test
    func createThenFindReturnsStoredJob() async throws {
        let repository = try makeTestRepository()
        let job = makeJob(jobId: "job-1")

        try await repository.create(job)
        let found = try await repository.find(jobId: job.jobId)

        #expect(found != nil)
        #expect(found?.jobId == job.jobId)
        #expect(found?.sourceType == job.sourceType)
        #expect(found?.sourceValue == job.sourceValue)
        #expect(found?.title == job.title)
        #expect(found?.durationSec == job.durationSec)
        #expect(found?.scheduledAt == job.scheduledAt)
        #expect(found?.timezone == job.timezone)
        #expect(found?.state == job.state)
    }

    @Test
    func findReturnsNilForUnknownId() async throws {
        let repository = try makeTestRepository()

        let found = try await repository.find(jobId: "missing")

        #expect(found == nil)
    }

    @Test
    func findAllReturnsAllCreatedJobs() async throws {
        let repository = try makeTestRepository()
        let job1 = makeJob(
            jobId: "job-1",
            scheduledAt: Date(timeIntervalSince1970: 1_700_000_000)
        )
        let job2 = makeJob(
            jobId: "job-2",
            scheduledAt: Date(timeIntervalSince1970: 1_700_010_000)
        )

        try await repository.create(job1)
        try await repository.create(job2)
        let jobs = try await repository.findAll()

        #expect(jobs.count == 2)
        #expect(Set(jobs.map(\.jobId)) == Set([job1.jobId, job2.jobId]))
    }

    @Test
    func deleteReturnsTrueForExistingJobAndRemovesIt() async throws {
        let repository = try makeTestRepository()
        let job = makeJob(jobId: "job-1")
        try await repository.create(job)

        let deleted = try await repository.delete(jobId: job.jobId)
        let found = try await repository.find(jobId: job.jobId)

        #expect(deleted)
        #expect(found == nil)
    }

    @Test
    func deleteReturnsFalseForUnknownJob() async throws {
        let repository = try makeTestRepository()

        let deleted = try await repository.delete(jobId: "missing")

        #expect(!deleted)
    }

    @Test
    func updateReturnsTrueAndPersistsChanges() async throws {
        let repository = try makeTestRepository()
        let original = makeJob(jobId: "job-1")
        try await repository.create(original)

        let updatedJob = Job(
            jobId: original.jobId,
            sourceType: "url",
            sourceValue: "https://example.com/updated",
            title: "Updated Title",
            durationSec: 120,
            scheduledAt: Date(timeIntervalSince1970: 1_700_123_456),
            timezone: "UTC",
            state: .paused
        )

        let updated = try await repository.update(updatedJob)
        let found = try await repository.find(jobId: original.jobId)

        #expect(updated)
        #expect(found?.jobId == updatedJob.jobId)
        #expect(found?.sourceType == updatedJob.sourceType)
        #expect(found?.sourceValue == updatedJob.sourceValue)
        #expect(found?.title == updatedJob.title)
        #expect(found?.durationSec == updatedJob.durationSec)
        #expect(found?.scheduledAt == updatedJob.scheduledAt)
        #expect(found?.timezone == updatedJob.timezone)
        #expect(found?.state == updatedJob.state)
    }

    @Test
    func updateReturnsFalseForUnknownJob() async throws {
        let repository = try makeTestRepository()

        let updated = try await repository.update(makeJob(jobId: "missing"))

        #expect(!updated)
    }
}

private func makeTestRepository() throws -> GRDBJobRepository {
    let dbQueue = try DatabaseQueue(path: ":memory:")
    let migrator = makeDatabaseMigrator()
    try migrator.migrate(dbQueue)
    return GRDBJobRepository(dbQueue: dbQueue)
}

private func makeJob(
    jobId: String,
    scheduledAt: Date = Date(timeIntervalSince1970: 1_700_000_000),
    state: JobState = .active
) -> Job {
    Job(
        jobId: jobId,
        sourceType: "youtube",
        sourceValue: "https://example.com/\(jobId)",
        title: "Title \(jobId)",
        durationSec: 60,
        scheduledAt: scheduledAt,
        timezone: "Asia/Tokyo",
        state: state
    )
}
