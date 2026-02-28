protocol RunUseCase: Sendable {
    func create(_ run: Run) async throws
    func get(runId: String) async throws -> Run
    func list() async throws -> [Run]
    func update(_ run: Run) async throws -> Run
    func delete(runId: String) async throws
}

enum RunUseCaseError: Error, Sendable {
    case duplicateRun(reason: String)
    case runNotFound(runId: String)
    case createFailed(reason: String)
    case getFailed(reason: String)
    case listFailed(reason: String)
    case updateFailed(reason: String)
    case deleteFailed(reason: String)

    var reason: String {
        switch self {
        case let .duplicateRun(reason):
            return reason
        case let .createFailed(reason):
            return reason
        case let .getFailed(reason):
            return reason
        case let .listFailed(reason):
            return reason
        case let .updateFailed(reason):
            return reason
        case let .deleteFailed(reason):
            return reason
        case let .runNotFound(runId):
            return "run not found: \(runId)"
        }
    }
}

struct DefaultRunUseCase: RunUseCase {
    private let runRepository: any RunRepository

    init(runRepository: any RunRepository) {
        self.runRepository = runRepository
    }

    func create(_ run: Run) async throws {
        do {
            try await runRepository.create(run)
        } catch let error as RunRepositoryError {
            switch error {
            case let .duplicateRun(reason):
                throw RunUseCaseError.duplicateRun(reason: reason)
            }
        } catch {
            throw RunUseCaseError.createFailed(reason: String(describing: error))
        }
    }

    func get(runId: String) async throws -> Run {
        do {
            guard let run = try await runRepository.find(runId: runId) else {
                throw RunUseCaseError.runNotFound(runId: runId)
            }
            return run
        } catch let error as RunUseCaseError {
            throw error
        } catch {
            throw RunUseCaseError.getFailed(reason: String(describing: error))
        }
    }

    func list() async throws -> [Run] {
        do {
            return try await runRepository.findAll()
        } catch {
            throw RunUseCaseError.listFailed(reason: String(describing: error))
        }
    }

    func update(_ run: Run) async throws -> Run {
        do {
            let updated = try await runRepository.update(run)
            guard updated else {
                throw RunUseCaseError.runNotFound(runId: run.runId)
            }
            return run
        } catch let error as RunUseCaseError {
            throw error
        } catch {
            throw RunUseCaseError.updateFailed(reason: String(describing: error))
        }
    }

    func delete(runId: String) async throws {
        do {
            let deleted = try await runRepository.delete(runId: runId)
            guard deleted else {
                throw RunUseCaseError.runNotFound(runId: runId)
            }
        } catch let error as RunUseCaseError {
            throw error
        } catch {
            throw RunUseCaseError.deleteFailed(reason: String(describing: error))
        }
    }
}
