import GRDB

protocol JobUseCase: Sendable {
    func create(_ job: Job) async throws
    func get(jobId: String) async throws -> Job
    func list() async throws -> [Job]
    func update(_ job: Job) async throws -> Job
    func delete(jobId: String) async throws
}

enum JobUseCaseError: Error, Sendable {
    case duplicateJob(reason: String)
    case jobNotFound(jobId: String)
    case createFailed(reason: String)
    case getFailed(reason: String)
    case listFailed(reason: String)
    case updateFailed(reason: String)
    case deleteFailed(reason: String)

    var reason: String {
        switch self {
        case let .duplicateJob(reason):
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
        case let .jobNotFound(jobId):
            return "job not found: \(jobId)"
        }
    }
}

struct DefaultJobUseCase: JobUseCase {
    private let jobRepository: any JobRepository

    init(jobRepository: any JobRepository) {
        self.jobRepository = jobRepository
    }

    func create(_ job: Job) async throws {
        do {
            try await jobRepository.create(job)
        } catch let error as DatabaseError where error.resultCode == .SQLITE_CONSTRAINT {
            throw JobUseCaseError.duplicateJob(
                reason: error.message ?? "constraint violation"
            )
        } catch {
            throw JobUseCaseError.createFailed(reason: String(describing: error))
        }
    }

    func get(jobId: String) async throws -> Job {
        do {
            guard let job = try await jobRepository.find(jobId: jobId) else {
                throw JobUseCaseError.jobNotFound(jobId: jobId)
            }
            return job
        } catch let error as JobUseCaseError {
            throw error
        } catch {
            throw JobUseCaseError.getFailed(reason: String(describing: error))
        }
    }

    func list() async throws -> [Job] {
        do {
            return try await jobRepository.findAll()
        } catch {
            throw JobUseCaseError.listFailed(reason: String(describing: error))
        }
    }

    func update(_ job: Job) async throws -> Job {
        do {
            let updated = try await jobRepository.update(job)
            guard updated else {
                throw JobUseCaseError.jobNotFound(jobId: job.jobId)
            }
            return job
        } catch let error as JobUseCaseError {
            throw error
        } catch {
            throw JobUseCaseError.updateFailed(reason: String(describing: error))
        }
    }

    func delete(jobId: String) async throws {
        do {
            let deleted = try await jobRepository.delete(jobId: jobId)
            guard deleted else {
                throw JobUseCaseError.jobNotFound(jobId: jobId)
            }
        } catch let error as JobUseCaseError {
            throw error
        } catch {
            throw JobUseCaseError.deleteFailed(reason: String(describing: error))
        }
    }
}
