import Foundation
import Vapor
import Queues

extension Application {
    var radioCore: RadioCore {
        .init(application: self)
    }

    var radio: RadioCore.Radio {
        self.radioCore.radio
    }

    internal struct RadioCore {
        final class Radio {
            var recorders: [Recorder]
            let rootSaveDir: URL

            init() {
                self.recorders = []
                let pwd = FileManager.default.currentDirectoryPath
                let pwdUrl = URL(fileURLWithPath: pwd)
                self.rootSaveDir = pwdUrl.appendingPathComponent("saveData")
                try! FileManager.default.createDirectory(at: self.rootSaveDir, withIntermediateDirectories: true)
            }

            func addRecorder(_ recorder: Recorder) {
                self.recorders.append(recorder)
                let recorderSaveDir = self.rootSaveDir.appendingPathComponent(recorder.saveDirName)
                try! FileManager.default.createDirectory(at: recorderSaveDir, withIntermediateDirectories: true)
            }
        }

        let application: Application
        
        struct Key: StorageKey {
            typealias Value = Radio
        }

        var radio: Radio {
            guard let radio = self.application.storage[Key.self] else {
                fatalError("Core not configured. Configure with app.core.initialize()")
            }
            return radio
        }

        func initialize() {
            self.application.storage[Key.self] = .init()
        }
    }
}
