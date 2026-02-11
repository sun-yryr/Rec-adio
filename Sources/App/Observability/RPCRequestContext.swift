import GRPCCore
import Logging

struct RPCRequestContext: Sendable {
  let traceContext: TraceContext
  let grpcService: String
  let grpcMethod: String
  let grpcMethodFQN: String
  let remotePeer: String
  let localPeer: String

  init(serverContext: ServerContext, traceContext: TraceContext) {
    self.traceContext = traceContext
    self.grpcService = serverContext.descriptor.service.fullyQualifiedService
    self.grpcMethod = serverContext.descriptor.method
    self.grpcMethodFQN = serverContext.descriptor.fullyQualifiedMethod
    self.remotePeer = serverContext.remotePeer
    self.localPeer = serverContext.localPeer
  }

  var traceID: String {
    traceContext.traceID
  }

  var spanID: String {
    traceContext.spanID
  }

  var traceparent: String {
    traceContext.traceparent
  }

  var logMetadata: Logger.Metadata {
    [
      "trace_id": .string(self.traceID),
      "span_id": .string(self.spanID),
      "traceparent": .string(self.traceparent),
      "grpc_service": .string(self.grpcService),
      "grpc_method": .string(self.grpcMethod),
      "grpc_method_fqn": .string(self.grpcMethodFQN),
      "remote_peer": .string(self.remotePeer),
      "local_peer": .string(self.localPeer),
    ]
  }
}
