import Logging

extension Logger {
    func withRPCContext(_ context: RPCRequestContext) -> Logger {
        var logger = self
        for (key, value) in context.logMetadata {
            logger[metadataKey: key] = value
        }
        return logger
    }

    func rpc() -> Logger {
        guard let context = RPCRequestContextTaskLocal.current else {
            return self
        }
        return withRPCContext(context)
    }
}
