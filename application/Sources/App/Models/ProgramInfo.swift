import FluentKit

final class ProgramInfo: Model {
    static var schema: String = "program_infos"
    init() {}

    @ID(custom: .id, generatedBy: .database)
    var id: Int?

    @Field(key: "title")
    var title: String

    @Field(key: "platform")
    var platform: String

    @Field(key: "pfm")
    var pfm: String

    @Field(key: "program_url")
    var programUrl: String

    init(id: Int? = nil, title: String, platform: String, pfm: String, programUrl: String) {
        self.id = id
        self.title = title
        self.platform = platform
        self.pfm = pfm
        self.programUrl = programUrl
    }
}