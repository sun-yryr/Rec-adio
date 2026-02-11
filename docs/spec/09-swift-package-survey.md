# 09. Swift実装向け外部パッケージ調査（Hummingbird前提）

更新日: 2026-02-10

## 1. 前提

- サーバーフレームワークは `Hummingbird` を採用する
- 録音設計は以下を前提:
  - インメモリスケジューラー
  - インメモリキュー
  - SQLite永続化
  - gRPC API
  - CLIクライアント
  - 繰り返しJob（RRULE）

## 2. 採用推奨（コア）

| 用途 | パッケージ | 推奨バージョン方針 | 採用理由 | 参考 |
|---|---|---|---|---|
| HTTPサーバー | `hummingbird-project/hummingbird` | `from: "2.20.0"` | Hummingbird本体。Swift on Server向けでアクティブ。 | https://swiftpackageindex.com/hummingbird-project/hummingbird |
| SQLite | `groue/GRDB.swift` | `from: "7.9.0"` | SQLite運用・マイグレーション・並行アクセス対応が強い。 | https://swiftpackageindex.com/groue/GRDB.swift |
| gRPC core | `grpc/grpc-swift-2` | `from: "2.2.1"` | gRPC Swift v2本体。 | https://swiftpackageindex.com/grpc/grpc-swift-2 |
| gRPC transport | `grpc/grpc-swift-nio-transport` | `from: "2.4.1"` | gRPCのNIO HTTP/2 transport。 | https://swiftpackageindex.com/grpc/grpc-swift-nio-transport |
| gRPC + Protobuf連携 | `grpc/grpc-swift-protobuf` | `from: "2.1.2"` | `.proto` からのサービス連携実装に必要。 | https://swiftpackageindex.com/grpc/grpc-swift-protobuf |
| Protobuf runtime/plugin | `apple/swift-protobuf` | `from: "1.33.3"` | protobufコード生成とランタイム。 | https://swiftpackageindex.com/apple/swift-protobuf |
| 非同期キュー補助 | `apple/swift-async-algorithms` | `from: "1.1.1"` | `AsyncChannel` 等でインメモリキュー実装を簡潔化。 | https://swiftpackageindex.com/apple/swift-async-algorithms |
| ログAPI | `apple/swift-log` | `from: "1.9.1"` | サーバーエコシステム標準。 | https://swiftpackageindex.com/apple/swift-log |
| メトリクスAPI | `apple/swift-metrics` | `from: "2.7.1"` | サーバーエコシステム標準。 | https://swiftpackageindex.com/apple/swift-metrics |
| ライフサイクル管理 | `swift-server/swift-service-lifecycle` | `from: "2.9.1"` | 起動/停止シーケンス管理とgraceful shutdownを標準化。 | https://swiftpackageindex.com/swift-server/swift-service-lifecycle |
| CLI | `apple/swift-argument-parser` | `from: "1.7.0"` | CLI実装の実質標準。 | https://swiftpackageindex.com/apple/swift-argument-parser |

## 3. 採用推奨（設定・運用補助）

| 用途 | パッケージ | 推奨バージョン方針 | 備考 | 参考 |
|---|---|---|---|---|
| TOML設定 | `LebJe/TOMLKit` | `from: "0.6.0"` | `daemon.toml` をSwift側でも継続利用する場合に採用。 | https://swiftpackageindex.com/LebJe/TOMLKit |

## 4. 条件付き採用（RRULE）

| 用途 | パッケージ | 推奨バージョン方針 | 判断 |
|---|---|---|---|
| RFC5545 RRULE | `kubens/RRuleKit` | `from: "1.0.0"` | **要PoC**。実装規模が小さく更新頻度も高くないため、本番採用は評価後に決定。 |

参考:

- https://swiftpackageindex.com/kubens/RRuleKit

補足:

- 初期要件が `FREQ/INTERVAL/BYDAY/BYHOUR/BYMINUTE` に限定されているため、PoC結果が不十分なら自前パーサ実装へ切替可能。

## 5. オブザーバビリティ（任意）

| 用途 | パッケージ | 推奨バージョン方針 | 参考 |
|---|---|---|---|
| Prometheusバックエンド | `swift-server/swift-prometheus` | `from: "2.2.0"` | https://swiftpackageindex.com/swift-server/swift-prometheus |
| OTLPバックエンド | `swift-otel/swift-otel` | `from: "1.0.4"` | https://swiftpackageindex.com/swift-otel/swift-otel |

## 6. Package.swift たたき台

```swift
// swift-tools-version: 6.0
import PackageDescription

let package = Package(
    name: "RecotoSwift",
    platforms: [.macOS("15.0")],
    dependencies: [
        .package(url: "https://github.com/hummingbird-project/hummingbird.git", from: "2.20.0"),
        .package(url: "https://github.com/groue/GRDB.swift.git", from: "7.9.0"),
        .package(url: "https://github.com/grpc/grpc-swift-2.git", from: "2.2.1"),
        .package(url: "https://github.com/grpc/grpc-swift-nio-transport.git", from: "2.4.1"),
        .package(url: "https://github.com/grpc/grpc-swift-protobuf.git", from: "2.1.2"),
        .package(url: "https://github.com/apple/swift-protobuf.git", from: "1.33.3"),
        .package(url: "https://github.com/apple/swift-async-algorithms", from: "1.1.1"),
        .package(url: "https://github.com/apple/swift-log", from: "1.9.1"),
        .package(url: "https://github.com/apple/swift-metrics.git", from: "2.7.1"),
        .package(url: "https://github.com/swift-server/swift-service-lifecycle.git", from: "2.9.1"),
        .package(url: "https://github.com/apple/swift-argument-parser", from: "1.7.0"),
        .package(url: "https://github.com/LebJe/TOMLKit.git", from: "0.6.0"),
        .package(url: "https://github.com/kubens/RRuleKit", from: "1.0.0"),
    ],
    targets: []
)
```

## 7. 採用方針メモ

- まずは「コア」パッケージのみで成立する実装を優先する。
- `RRuleKit` は早期にPoCを行い、問題があれば自前RRULEパーサへ切替える。
- observability backend（Prometheus/OTLP）は初期リリースでは任意導入とする。

