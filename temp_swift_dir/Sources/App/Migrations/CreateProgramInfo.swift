import FluentKit

struct CreateProgramInfo: AsyncMigration {
    func prepare(on database: Database) async throws {
        try await database.schema("program_infos")
            .field(.id, .int, .identifier(auto: true))
            .field("title", .string, .required)
            .field("platform", .string, .required)
            .field("pfm", .string, .required)
            .field("program_url", .string, .required)
            .create()
    }

    func revert(on database: Database) async throws {
        try await database.schema("program_infos").delete()
    }
}