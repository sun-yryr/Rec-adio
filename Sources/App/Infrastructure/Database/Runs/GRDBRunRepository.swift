import GRDB

actor GRDBRunRepository: RunRepository {
    private let dbQueue: DatabaseQueue

    init(dbQueue: DatabaseQueue) {
        self.dbQueue = dbQueue
    }

    func create(_ run: Run) async throws {
        do {
            try await dbQueue.write { database in
                try run.insert(database)
            }
        } catch let error as DatabaseError where error.resultCode == .SQLITE_CONSTRAINT {
            if let message = error.message, message.contains("UNIQUE constraint failed") {
                throw RunRepositoryError.duplicateRun(reason: message)
            }
            throw error
        }
    }

    func find(runId: String) async throws -> Run? {
        try await dbQueue.read { database in
            try Run.fetchOne(database, key: runId)
        }
    }

    func findAll() async throws -> [Run] {
        try await dbQueue.read { database in
            try Run
                .order(Run.Columns.plannedAt)
                .fetchAll(database)
        }
    }

    func update(_ run: Run) async throws -> Bool {
        try await dbQueue.write { database in
            try Run
                .filter(Run.Columns.runId == run.runId)
                .updateAll(
                    database,
                    [
                        Run.Columns.jobId.set(to: run.jobId),
                        Run.Columns.plannedAt.set(to: run.plannedAt),
                        Run.Columns.startedAt.set(to: run.startedAt),
                        Run.Columns.finishedAt.set(to: run.finishedAt),
                        Run.Columns.state.set(to: run.state.rawValue),
                        Run.Columns.outputPath.set(to: run.outputPath),
                        Run.Columns.errorMessage.set(to: run.errorMessage),
                    ]
                ) == 1
        }
    }

    func delete(runId: String) async throws -> Bool {
        try await dbQueue.write { database in
            try Run.deleteOne(database, key: runId)
        }
    }
}
