import Foundation
import FluentKit

final class Schedule: Model {
    static let schema = "schedules"
    init() { }

    @ID(custom: .id, generatedBy: .database)
    var id: Int?

    @Field(key: "title")
    var title: String

    @Field(key: "start_datetime")
    var startDatetime: Date
    
    @Field(key: "duration")
    var duration: Int

    @Field(key: "is_processing")
    var isProcessing: Bool

    @Field(key: "platform")
    var platform: String

    @Parent(key: "program_info_id")
    var programInfo: ProgramInfo

    @OptionalField(key: "extra_field")
    var extraField: String?

    @Timestamp(key: "deleted_at", on: .delete)
    var deletedAt: Date?
    
    init(id: Int? = nil, title: String, startDatetime: Date, duration: Int, isProcessing: Bool = false, platform: String, programInfoId: ProgramInfo.IDValue, extraField: String? = nil) {
        self.id = id
        self.title = title
        self.startDatetime = startDatetime
        self.duration = duration
        self.isProcessing = isProcessing
        self.platform = platform
        self.$programInfo.id = programInfoId
        self.extraField = extraField
    }
}
