import Foundation
import Vapor

struct AgqrRecorder: Recorder {
    let supportPlatform: String = "agqr"
    let saveDirName: String
    let streamingUrl: String

    init(saveDirName: String, streamingUrl: String) {
        self.saveDirName = saveDirName
        self.streamingUrl = streamingUrl
    }

    func run(application: Application, schedule: Schedule) -> Result<Record, RecordingError> {
        let dateFormatter = DateFormatter.originalFormatter()
        let outputFileName = "\(dateFormatter.string(from: schedule.startDatetime))_\(schedule.title).mp3"
        let outputPath = application.radio.rootSaveDir
            .appendingPathComponent(self.saveDirName)
            .appendingPathComponent(outputFileName)
        let process = Process()
        let pipeErr = Pipe()
        process.standardError = pipeErr
        process.standardOutput = Pipe()
        process.standardInput = Pipe()
        process.executableURL = URL(fileURLWithPath: "/usr/bin/ffmpeg")
        process.arguments = [
            "-i",
            self.streamingUrl,
            "-t",
            String(schedule.duration * 60),
            outputPath.absoluteString,
        ]
        do {
            try process.run()
        } catch {
            return .failure(.processError(error.localizedDescription))
        }
        process.waitUntilExit()
        if (process.terminationStatus != 0) {
            let stdErrData = try! pipeErr.fileHandleForReading.readToEnd()
            let stdErr = String(data: stdErrData ?? Data(), encoding: .utf8) ?? "nil"
            return .failure(.processError(stdErr))
        }
        return .success(Record(schedule: schedule, path: outputPath.absoluteString, recordDatetime: Date()))
    }
}
