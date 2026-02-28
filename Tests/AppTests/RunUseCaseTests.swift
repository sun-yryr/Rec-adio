@testable import App
import Foundation
import Testing

@Suite
struct RunUseCaseTests {
    @Test
    func createStoresRunThroughRepository() async throws {
        let repository = RunRepositoryDouble()
        let useCase = DefaultRunUseCase(runRepository: repository)
        let run = makeRun(runId: "run-1")

        try await useCase.create(run)
        let stored = try await repository.find(runId: run.runId)

        #expect(stored?.runId == run.runId)
    }

    @Test
    func createTranslatesConstraintViolationToDuplicateRunError() async {
        let repository = RunRepositoryDouble(createBehavior: .throwConstraintViolation)
        let useCase = DefaultRunUseCase(runRepository: repository)

        do {
            try await useCase.create(makeRun(runId: "run-1"))
            Issue.record("Expected create to throw")
        } catch let error as RunUseCaseError {
            switch error {
            case let .duplicateRun(reason):
                #expect(reason.contains("UNIQUE constraint failed"))
            default:
                Issue.record("Expected duplicateRun but got: \(error)")
            }
        } catch {
            Issue.record("Expected RunUseCaseError but got: \(error)")
        }
    }

    @Test
    func getReturnsStoredRun() async throws {
        let repository = RunRepositoryDouble(initialRuns: [makeRun(runId: "run-1")])
        let useCase = DefaultRunUseCase(runRepository: repository)

        let run = try await useCase.get(runId: "run-1")

        #expect(run.runId == "run-1")
    }

    @Test
    func getThrowsNotFoundForUnknownRun() async {
        let repository = RunRepositoryDouble()
        let useCase = DefaultRunUseCase(runRepository: repository)

        do {
            _ = try await useCase.get(runId: "missing")
            Issue.record("Expected get to throw")
        } catch let error as RunUseCaseError {
            switch error {
            case let .runNotFound(runId):
                #expect(runId == "missing")
            default:
                Issue.record("Expected runNotFound but got: \(error)")
            }
        } catch {
            Issue.record("Expected RunUseCaseError but got: \(error)")
        }
    }

    @Test
    func listReturnsRunsFromRepository() async throws {
        let run1 = makeRun(runId: "run-1")
        let run2 = makeRun(runId: "run-2", plannedAt: Date(timeIntervalSince1970: 1_700_000_120))
        let repository = RunRepositoryDouble(initialRuns: [run1, run2])
        let useCase = DefaultRunUseCase(runRepository: repository)

        let runs = try await useCase.list()

        #expect(runs.count == 2)
        #expect(Set(runs.map(\.runId)) == Set(["run-1", "run-2"]))
    }

    @Test
    func updatePersistsChangesAndReturnsUpdatedRun() async throws {
        let initial = makeRun(runId: "run-1")
        let repository = RunRepositoryDouble(initialRuns: [initial])
        let useCase = DefaultRunUseCase(runRepository: repository)
        let updated = Run(
            runId: initial.runId,
            jobId: initial.jobId,
            plannedAt: initial.plannedAt,
            startedAt: Date(timeIntervalSince1970: 1_700_000_005),
            finishedAt: Date(timeIntervalSince1970: 1_700_000_100),
            state: .failed,
            outputPath: initial.outputPath,
            errorMessage: "ffmpeg exited with code 1"
        )

        let result = try await useCase.update(updated)
        let stored = try await repository.find(runId: initial.runId)

        #expect(result.runId == updated.runId)
        #expect(result.state == .failed)
        #expect(stored?.errorMessage == "ffmpeg exited with code 1")
    }

    @Test
    func updateThrowsNotFoundForUnknownRun() async {
        let repository = RunRepositoryDouble(updateBehavior: .returnFalse)
        let useCase = DefaultRunUseCase(runRepository: repository)

        do {
            _ = try await useCase.update(makeRun(runId: "missing"))
            Issue.record("Expected update to throw")
        } catch let error as RunUseCaseError {
            switch error {
            case let .runNotFound(runId):
                #expect(runId == "missing")
            default:
                Issue.record("Expected runNotFound but got: \(error)")
            }
        } catch {
            Issue.record("Expected RunUseCaseError but got: \(error)")
        }
    }

