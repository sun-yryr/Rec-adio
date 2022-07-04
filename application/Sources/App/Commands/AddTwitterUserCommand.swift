import Vapor

struct AddTwitterUserCommand: Command {
    struct Signature: CommandSignature {
        @Argument(name: "name")
        var name: String

        @Argument(name: "userId")
        var userId: String
    }

    var help = "add data to twitter_user table"

    func run(using context: CommandContext, signature: Signature) throws {
        var info = try ProgramInfo.query(on: context.application.db)
            .filter(\.$platform, .equal, "space")
            .filter(\.$pfm, .equal, signature.name)
            .first()
            .wait()
        if (info == nil) {
            info = ProgramInfo(
                title: "space_\(signature.name)",
                platform: "space",
                pfm: signature.name,
                programUrl: ""
            )
            try info!.create(on: context.application.db).wait()
        }
        try TwitterUser(
            name: signature.name,
            userId: signature.userId,
            programInfoId: info!.id!
        ).create(on: context.application.db).wait()
    }
}