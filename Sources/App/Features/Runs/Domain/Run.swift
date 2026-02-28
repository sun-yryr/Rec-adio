import Foundation

struct Run: Codable {
    let runId: String
    let jobId: String
    let plannedAt: Date
    let startedAt: Date?
    let finishedAt: Date?
    let state: RunState
    let outputPath: String
    let errorMessage: String?
    let createdAt: Date
    let updatedAt: Date

    init(
        runId: String,
        jobId: String,
        plannedAt: Date,
        startedAt: Date? = nil,
        finishedAt: Date? = nil,
        state: RunState,
        outputPath: String,
        errorMessage: String? = nil,
        createdAt: Date = Date(),
        updatedAt: Date = Date()
    ) {
        self.runId = runId
        self.jobId = jobId
        self.plannedAt = plannedAt
        self.startedAt = startedAt
        self.finishedAt = finishedAt
        self.state = state
        self.outputPath = outputPath
        self.errorMessage = errorMessage
        self.createdAt = createdAt
        self.updatedAt = updatedAt
    }
}

enum RunState: String, Codable {
    case queued, running, succeeded, failed, canceled
}
