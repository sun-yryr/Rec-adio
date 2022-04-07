import Vapor

struct AddScheduleCommand: Command {
    struct Signature: CommandSignature {
        @Argument(name: "title")
        var title: String

        @Argument(name: "platform")
        var platform: String

        @Argument(name: "startAt")
        var startAt: String

        @Argument(name: "duration")
        var duration: Int
    }

    var help: String = "add data to schedule table"

    func run(using context: CommandContext, signature: Signature) throws {
        var info = try ProgramInfo.query(on: context.application.db)
            .filter(\.$title, .equal, signature.title)
            .filter(\.$platform, .equal, signature.platform)
            .first()
            .wait()
        if (info == nil) {
            // 新規作成
            info = ProgramInfo(
                title: signature.title, 
                platform: signature.platform, 
                pfm: "", 
                programUrl: ""
            )
            try info!.create(on: context.application.db).wait()
        }
        let df = DateFormatter()
        df.dateFormat = "yyyy-MM-dd'T'HH:mm:ssXXX"
        df.locale = Locale(identifier: "en_US_POSIX")
        try Schedule(
            title: signature.title,
            startDatetime: df.date(from: signature.startAt)!, 
            duration: signature.duration, 
            platform: signature.platform, 
            programInfoId: info!.id!
        ).create(on: context.application.db).wait()
    }
}
