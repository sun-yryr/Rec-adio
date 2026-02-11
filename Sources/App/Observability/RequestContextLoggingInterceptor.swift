import Foundation
import GRPCCore
import Logging

struct RequestContextLoggingInterceptor: ServerInterceptor {
  let baseLogger: Logger

  func intercept<Input: Sendable, Output: Sendable>(
    request: StreamingServerRequest<Input>,
    context: ServerContext,
    next:
      @Sendable (
        _ request: StreamingServerRequest<Input>,
        _ context: ServerContext
      ) async throws -> StreamingServerResponse<Output>
  ) async throws -> StreamingServerResponse<Output> {
    let traceContext = self.makeTraceContext(from: request.metadata)
    let requestContext = RPCRequestContext(serverContext: context, traceContext: traceContext)

    return try await RPCRequestContextTaskLocal.$current.withValue(requestContext) {
      let logger = self.baseLogger.withRPCContext(requestContext)
      logger.info("rpc.started")

      let startedAt = Date()
      do {
        var response = try await next(request, context)
        self.attachTraceparent(traceContext.traceparent, to: &response)

        logger.info(
          "rpc.finished",
          metadata: [
            "grpc_status": Self.status(from: response),
            "duration_ms": .string(Self.elapsedMilliseconds(since: startedAt)),
          ]
        )
        return response
      } catch var error as RPCError {
        error.metadata.replaceOrAddString(traceContext.traceparent, forKey: TraceContext.traceparentHeader)
        logger.error(
          "rpc.failed",
          metadata: [
            "grpc_status": .string(String(describing: error.code)),
            "duration_ms": .string(Self.elapsedMilliseconds(since: startedAt)),
          ]
        )
        throw error
      } catch {
        logger.error(
          "rpc.failed",
          metadata: [
            "grpc_status": .string("unknown"),
            "duration_ms": .string(Self.elapsedMilliseconds(since: startedAt)),
          ]
        )
        throw error
      }
    }
  }

  private func makeTraceContext(from metadata: Metadata) -> TraceContext {
    if let traceparent = metadata[stringValues: TraceContext.traceparentHeader].first(where: { _ in true }),
      let parsed = TraceContext.parse(traceparent: traceparent)
    {
      return parsed
    }

    return TraceContext.makeNew()
  }

  private func attachTraceparent<Output>(
    _ traceparent: String,
    to response: inout StreamingServerResponse<Output>
  ) {
    var metadata = response.metadata
    metadata.replaceOrAddString(traceparent, forKey: TraceContext.traceparentHeader)
    response.metadata = metadata
  }

  private static func status<Output>(
    from response: StreamingServerResponse<Output>
  ) -> Logger.MetadataValue {
    switch response.accepted {
    case .success:
      return .string("ok")
    case .failure(let error):
      return .string(String(describing: error.code))
    }
  }

  private static func elapsedMilliseconds(since startedAt: Date) -> String {
    let milliseconds = max(0, Int(Date().timeIntervalSince(startedAt) * 1_000))
    return String(milliseconds)
  }
}
