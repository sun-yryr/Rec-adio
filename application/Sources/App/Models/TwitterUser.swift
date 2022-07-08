import FluentKit

final class TwitterUser: Model {
    static var schema: String = "twitter_users"
    init() {}

    @ID(custom: .id, generatedBy: .database)
    var id: Int?

    @Field(key: "name")
    var name: String

    @Field(key: "user_id")
    var userId: String

    @Parent(key: "program_info_id")
    var programInfo: ProgramInfo

    init(id: Int? = nil, name: String, userId: String, programInfoId: ProgramInfo.IDValue) {
        self.id = id
        self.name = name
        self.userId = userId
        self.$programInfo.id = programInfoId
    }
}
