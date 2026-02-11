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
    try await dbQueue.write { db in
      try job.insert(db)
    }
  }

  func find(jobId: String) async throws -> Job? {
    try await dbQueue.read { db in
      try Job.fetchOne(db, key: jobId)
    }
  }

  func findAll() async throws -> [Job] {
    try await dbQueue.read { db in
      try Job
        .order(Job.Columns.scheduledAt)
        .fetchAll(db)
    }
  }

  func delete(jobId: String) async throws -> Bool {
    try await dbQueue.write { db in
      try Job.deleteOne(db, key: jobId)
    }
  }
}
