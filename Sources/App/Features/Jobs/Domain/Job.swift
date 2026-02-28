import Foundation

struct Job: Codable {
    let jobId: String
    let sourceType: String
    let sourceValue: String
    let title: String
    let durationSec: Int64
    let scheduledAt: Date
    let timezone: String
    let state: JobState
    let createdAt: Date
    let updatedAt: Date

    init(
        jobId: String,
        sourceType: String,
        sourceValue: String,
        title: String,
        durationSec: Int64,
        scheduledAt: Date,
        timezone: String,
        state: JobState,
        createdAt: Date = Date(),
        updatedAt: Date = Date()
    ) {
        self.jobId = jobId
        self.sourceType = sourceType
        self.sourceValue = sourceValue
        self.title = title
        self.durationSec = durationSec
        self.scheduledAt = scheduledAt
        self.timezone = timezone
        self.state = state
        self.createdAt = createdAt
        self.updatedAt = updatedAt
    }
}

enum JobState: String, Codable {
    case active
    case paused
    case deleted
}

extension Job {
    struct Updatable: Sendable {
        let title: String?
        let durationSec: Int64?
        let scheduledAt: Date?
        let timezone: String?

        var isEmpty: Bool {
            title == nil && durationSec == nil && scheduledAt == nil && timezone == nil
        }
    }
}
