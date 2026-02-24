import GRDB

protocol JobRepository: Sendable {
    func create(_ job: Job) async throws
    func find(jobId: String) async throws -> Job?
    func findAll() async throws -> [Job]
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

    func delete(jobId: String) async throws -> Bool {
        try await dbQueue.write { database in
            try Job.deleteOne(database, key: jobId)
        }
    }
}
