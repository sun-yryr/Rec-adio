import Fluent

struct CreateTwitterUsers: AsyncMigration {
    func prepare(on database: Database) async throws {
        try await database.schema("twitter_users")
            .field(.id, .int, .identifier(auto: true))
            .field("name", .string, .required)
            .field("user_id", .string, .required)
            .field("program_info_id", .int, .required, .references("program_infos", "id"))
            .create()
    }

    func revert(on database: Database) async throws {
        try await database.schema("twitter_users").delete()
    }
}
