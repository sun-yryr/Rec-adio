import Fluent

struct CreateSchedule: AsyncMigration {
    func prepare(on database: Database) async throws {
        try await database.schema("schedules")
            .field(.id, .int, .identifier(auto: true))
            .field("title", .string, .required)
            .field("start_datetime", .datetime, .required)
            .field("duration", .int, .required)
            .field("is_processing", .bool, .required)
            .field("platform", .string, .required)
            .create()
    }

    func revert(on database: Database) async throws {
        try await database.schema("schedules").delete()
    }
}