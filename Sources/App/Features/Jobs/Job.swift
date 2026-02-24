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
}

enum JobState: String, Codable {
    case active
    case paused
    case deleted
}
