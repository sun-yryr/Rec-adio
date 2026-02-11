import GRPCCore
import Logging
import Testing

@testable import App

@Suite
struct RequestContextLoggingInterceptorTests {
  @Test
  func interceptorUsesIncomingTraceparentAndPropagatesItToResponse() async throws {
    let incomingTraceparent = "00-0AF7651916CD43DD8448EB211C80319C-B7AD6B7169203331-01"
    var metadata: Metadata = [:]
    metadata.addString(incomingTraceparent, forKey: TraceContext.traceparentHeader)

    let request = StreamingServerRequest(
      single: ServerRequest(metadata: metadata, message: "ping")
    )
    let interceptor = RequestContextLoggingInterceptor(baseLogger: Logger(label: "test"))

    let response = try await interceptor.intercept(
      request: request,
      context: makeServerContext()
    ) { _, _ in
      #expect(
        RPCRequestContextTaskLocal.current?.traceID == "0af7651916cd43dd8448eb211c80319c"
      )
      #expect(
        RPCRequestContextTaskLocal.current?.spanID == "b7ad6b7169203331"
      )
      return StreamingServerResponse(of: String.self) { _ in [:] }
    }

    let responseTraceparent = response.metadata[stringValues: TraceContext.traceparentHeader]
      .first(where: { _ in true })

    #expect(
      responseTraceparent == "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
    )
  }

  @Test
  func interceptorGeneratesTraceparentWhenMissingOrInvalid() async throws {
    var metadata: Metadata = [:]
    metadata.addString("invalid-traceparent", forKey: TraceContext.traceparentHeader)

    let request = StreamingServerRequest(
      single: ServerRequest(metadata: metadata, message: "ping")
    )
    let interceptor = RequestContextLoggingInterceptor(baseLogger: Logger(label: "test"))

    let response = try await interceptor.intercept(
      request: request,
      context: makeServerContext()
    ) { _, _ in
      #expect(RPCRequestContextTaskLocal.current != nil)
      return StreamingServerResponse(of: String.self) { _ in [:] }
    }

    let responseTraceparent = response.metadata[stringValues: TraceContext.traceparentHeader]
      .first(where: { _ in true })

    #expect(responseTraceparent != nil)
    #expect(responseTraceparent != "invalid-traceparent")
    #expect(TraceContext.parse(traceparent: responseTraceparent ?? "") != nil)
  }
}

private func makeServerContext() -> ServerContext {
  ServerContext(
    descriptor: MethodDescriptor(
      fullyQualifiedService: "recoto.health.v1.HealthService",
      method: "Check"
    ),
    remotePeer: "ipv4:127.0.0.1:50000",
    localPeer: "ipv4:127.0.0.1:50051",
    cancellation: .init()
  )
}
