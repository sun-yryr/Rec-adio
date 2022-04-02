import FluentKit

struct CreateRecord: AsyncMigration {
    func prepare(on database: Database) async throws {
        try await database.schema("records")
            .field(.id, .int, .identifier(auto: true))
            .field("title", .string, .required)
            .field("record_datetime", .datetime, .required)
            .field("duration", .int, .required)
            .field("path", .string, .required)
            .field("platform", .string, .required)
            .field("program_info_id", .int, .required, .references("program_infos", "id"))
            .create()
    }

    func revert(on database: Database) async throws {
        try await database.schema("records").delete()
    }
}