import GRDB

extension Run: Identifiable, TableRecord, FetchableRecord, PersistableRecord {
    static let databaseTableName = "runs"
    static let databaseColumnEncodingStrategy: DatabaseColumnEncodingStrategy = .convertToSnakeCase
    static let databaseColumnDecodingStrategy: DatabaseColumnDecodingStrategy = .convertFromSnakeCase
    static let databaseDateEncodingStrategy = DatabaseDateEncodingStrategy.iso8601
    static let databaseDateDecodingStrategy = DatabaseDateDecodingStrategy.iso8601

    var id: String {
        runId
    }

    enum Columns {
        static let runId = Column("run_id")
        static let jobId = Column("job_id")
        static let plannedAt = Column("planned_at")
        static let startedAt = Column("started_at")
        static let finishedAt = Column("finished_at")
        static let state = Column("state")
        static let outputPath = Column("output_path")
        static let errorMessage = Column("error_message")
        static let createdAt = Column("created_at")
        static let updatedAt = Column("updated_at")
    }
}
