import Foundation
import GRDB

actor GRDBJobRepository: JobRepository {
    private let dbQueue: DatabaseQueue

    init(dbQueue: DatabaseQueue) {
        self.dbQueue = dbQueue
    }

    func create(_ job: Job) async throws {
        do {
            try await dbQueue.write { database in
                try job.insert(database)
            }
        } catch let error as DatabaseError where error.resultCode == .SQLITE_CONSTRAINT {
            if let message = error.message, message.contains("UNIQUE constraint failed") {
                throw JobRepositoryError.duplicateJob(reason: message)
            }
            throw error
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

    func update(jobId: String, updatable: Job.Updatable) async throws -> Bool {
        try await dbQueue.write { database in
            var columns: [ColumnAssignment] = []
            if let title = updatable.title {
                columns.append(Job.Columns.title.set(to: title))
            }
            if let durationSec = updatable.durationSec {
                columns.append(Job.Columns.durationSec.set(to: durationSec))
            }
            if let scheduledAt = updatable.scheduledAt {
                columns.append(Job.Columns.scheduledAt.set(to: scheduledAt))
            }
            if let timezone = updatable.timezone {
                columns.append(Job.Columns.timezone.set(to: timezone))
            }
            columns.append(Job.Columns.updatedAt.set(to: Date()))

            return try Job
                .filter(Job.Columns.jobId == jobId)
                .updateAll(database, columns) == 1
        }
    }

    func delete(jobId: String) async throws -> Bool {
        try await dbQueue.write { database in
            try Job.deleteOne(database, key: jobId)
        }
    }
}
