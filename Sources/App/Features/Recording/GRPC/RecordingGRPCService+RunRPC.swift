import GRPCCore

extension RecordingGRPCService {
    func cancelRun(request _: Recoto_Recording_V1_CancelRunRequest, context _: ServerContext)
        async throws
        -> Recoto_Recording_V1_CancelRunResponse
    {
        try unimplementedRPC("cancelRun")
    }

    func getRun(request _: Recoto_Recording_V1_GetRunRequest, context _: ServerContext)
        async throws
        -> Recoto_Recording_V1_GetRunResponse
    {
        try unimplementedRPC("getRun")
    }

    func listRuns(request _: Recoto_Recording_V1_ListRunsRequest, context _: ServerContext)
        async throws
        -> Recoto_Recording_V1_ListRunsResponse
    {
        try unimplementedRPC("listRuns")
    }
}
