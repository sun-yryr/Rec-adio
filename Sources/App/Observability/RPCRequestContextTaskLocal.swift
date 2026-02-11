enum RPCRequestContextTaskLocal {
  @TaskLocal static var current: RPCRequestContext?
}
