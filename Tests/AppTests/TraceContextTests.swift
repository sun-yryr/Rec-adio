@testable import App
import Testing

@Suite
struct TraceContextTests {
    @Test
    func parseAcceptsValidTraceparentAndNormalizesCase() {
        let parsed = TraceContext.parse(
            traceparent: "00-0AF7651916CD43DD8448EB211C80319C-B7AD6B7169203331-01"
        )

        #expect(parsed != nil)
        #expect(parsed?.version == "00")
        #expect(parsed?.traceID == "0af7651916cd43dd8448eb211c80319c")
        #expect(parsed?.spanID == "b7ad6b7169203331")
        #expect(parsed?.traceFlags == "01")
    }

    @Test
    func parseRejectsInvalidTraceparent() {
        #expect(
            TraceContext.parse(
                traceparent: "00-00000000000000000000000000000000-b7ad6b7169203331-01"
            ) == nil
        )
        #expect(
            TraceContext.parse(
                traceparent: "00-0af7651916cd43dd8448eb211c80319c-0000000000000000-01"
            ) == nil
        )
        #expect(
            TraceContext.parse(
                traceparent: "zz-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
            ) == nil
        )
    }

    @Test
    func makeNewCreatesParseableTraceparent() {
        let context = TraceContext.makeNew()
        let parsed = TraceContext.parse(traceparent: context.traceparent)

        #expect(parsed != nil)
        #expect(parsed?.traceID == context.traceID)
        #expect(parsed?.spanID == context.spanID)
    }
}
