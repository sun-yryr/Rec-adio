import Foundation

struct TraceContext: Sendable, Equatable {
    static let traceparentHeader = "traceparent"
    static let supportedVersion = "00"
    static let sampledFlags = "01"

    private static let traceIDLength = 32
    private static let spanIDLength = 16
    private static let traceFlagsLength = 2
    private static let hexAlphabet = Array("0123456789abcdef".utf8)

    let version: String
    let traceID: String
    let spanID: String
    let traceFlags: String

    var traceparent: String {
        "\(version)-\(traceID)-\(spanID)-\(traceFlags)"
    }

    static func parse(traceparent value: String) -> TraceContext? {
        let components = value
            .trimmingCharacters(in: .whitespacesAndNewlines)
            .split(separator: "-", omittingEmptySubsequences: false)

        guard components.count == 4 else {
            return nil
        }

        let version = String(components[0]).lowercased()
        let traceID = String(components[1]).lowercased()
        let spanID = String(components[2]).lowercased()
        let traceFlags = String(components[3]).lowercased()

        guard version != "ff", Self.isHex(version, length: 2) else {
            return nil
        }
        guard Self.isHex(traceID, length: Self.traceIDLength), !Self.isAllZero(traceID) else {
            return nil
        }
        guard Self.isHex(spanID, length: Self.spanIDLength), !Self.isAllZero(spanID) else {
            return nil
        }
        guard Self.isHex(traceFlags, length: Self.traceFlagsLength) else {
            return nil
        }

        return TraceContext(
            version: version,
            traceID: traceID,
            spanID: spanID,
            traceFlags: traceFlags
        )
    }

    static func makeNew() -> TraceContext {
        TraceContext(
            version: supportedVersion,
            traceID: randomHex(byteCount: 16),
            spanID: randomHex(byteCount: 8),
            traceFlags: sampledFlags
        )
    }

    private static func randomHex(byteCount: Int) -> String {
        var bytes: [UInt8] = []
        bytes.reserveCapacity(byteCount * 2)
        var generator = SystemRandomNumberGenerator()

        for _ in 0 ..< byteCount {
            let value = UInt8.random(in: UInt8.min ... UInt8.max, using: &generator)
            bytes.append(Self.hexAlphabet[Int(value >> 4)])
            bytes.append(Self.hexAlphabet[Int(value & 0x0F)])
        }

        guard let hex = String(bytes: bytes, encoding: .utf8) else {
            preconditionFailure("randomHex produced non-UTF8 bytes")
        }
        return hex
    }

    private static func isHex(_ value: String, length: Int) -> Bool {
        guard value.utf8.count == length else {
            return false
        }

        return value.utf8.allSatisfy { byte in
            switch byte {
            case UInt8(ascii: "0") ... UInt8(ascii: "9"),
                 UInt8(ascii: "a") ... UInt8(ascii: "f"):
                return true
            default:
                return false
            }
        }
    }

    private static func isAllZero(_ value: String) -> Bool {
        value.utf8.allSatisfy { $0 == UInt8(ascii: "0") }
    }
}
