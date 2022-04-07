import FluentKit
import Foundation

final class Record: Model {
    static var schema: String = "records"
    init() {}

    @ID(custom: .id, generatedBy: .database)
    var id: Int?

    @Field(key: "title")
    var title: String

    @Field(key: "record_datetime")
    var recordDatetime: Date

    @Field(key: "duration")
    var duration: Int

    @Field(key: "path")
    var path: String

    @Field(key: "platform")
    var platform: String

    @Parent(key: "program_info_id")
    var programInfo: ProgramInfo

    init(id: Int? = nil, title: String, recordDatetime: Date, duration: Int, path: String, platform: String, programInfoId: ProgramInfo.IDValue) {
        self.id = id
        self.title = title
        self.recordDatetime = recordDatetime
        self.duration = duration
        self.path = path
        self.platform = platform
        self.$programInfo.id = programInfoId
    }

    convenience init(schedule: Schedule, path: String, recordDatetime: Date = Date()) {
        self.init(title: schedule.title, recordDatetime: recordDatetime, duration: schedule.duration, path: path, platform: schedule.platform, programInfoId: schedule.$programInfo.id)
    }
}