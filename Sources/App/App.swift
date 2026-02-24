import Configuration
import Hummingbird
import Logging

@main
struct App {
    static func main() async throws {
        // Application will read configuration from the following in the order listed
        // Command line, Environment variables, dotEnv file, defaults provided in memory
        let reader = try await ConfigReader(providers: [
            CommandLineArgumentsProvider(),
            EnvironmentVariablesProvider(),
            EnvironmentVariablesProvider(environmentFilePath: ".env", allowMissing: true),
            InMemoryProvider(values: [
                "http.serverName": "Recoto",
                "grpc.host": "127.0.0.1",
                "grpc.port": "50051",
                "grpc.reflection": "true",
            ]),
        ])
        try await runGRPCDaemon(reader: reader)
    }
}
