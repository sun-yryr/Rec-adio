import Foundation
import Vapor

protocol Recorder {
    var supportPlatform: String { get }
    var saveDirName: String { get }
    func run(application: Application, schedule: Schedule) -> Result<Record, RecordingError>
}

enum RecordingError: Error {
    case defaultError
    case processError(String)
    case NotFoundPlatform(String)
}
