// swift-tools-version:6.2
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
  name: "Recoto",
  platforms: [.macOS(.v15), .iOS(.v18), .tvOS(.v18)],
  products: [
    .executable(name: "App", targets: ["App"])
  ],
  dependencies: [
    .package(url: "https://github.com/hummingbird-project/hummingbird.git", from: "2.0.0"),
    .package(
      url: "https://github.com/apple/swift-configuration.git", from: "1.0.0",
      traits: [.defaults, "CommandLineArguments"]),
    .package(url: "https://github.com/groue/GRDB.swift.git", from: "7.9.0"),
    .package(url: "https://github.com/grpc/grpc-swift-2.git", from: "2.2.1"),
    .package(url: "https://github.com/grpc/grpc-swift-nio-transport.git", from: "2.4.1"),
    .package(url: "https://github.com/grpc/grpc-swift-protobuf.git", from: "2.1.2"),
    .package(url: "https://github.com/apple/swift-protobuf.git", from: "1.33.3"),
    .package(url: "https://github.com/grpc/grpc-swift-extras.git", from: "2.1.1"),
    .package(url: "https://github.com/swift-server/swift-service-lifecycle.git", from: "2.9.1"),
  ],
  targets: [
    .executableTarget(
      name: "App",
      dependencies: [
        .product(name: "Configuration", package: "swift-configuration"),
        .product(name: "Hummingbird", package: "hummingbird"),
        .product(name: "GRDB", package: "GRDB.swift"),
        .product(name: "GRPCCore", package: "grpc-swift-2"),
        .product(name: "GRPCNIOTransportHTTP2", package: "grpc-swift-nio-transport"),
        .product(name: "GRPCProtobuf", package: "grpc-swift-protobuf"),
        .product(name: "SwiftProtobuf", package: "swift-protobuf"),
        .product(name: "GRPCReflectionService", package: "grpc-swift-extras"),
        .product(name: "ServiceLifecycle", package: "swift-service-lifecycle"),
        .product(name: "UnixSignals", package: "swift-service-lifecycle"),
      ],
      path: "Sources/App",
      resources: [
        .copy("Resources/Reflection/recoto.pb"),
      ]
    ),
    .testTarget(
      name: "AppTests",
      dependencies: [
        .byName(name: "App"),
        .product(name: "HummingbirdTesting", package: "hummingbird"),
      ],
      path: "Tests/AppTests"
    ),
  ]
)
