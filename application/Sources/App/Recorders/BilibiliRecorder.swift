import Foundation
import Vapor

struct BilibiliRecorder: Recorder {
    let supportPlatform: String = "bilibili"
    let saveDirName: String

    func run(application: Application, schedule: Schedule) -> Result<Record, RecordingError> {
        guard schedule.extraField != nil else {
            fatalError("require extraField")
        }
        let dateFormatter = DateFormatter.originalFormatter()
        let outputFileName = "\(dateFormatter.string(from: schedule.startDatetime))_\(schedule.title).mp4"
        let outputPath = application.radio.rootSaveDir
            .appendingPathComponent(self.saveDirName)
            .appendingPathComponent(outputFileName)
        let process = Process()
        let pipeErr = Pipe()
        process.standardError = pipeErr
        process.standardOutput = Pipe()
        process.standardInput = Pipe()
        process.executableURL = URL(fileURLWithPath: "/usr/local/bin/yt-dlp")
        process.arguments = [
            "--hls-use-mpegts",
            "-o",
            outputPath.absoluteString,
            schedule.extraField!, // url
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