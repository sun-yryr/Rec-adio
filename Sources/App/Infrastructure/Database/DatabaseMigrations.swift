import GRDB

func makeDatabaseMigrator() -> DatabaseMigrator {
    var migrator = DatabaseMigrator()
    migrator.registerMigration("create_jobs_table") { db in
        try db.create(table: "jobs") { t in
            t.primaryKey("job_id", .text)
            t.column("source_type", .text).notNull()
            t.column("source_value", .text).notNull()
            t.column("title", .text).notNull()
            t.column("duration_sec", .integer).notNull().check { column in
                column > 0
            }
            t.column("scheduled_at", .datetime).notNull()
            t.column("timezone", .text).notNull()
            t.column("state", .text).notNull()
        }
    }

    migrator.registerMigration(
        "create_runs_table",
        migrate: { db in
            try db.create(table: "runs") { t in
                t.primaryKey("run_id", .text)
                t.column("job_id", .text).notNull()
                t.column("planned_at", .datetime).notNull()
                t.column("started_at", .datetime)
                t.column("finished_at", .datetime)
                t.column("output_path", .text).notNull()
                t.column("error_message", .text)
                t.column("state", .text).notNull()
            }
        }
    )
    return migrator
}
