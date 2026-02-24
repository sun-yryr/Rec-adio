import Configuration
import Foundation
import GRDB
import Hummingbird
import JSONLogger
import Logging

/// Request context used by application
typealias AppRequestContext = BasicRequestContext

struct AppCore {
    let logger: Logger
    let dbQueue: DatabaseQueue
    let jobRepository: any JobRepository
}

///  Build application
/// - Parameter reader: configuration reader
func buildApplication(reader: ConfigReader) async throws -> some ApplicationProtocol {
    let core = try buildAppCore(reader: reader)
    let router = try buildRouter(jobRepository: core.jobRepository)
    return Application(
        router: router,
        configuration: ApplicationConfiguration(reader: reader.scoped(to: "http")),
        logger: core.logger
    )
}

func buildAppCore(reader: ConfigReader) throws -> AppCore {
    let logger = makeLogger(reader: reader)
    let dbQueue = try makeDatabaseQueue(reader: reader)
    let migrator = makeDatabaseMigrator()
    try migrator.migrate(dbQueue)
    return AppCore(
        logger: logger,
        dbQueue: dbQueue,
        jobRepository: GRDBJobRepository(dbQueue: dbQueue)
    )
}

private func makeLogger(reader: ConfigReader) -> Logger {
    LoggingSystem.bootstrap(JSONLogger.init, metadataProvider: nil)
    var logger = Logger(label: "Recoto")
    logger.logLevel = reader.string(forKey: "log.level", as: Logger.Level.self, default: .info)
    return logger
}

func makeDatabaseQueue(reader: ConfigReader) throws -> DatabaseQueue {
    let dbPath = reader.string(forKey: "db.path", default: "./data/database.sqlite")
    if dbPath != ":memory:" {
        let directoryPath = (dbPath as NSString).deletingLastPathComponent
        if !directoryPath.isEmpty, directoryPath != "." {
            try FileManager.default.createDirectory(
                atPath: directoryPath,
                withIntermediateDirectories: true
            )
        }
    }
    return try DatabaseQueue(path: dbPath)
}

/// Build router
func buildRouter(jobRepository _: any JobRepository) throws -> Router<AppRequestContext> {
    let router = Router(context: AppRequestContext.self)
    // Add middleware
    router.addMiddleware {
        // logging middleware
        LogRequestsMiddleware(.info)
    }
    // Add default endpoint
    router.get("/") { _, _ in
        "Hello!"
    }
    return router
}
