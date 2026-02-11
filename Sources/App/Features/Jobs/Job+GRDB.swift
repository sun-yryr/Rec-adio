import GRDB

extension Job: Identifiable, TableRecord, FetchableRecord, PersistableRecord {
  static let databaseTableName = "jobs"
  static let databaseColumnEncodingStrategy = DatabaseColumnEncodingStrategy.convertToSnakeCase
  static let databaseColumnDecodingStrategy = DatabaseColumnDecodingStrategy.convertFromSnakeCase

  var id: String { jobId }

  enum Columns {
    static let jobId = Column("job_id")
    static let sourceType = Column("source_type")
    static let sourceValue = Column("source_value")
    static let title = Column("title")
    static let durationSec = Column("duration_sec")
    static let scheduledAt = Column("scheduled_at")
    static let timezone = Column("timezone")
    static let state = Column("state")
  }
}