    @Test
    func deleteRemovesExistingRun() async throws {
        let initial = makeRun(runId: "run-1")
        let repository = RunRepositoryDouble(initialRuns: [initial])
        let useCase = DefaultRunUseCase(runRepository: repository)

        try await useCase.delete(runId: initial.runId)
        let found = try await repository.find(runId: initial.runId)

        #expect(found == nil)
    }

    @Test
    func deleteThrowsNotFoundForUnknownRun() async {
        let repository = RunRepositoryDouble(deleteBehavior: .returnFalse)
        let useCase = DefaultRunUseCase(runRepository: repository)

        do {
            try await useCase.delete(runId: "missing")
            Issue.record("Expected delete to throw")
        } catch let error as RunUseCaseError {
            switch error {
            case let .runNotFound(runId):
                #expect(runId == "missing")
            default:
                Issue.record("Expected runNotFound but got: \(error)")
            }
        } catch {
            Issue.record("Expected RunUseCaseError but got: \(error)")
        }
    }
}

private func makeRun(
    runId: String,
    plannedAt: Date = Date(timeIntervalSince1970: 1_700_000_000),
    state: RunState = .queued
) -> Run {
    Run(
        runId: runId,
        jobId: "job-1",
        plannedAt: plannedAt,
        startedAt: nil,
        finishedAt: nil,
        state: state,
        outputPath: "/tmp/recordings/\(runId).m4a",
        errorMessage: ""
    )
}

private actor RunRepositoryDouble: RunRepository {
    enum CreateBehavior {
        case succeed
        case throwConstraintViolation
        case throwUnexpectedError
    }

    enum UpdateBehavior {
        case useStorage
        case returnFalse
        case throwUnexpectedError
    }

    enum DeleteBehavior {
        case useStorage
        case returnFalse
        case throwUnexpectedError
    }

    private let createBehavior: CreateBehavior
    private let updateBehavior: UpdateBehavior
    private let deleteBehavior: DeleteBehavior
    private var runs: [String: Run]

    init(
        createBehavior: CreateBehavior = .succeed,
        updateBehavior: UpdateBehavior = .useStorage,
        deleteBehavior: DeleteBehavior = .useStorage,
        initialRuns: [Run] = []
    ) {
        self.createBehavior = createBehavior
        self.updateBehavior = updateBehavior
        self.deleteBehavior = deleteBehavior
        self.runs = Dictionary(uniqueKeysWithValues: initialRuns.map { ($0.runId, $0) })
    }

    func create(_ run: Run) async throws {
        switch createBehavior {
        case .succeed:
            runs[run.runId] = run
        case .throwConstraintViolation:
            throw RunRepositoryError.duplicateRun(
                reason: "UNIQUE constraint failed: runs.run_id"
            )
        case .throwUnexpectedError:
            throw RunUseCaseTestError.unexpected
        }
    }

    func find(runId: String) async throws -> Run? {
        runs[runId]
    }

    func findAll() async throws -> [Run] {
        runs.values.sorted { $0.plannedAt < $1.plannedAt }
    }

    func update(_ run: Run) async throws -> Bool {
        switch updateBehavior {
        case .useStorage:
            guard runs[run.runId] != nil else {
                return false
            }
            runs[run.runId] = run
            return true
        case .returnFalse:
            return false
        case .throwUnexpectedError:
            throw RunUseCaseTestError.unexpected
        }
    }

    func delete(runId: String) async throws -> Bool {
        switch deleteBehavior {
        case .useStorage:
            guard runs[runId] != nil else {
                return false
            }
            runs[runId] = nil
            return true
        case .returnFalse:
            return false
        case .throwUnexpectedError:
            throw RunUseCaseTestError.unexpected
        }
    }
}

private enum RunUseCaseTestError: Error {
    case unexpected
}
