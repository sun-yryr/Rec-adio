import Vapor

struct WebController: RouteCollection {
    func boot(routes: RoutesBuilder) throws {
        routes.get("hello", use: hello(req:))
    }

    func hello(req: Request) throws -> EventLoopFuture<View> {
        return req.view.render("hello", ["name": "Leaf"])
    }
}
