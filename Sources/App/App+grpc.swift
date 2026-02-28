import Configuration
import Foundation
import GRDB
import GRPCCore
import GRPCNIOTransportHTTP2
import GRPCReflectionService
import Logging
import ServiceLifecycle
import UnixSignals

struct GRPCServerConfig {
    let host: String
    let port: Int
    let reflectionEnabled: Bool
}

private struct GRPCServerRunner: Service {
    let server: GRPCServer<HTTP2ServerTransport.Posix>
    let logger: Logger

    func run() async throws {
        try await withGracefulShutdownHandler {
            try await server.serve()
        } onGracefulShutdown: {
            logger.info("Received graceful shutdown signal. Stopping gRPC server.")
            server.beginGracefulShutdown()
        }
    }
}

enum GRPCConfigurationError: Error, CustomStringConvertible {
    case missingDescriptorSet

    var description: String {
        switch self {
        case .missingDescriptorSet:
            return "Failed to load reflection descriptor set file."
        }
    }
}

func runGRPCDaemon(reader: ConfigReader) async throws {
    let core = try buildAppCore(reader: reader)
    let config = readGRPCServerConfig(reader: reader)
    let server = try buildGRPCServer(config: config, dbQueue: core.dbQueue, logger: core.logger)

    core.logger.info(
        "Starting gRPC server",
        metadata: [
            "host": .string(config.host),
            "port": .string("\(config.port)"),
            "reflection": .string("\(config.reflectionEnabled)"),
        ]
    )

    let serviceGroup = ServiceGroup(
        services: [GRPCServerRunner(server: server, logger: core.logger)],
        gracefulShutdownSignals: [.sigterm, .sigint],
        logger: core.logger
    )
    try await serviceGroup.run()
}

private func readGRPCServerConfig(reader: ConfigReader) -> GRPCServerConfig {
    GRPCServerConfig(
        host: reader.string(forKey: "grpc.host", default: "127.0.0.1"),
        port: reader.int(forKey: "grpc.port", default: 50051),
        reflectionEnabled: reader.bool(forKey: "grpc.reflection", default: true)
    )
}

private func buildGRPCServer(
    config: GRPCServerConfig,
    dbQueue: DatabaseQueue,
    logger: Logger
) throws -> GRPCServer<HTTP2ServerTransport.Posix> {
    let jobRepository = GRDBJobRepository(dbQueue: dbQueue)
    let jobUseCase = DefaultJobUseCase(jobRepository: jobRepository)
    var services: [any RegistrableRPCService] = [
        HealthService(dbQueue: dbQueue, logger: logger),
        RecordingGRPCService(jobUseCase: jobUseCase, logger: logger),
    ]

    if config.reflectionEnabled {
        try services.append(makeReflectionService())
    }

    let transport = HTTP2ServerTransport.Posix(
        address: .ipv4(host: config.host, port: config.port),
        transportSecurity: .plaintext
    )

    return GRPCServer(
        transport: transport,
        services: services,
        interceptors: [RequestContextLoggingInterceptor(baseLogger: logger)]
    )
}

private func makeReflectionService() throws -> ReflectionService {
    guard let descriptorSetPath = Bundle.module.path(forResource: "recoto", ofType: "pb") else {
        throw GRPCConfigurationError.missingDescriptorSet
    }
    return try ReflectionService(descriptorSetFilePaths: [descriptorSetPath])
}
