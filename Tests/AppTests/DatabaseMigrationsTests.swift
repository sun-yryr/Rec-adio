@testable import App
import GRDB
import Testing

@Suite
struct DatabaseMigrationsTests {
    @Test
    func createsJobsAndRunsSchemaAndIndexes() throws {
        let dbQueue = try makeMigratedDatabaseQueue()

        try dbQueue.read { database in
            #expect(try database.tableExists("jobs"))
            #expect(try database.tableExists("runs"))

            #expect(try hasColumn("created_at", in: "jobs", database: database))
            #expect(try hasColumn("updated_at", in: "jobs", database: database))
            #expect(try hasColumn("created_at", in: "runs", database: database))
            #expect(try hasColumn("updated_at", in: "runs", database: database))

            #expect(try indexExists("idx_jobs_state_scheduled_at", database: database))
            #expect(try indexExists("idx_runs_job_id_created_at", database: database))
            #expect(try indexExists("idx_runs_state_updated_at", database: database))

            let foreignKeys = try Row.fetchAll(database, sql: "PRAGMA foreign_key_list(runs);")
            #expect(foreignKeys.count == 1)
            #expect((foreignKeys.first?["table"] as String?) == "jobs")
            #expect((foreignKeys.first?["from"] as String?) == "job_id")
            #expect((foreignKeys.first?["to"] as String?) == "job_id")
        }
    }

    @Test
    func enforcesStateChecksAndForeignKeyConstraint() throws {
        let dbQueue = try makeMigratedDatabaseQueue()

        try dbQueue.write { database in
            do {
                try database.execute(
                    sql: """
                    INSERT INTO jobs (
                      job_id,
                      source_type,
                      source_value,
                      title,
                      duration_sec,
                      scheduled_at,
                      timezone,
                      state,
                      created_at,
                      updated_at
                    ) VALUES (
                      'job-invalid',
                      'url',
                      'https://example.com/live.m3u8',
                      'Invalid Job',
                      120,
                      '2026-02-10T00:00:00Z',
                      'Asia/Tokyo',
                      'invalid-state',
                      '2026-02-10T00:00:00Z',
                      '2026-02-10T00:00:00Z'
                    );
                    """
                )
                Issue.record("Expected jobs.state CHECK constraint violation")
            } catch let error as DatabaseError {
                #expect(error.resultCode == .SQLITE_CONSTRAINT)
            }

            try database.execute(
                sql: """
                INSERT INTO jobs (
                  job_id,
                  source_type,
                  source_value,
                  title,
                  duration_sec,
                  scheduled_at,
                  timezone,
                  state,
                  created_at,
                  updated_at
                ) VALUES (
                  'job-1',
                  'url',
                  'https://example.com/live.m3u8',
                  'Valid Job',
                  120,
                  '2026-02-10T00:00:00Z',
                  'Asia/Tokyo',
                  'active',
                  '2026-02-10T00:00:00Z',
                  '2026-02-10T00:00:00Z'
                );
                """
            )

            do {
                try database.execute(
                    sql: """
                    INSERT INTO runs (
                      run_id,
                      job_id,
                      planned_at,
                      started_at,
                      finished_at,
                      state,
                      output_path,
                      error_message,
                      created_at,
                      updated_at
                    ) VALUES (
                      'run-invalid',
                      'missing-job',
                      '2026-02-10T00:10:00Z',
                      NULL,
                      NULL,
                      'queued',
                      '/tmp/recoto-output.mp4',
                      NULL,
                      '2026-02-10T00:10:00Z',
                      '2026-02-10T00:10:00Z'
                    );
                    """
                )
                Issue.record("Expected runs.job_id FOREIGN KEY constraint violation")
            } catch let error as DatabaseError {
                #expect(error.resultCode == .SQLITE_CONSTRAINT)
            }
        }
    }
}

private func makeMigratedDatabaseQueue() throws -> DatabaseQueue {
    var configuration = Configuration()
    configuration.prepareDatabase { database in
        try database.execute(sql: "PRAGMA foreign_keys=ON;")
    }

    let dbQueue = try DatabaseQueue(path: ":memory:", configuration: configuration)
    try makeDatabaseMigrator().migrate(dbQueue)
    return dbQueue
}

private func hasColumn(_ name: String, in table: String, database: Database) throws -> Bool {
    let columns = try database.columns(in: table).map { $0.name.lowercased() }
    return columns.contains(name.lowercased())
}

private func indexExists(_ name: String, database: Database) throws -> Bool {
    let count = try Int.fetchOne(
        database,
        sql: "SELECT COUNT(*) FROM sqlite_master WHERE type = 'index' AND name = ?;",
        arguments: [name]
    ) ?? 0
    return count == 1
}
