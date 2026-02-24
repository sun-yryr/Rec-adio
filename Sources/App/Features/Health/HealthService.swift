import GRDB
import GRPCCore
import Logging

protocol DatabaseHealthChecking: Sendable {
    func check() async throws
}

actor GRDBDatabaseHealthChecker: DatabaseHealthChecking {
    private let dbQueue: DatabaseQueue

    init(dbQueue: DatabaseQueue) {
        self.dbQueue = dbQueue
    }

    func check() async throws {
        try await dbQueue.read { database in
            _ = try Row.fetchOne(database, sql: "SELECT 1")
        }
    }
}

struct HealthService: Recoto_Health_V1_HealthService.SimpleServiceProtocol {
    private let databaseChecker: any DatabaseHealthChecking
    private let logger: Logger

    init(
        dbQueue: DatabaseQueue,
        logger: Logger = Logger(label: "Recoto.HealthService")
    ) {
        databaseChecker = GRDBDatabaseHealthChecker(dbQueue: dbQueue)
        self.logger = logger
    }

    init(
        databaseChecker: any DatabaseHealthChecking,
        logger: Logger = Logger(label: "Recoto.HealthService")
    ) {
        self.databaseChecker = databaseChecker
        self.logger = logger
    }

    func check(
        request _: Recoto_Health_V1_CheckRequest,
        context _: ServerContext
    ) async throws -> Recoto_Health_V1_CheckResponse {
        let logger = self.logger.rpc()
        logger.debug("health.check.started")

        let response = await makeCheckResponse()
        logger.info(
            "health.check.finished",
            metadata: ["health_status": .string(response.status)]
        )

        return response
    }

    func makeCheckResponse() async -> Recoto_Health_V1_CheckResponse {
        var databaseResult = Recoto_Health_V1_CheckResult()
        databaseResult.name = "database"

        do {
            try await databaseChecker.check()
            databaseResult.ok = true
            databaseResult.error = ""
        } catch {
            databaseResult.ok = false
            databaseResult.error = String(describing: error)
        }

        var response = Recoto_Health_V1_CheckResponse()
        response.status = databaseResult.ok ? "ok" : "degraded"
        response.results = [databaseResult]
        return response
    }
}
