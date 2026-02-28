import GRDB

protocol JobRepository: Sendable {
    func create(_ job: Job) async throws
    func find(jobId: String) async throws -> Job?
    func findAll() async throws -> [Job]
    @discardableResult
    func update(_ job: Job) async throws -> Bool
    @discardableResult
    func delete(jobId: String) async throws -> Bool
}

actor GRDBJobRepository: JobRepository {
    private let dbQueue: DatabaseQueue

    init(dbQueue: DatabaseQueue) {
        self.dbQueue = dbQueue
    }

    func create(_ job: Job) async throws {
        try await dbQueue.write { database in
            try job.insert(database)
        }
    }

    func find(jobId: String) async throws -> Job? {
        try await dbQueue.read { database in
            try Job.fetchOne(database, key: jobId)
        }
    }

    func findAll() async throws -> [Job] {
        try await dbQueue.read { database in
            try Job
                .order(Job.Columns.scheduledAt)
                .fetchAll(database)
        }
    }

    func update(_ job: Job) async throws -> Bool {
        try await dbQueue.write { database in
            try Job
                .filter(Job.Columns.jobId == job.jobId)
                .updateAll(
                    database,
                    [
                        Job.Columns.sourceType.set(to: job.sourceType),
                        Job.Columns.sourceValue.set(to: job.sourceValue),
                        Job.Columns.title.set(to: job.title),
                        Job.Columns.durationSec.set(to: job.durationSec),
                        Job.Columns.scheduledAt.set(to: job.scheduledAt),
                        Job.Columns.timezone.set(to: job.timezone),
                        Job.Columns.state.set(to: job.state.rawValue),
                    ]
                ) == 1
        }
    }

    func delete(jobId: String) async throws -> Bool {
        try await dbQueue.write { database in
            try Job.deleteOne(database, key: jobId)
        }
    }
}
