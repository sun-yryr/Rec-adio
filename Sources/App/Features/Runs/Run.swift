import Foundation

struct Run: Codable {
    let runId: String
    let jobId: String
    let plannedAt: Date
    let startedAt: Date?
    let finishedAt: Date?
    let state: RunState
    let outputPath: String
    let errorMessage: String
}

enum RunState: String, Codable {
    case queued, running, succeeded, failed, canceled
}
