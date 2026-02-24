import GRDB

protocol JobServicing: Sendable {
    func create(_ job: Job) async throws
    func list() async throws -> [Job]
}

enum JobServiceError: Error, Sendable {
    case duplicateJob(reason: String)
    case createFailed(reason: String)
    case listFailed(reason: String)

    var reason: String {
        switch self {
        case let .duplicateJob(reason):
            return reason
        case let .createFailed(reason):
            return reason
        case let .listFailed(reason):
            return reason
        }
    }
}

struct JobService: JobServicing {
    private let jobRepo: any JobRepository

    init(jobRepo: any JobRepository) {
        self.jobRepo = jobRepo
    }

    func create(_ job: Job) async throws {
        do {
            try await jobRepo.create(job)
        } catch let error as DatabaseError where error.resultCode == .SQLITE_CONSTRAINT {
            throw JobServiceError.duplicateJob(
                reason: error.message ?? "constraint violation"
            )
        } catch {
            throw JobServiceError.createFailed(reason: String(describing: error))
        }
    }

    func list() async throws -> [Job] {
        do {
            return try await jobRepo.findAll()
        } catch {
            throw JobServiceError.listFailed(reason: String(describing: error))
        }
    }
}
