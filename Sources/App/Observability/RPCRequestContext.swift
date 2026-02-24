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
        grpcService = serverContext.descriptor.service.fullyQualifiedService
        grpcMethod = serverContext.descriptor.method
        grpcMethodFQN = serverContext.descriptor.fullyQualifiedMethod
        remotePeer = serverContext.remotePeer
        localPeer = serverContext.localPeer
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
            "trace_id": .string(traceID),
            "span_id": .string(spanID),
            "traceparent": .string(traceparent),
            "grpc_service": .string(grpcService),
            "grpc_method": .string(grpcMethod),
            "grpc_method_fqn": .string(grpcMethodFQN),
            "remote_peer": .string(remotePeer),
            "local_peer": .string(localPeer),
        ]
    }
}
