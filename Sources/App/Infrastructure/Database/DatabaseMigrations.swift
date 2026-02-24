import GRDB

func makeDatabaseMigrator() -> DatabaseMigrator {
    var migrator = DatabaseMigrator()
    migrator.registerMigration("create_jobs_table") { database in
        try database.create(table: "jobs") { table in
            table.primaryKey("job_id", .text)
            table.column("source_type", .text).notNull()
            table.column("source_value", .text).notNull()
            table.column("title", .text).notNull()
            table.column("duration_sec", .integer).notNull().check { column in
                column > 0
            }
            table.column("scheduled_at", .datetime).notNull()
            table.column("timezone", .text).notNull()
            table.column("state", .text).notNull()
        }
    }

    migrator.registerMigration(
        "create_runs_table",
        migrate: { database in
            try database.create(table: "runs") { table in
                table.primaryKey("run_id", .text)
                table.column("job_id", .text).notNull()
                table.column("planned_at", .datetime).notNull()
                table.column("started_at", .datetime)
                table.column("finished_at", .datetime)
                table.column("output_path", .text).notNull()
                table.column("error_message", .text)
                table.column("state", .text).notNull()
            }
        }
    )
    return migrator
}
