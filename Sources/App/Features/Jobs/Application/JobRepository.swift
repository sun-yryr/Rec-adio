protocol JobRepository: Sendable {
    func create(_ job: Job) async throws
    func find(jobId: String) async throws -> Job?
    func findAll() async throws -> [Job]
    @discardableResult
    func update(jobId: String, updatable: Job.Updatable) async throws -> Bool
    @discardableResult
    func delete(jobId: String) async throws -> Bool
}

enum JobRepositoryError: Error, Sendable {
    case duplicateJob(reason: String)

    var reason: String {
        switch self {
        case let .duplicateJob(reason):
            return reason
        }
    }
}
