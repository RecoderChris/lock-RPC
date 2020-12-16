struct retType {
  1: i32 cid = 0,
  2: i32 retValue,
}

service lockServe {
    retType acquireLock(1: i32 clientId)
    retType releaseLock(1: i32 clientId)
}