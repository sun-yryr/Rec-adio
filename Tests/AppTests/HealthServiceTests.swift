import Testing

@testable import App

@Suite
struct HealthServiceTests {
  @Test
  func checkReturnsOkWhenDatabaseIsReachable() async {
    let service = HealthService(databaseChecker: SuccessfulDatabaseChecker())

    let response = await service.makeCheckResponse()

    #expect(response.status == "ok")
    #expect(response.results.count == 1)
    #expect(response.results[0].name == "database")
    #expect(response.results[0].ok)
    #expect(response.results[0].error.isEmpty)
  }

  @Test
  func checkReturnsDegradedWhenDatabaseIsUnreachable() async {
    let service = HealthService(databaseChecker: FailingDatabaseChecker())

    let response = await service.makeCheckResponse()

    #expect(response.status == "degraded")
    #expect(response.results.count == 1)
    #expect(response.results[0].name == "database")
    #expect(!response.results[0].ok)
    #expect(response.results[0].error.contains("database down"))
  }
}

private struct SuccessfulDatabaseChecker: DatabaseHealthChecking {
  func check() async throws {}
}

private struct FailingDatabaseChecker: DatabaseHealthChecking {
  func check() async throws {
    throw DatabaseDownError()
  }
}

private struct DatabaseDownError: Error, CustomStringConvertible {
  var description: String { "database down" }
}
