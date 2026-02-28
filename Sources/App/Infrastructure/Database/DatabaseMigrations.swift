import GRDB

func makeDatabaseMigrator() -> DatabaseMigrator {
    var migrator = DatabaseMigrator()
    migrator.registerMigration("create_jobs_table") { database in
        try createJobsTable(database)
        try createJobsIndexes(database)
    }

    migrator.registerMigration(
        "create_runs_table",
        migrate: { database in
            try createRunsTable(database)
            try createRunsIndexes(database)
        }
    )
    return migrator
}

private let validJobStates = "'active','paused','deleted'"
private let validRunStates = "'queued','running','succeeded','failed','canceled'"

private func createJobsTable(_ database: Database) throws {
    try database.execute(
        sql: """
        CREATE TABLE IF NOT EXISTS jobs (
          job_id TEXT PRIMARY KEY,
          source_type TEXT NOT NULL,
          source_value TEXT NOT NULL,
          title TEXT NOT NULL,
          duration_sec INTEGER NOT NULL CHECK(duration_sec > 0),
          scheduled_at TEXT NOT NULL,
          timezone TEXT NOT NULL,
          state TEXT NOT NULL CHECK(state IN (\(validJobStates))),
          created_at TEXT NOT NULL,
          updated_at TEXT NOT NULL
        );
        """
    )
}

private func createRunsTable(_ database: Database) throws {
    try database.execute(
        sql: """
        CREATE TABLE IF NOT EXISTS runs (
          run_id TEXT PRIMARY KEY,
          job_id TEXT NOT NULL,
          planned_at TEXT NOT NULL,
          started_at TEXT,
          finished_at TEXT,
          state TEXT NOT NULL CHECK(state IN (\(validRunStates))),
          output_path TEXT NOT NULL,
          error_message TEXT,
          created_at TEXT NOT NULL,
          updated_at TEXT NOT NULL,
          FOREIGN KEY(job_id) REFERENCES jobs(job_id)
        );
        """
    )
}

private func createJobsIndexes(_ database: Database) throws {
    try database.execute(
        sql: """
        CREATE INDEX IF NOT EXISTS idx_jobs_state_scheduled_at
          ON jobs(state, scheduled_at);
        """
    )
}

private func createRunsIndexes(_ database: Database) throws {
    try database.execute(
        sql: """
        CREATE INDEX IF NOT EXISTS idx_runs_job_id_created_at
          ON runs(job_id, created_at DESC);
        """
    )
    try database.execute(
        sql: """
        CREATE INDEX IF NOT EXISTS idx_runs_state_updated_at
          ON runs(state, updated_at DESC);
        """
    )
}
