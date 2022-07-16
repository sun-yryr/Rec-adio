import Fluent
import FluentMySQLDriver
import Vapor
import QueuesRedisDriver

// configures your application
public func configure(_ app: Application) throws {
    // uncomment to serve files from /Public folder
    // app.middleware.use(FileMiddleware(publicDirectory: app.directory.publicDirectory))

    var tlsConfigure: TLSConfiguration = .makeClientConfiguration()
    tlsConfigure.certificateVerification = .none
    app.databases.use(.mysql(
        hostname: Environment.get("DATABASE_HOST") ?? "localhost",
        port: Environment.get("DATABASE_PORT").flatMap(Int.init(_:)) ?? MySQLConfiguration.ianaPortNumber,
        username: Environment.get("DATABASE_USERNAME") ?? "vapor_username",
        password: Environment.get("DATABASE_PASSWORD") ?? "vapor_password",
        database: Environment.get("DATABASE_NAME") ?? "vapor_database",
        tlsConfiguration: tlsConfigure
    ), as: .mysql)

    app.migrations.add(CreateProgramInfo())
    app.migrations.add(CreateSchedule())
    app.migrations.add(CreateRecord())
    app.migrations.add(CreateTwitterUsers())
    app.migrations.add(CreateBilibiliUsers())
    app.migrations.add(AddDeletedAtScheduleModel())

    // commands
    app.commands.use(SampleCommand(), as: "sample")
    app.commands.use(AddScheduleCommand(), as: "schedule:add")
    app.commands.use(AddTwitterUserCommand(), as: "twitter_user:add")
    app.commands.use(AddBilibiliUserCommand(), as: "bilibili_user:add")

    // radio
    app.radioCore.initialize()
    app.radio.addRecorder(AgqrRecorder(saveDirName: "agqr", streamingUrl: "https://fms2.uniqueradio.jp/agqr10/aandg1.m3u8"))
    app.radio.addRecorder(SpaceRecorder(saveDirName: "space"))
    app.radio.addRecorder(BilibiliRecorder(saveDirName: "bilibili"))
    
    // job queue
    try app.queues.use(.redis(url: Environment.get("REDIS_URL")!))
    app.queues.schedule(RecordingScheduleJob())
        .minutely()
        .at(0)
    app.queues.schedule(SearchSpaceScheduleJob(bearerToken: Environment.get("TWITTER_BEARER_TOKEN")!))
        .minutely()
        .at(0)
    app.queues.schedule(SearchBilibiliScheduleJob())
        .minutely()
        .at(0)
    app.queues.add(RecordingJob())

    // register routes
    try routes(app)
}
