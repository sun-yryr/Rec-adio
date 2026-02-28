protocol RunRepository: Sendable {
    func create(_ run: Run) async throws
    func find(runId: String) async throws -> Run?
    func findAll() async throws -> [Run]
    @discardableResult
    func update(_ run: Run) async throws -> Bool
    @discardableResult
    func delete(runId: String) async throws -> Bool
}

enum RunRepositoryError: Error, Sendable {
    case duplicateRun(reason: String)

    var reason: String {
        switch self {
        case let .duplicateRun(reason):
            return reason
        }
    }
}
