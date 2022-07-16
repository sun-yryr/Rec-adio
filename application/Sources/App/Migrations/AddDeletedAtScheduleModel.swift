import Fluent

struct AddDeletedAtScheduleModel: AsyncMigration {
    func prepare(on database: Database) async throws {
        try await database.schema("schedules")
            .field("deleted_at", .datetime)
            .update()
    }

    func revert(on database: Database) async throws {
        try await database.schema("schedules")
            .deleteField("deleted_at")
            .update()
    }
}
