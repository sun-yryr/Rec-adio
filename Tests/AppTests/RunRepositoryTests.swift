@testable import App
import Foundation
import GRDB
import Testing

@Suite
struct RunRepositoryTests {
    @Test
    func createThenFindReturnsStoredRun() async throws {
        let repository = try makeTestRepository()
        let run = makeRun(runId: "run-1")

        try await repository.create(run)
        let found = try await repository.find(runId: run.runId)

        #expect(found != nil)
        #expect(found?.runId == run.runId)
        #expect(found?.jobId == run.jobId)
        #expect(found?.plannedAt == run.plannedAt)
        #expect(found?.startedAt == run.startedAt)
        #expect(found?.finishedAt == run.finishedAt)
        #expect(found?.state == run.state)
        #expect(found?.outputPath == run.outputPath)
        #expect(found?.errorMessage == run.errorMessage)
    }

    @Test
    func findReturnsNilForUnknownId() async throws {
        let repository = try makeTestRepository()

        let found = try await repository.find(runId: "missing")

        #expect(found == nil)
    }

    @Test
    func findAllReturnsAllCreatedRuns() async throws {
        let repository = try makeTestRepository()
        let run1 = makeRun(
            runId: "run-1",
            plannedAt: Date(timeIntervalSince1970: 1_700_000_000)
        )
        let run2 = makeRun(
            runId: "run-2",
            plannedAt: Date(timeIntervalSince1970: 1_700_000_600)
        )

        try await repository.create(run1)
        try await repository.create(run2)
        let runs = try await repository.findAll()

        #expect(runs.count == 2)
        #expect(runs.map(\.runId) == ["run-1", "run-2"])
    }

    @Test
    func updateReturnsTrueAndPersistsChanges() async throws {
        let repository = try makeTestRepository()
        let original = makeRun(runId: "run-1")
        try await repository.create(original)

        let updatedRun = Run(
            runId: original.runId,
            jobId: original.jobId,
            plannedAt: original.plannedAt,
            startedAt: Date(timeIntervalSince1970: 1_700_000_010),
            finishedAt: Date(timeIntervalSince1970: 1_700_000_120),
            state: .succeeded,
            outputPath: "/tmp/recordings/updated.m4a",
            errorMessage: ""
        )

        let updated = try await repository.update(updatedRun)
        let found = try await repository.find(runId: original.runId)

        #expect(updated)
        #expect(found?.runId == updatedRun.runId)
        #expect(found?.state == updatedRun.state)
        #expect(found?.startedAt == updatedRun.startedAt)
        #expect(found?.finishedAt == updatedRun.finishedAt)
        #expect(found?.outputPath == updatedRun.outputPath)
    }

    @Test
    func updateReturnsFalseForUnknownRun() async throws {
        let repository = try makeTestRepository()

        let updated = try await repository.update(makeRun(runId: "missing"))

        #expect(!updated)
    }

    @Test
    func deleteReturnsTrueForExistingRunAndRemovesIt() async throws {
        let repository = try makeTestRepository()
        let run = makeRun(runId: "run-1")
        try await repository.create(run)

        let deleted = try await repository.delete(runId: run.runId)
        let found = try await repository.find(runId: run.runId)

        #expect(deleted)
        #expect(found == nil)
    }

    @Test
    func deleteReturnsFalseForUnknownRun() async throws {
        let repository = try makeTestRepository()

        let deleted = try await repository.delete(runId: "missing")

        #expect(!deleted)
    }
}

private func makeTestRepository() throws -> GRDBRunRepository {
    let dbQueue = try DatabaseQueue(path: ":memory:")
    let migrator = makeDatabaseMigrator()
    try migrator.migrate(dbQueue)
    try dbQueue.write { database in
        try database.execute(
            sql: """
            INSERT INTO jobs (
              job_id, source_type, source_value, title, duration_sec, scheduled_at, timezone, state, created_at, updated_at
            ) VALUES (
              'job-1', 'url', 'https://example.com/job-1', 'Seed Job', 60, '2024-01-01T00:00:00Z', 'UTC', 'active', '2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z'
            );
            """
        )
    }
    return GRDBRunRepository(dbQueue: dbQueue)
}

private func makeRun(
    runId: String,
    plannedAt: Date = Date(timeIntervalSince1970: 1_700_000_000),
    state: RunState = .queued
) -> Run {
    Run(
        runId: runId,
        jobId: "job-1",
        plannedAt: plannedAt,
        startedAt: nil,
        finishedAt: nil,
        state: state,
        outputPath: "/tmp/recordings/\(runId).m4a",
        errorMessage: ""
    )
}
