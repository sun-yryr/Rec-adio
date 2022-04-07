import Vapor

struct SampleCommand: Command {
    struct Signature: CommandSignature {}

    var help: String = "sample command"

    func run(using context: CommandContext, signature: Signature) throws {
        let cmd = Process()
        cmd.executableURL = URL(fileURLWithPath: "/usr/bin/pwd")
        cmd.currentDirectoryURL = URL(fileURLWithPath: "/app", isDirectory: true)
        let pipe = Pipe()
        cmd.standardOutput = pipe
        do {
            try cmd.run()
            let readHandle = pipe.fileHandleForReading
            let data = try readHandle.readToEnd()
            context.console.print(String(data: data!, encoding: .utf8)!, newLine: true)
        } catch {
            context.console.error(error.localizedDescription, newLine: true)
        }
    }
}
