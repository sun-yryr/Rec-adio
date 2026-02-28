import GRPCCore
import Logging

struct RecordingGRPCService: Recoto_Recording_V1_RecordingService.SimpleServiceProtocol {
    let logger: Logger
    let jobUseCase: any JobUseCase

    init(
        jobUseCase: any JobUseCase,
        logger: Logger = Logger(label: "Recoto.RecordingGRPCService")
    ) {
        self.jobUseCase = jobUseCase
        self.logger = logger
    }

    init(jobRepo: any JobRepository, logger: Logger = Logger(label: "Recoto.RecordingGRPCService")) {
        self.init(jobUseCase: DefaultJobUseCase(jobRepository: jobRepo), logger: logger)
    }

    func unimplementedRPC<Response>(_ method: String) throws -> Response {
        throw RPCError(
            code: .unimplemented,
            message: "\(method) is not implemented yet"
        )
    }
}
