# rs.capnp
@0xfc318f08d018f351;
$import "/go.capnp".package("proto");
$import "/go.capnp".import("proto");
struct Record @0xe3dd34c47de8e014 {  # 8 bytes, 4 ptrs
  id @0 :Text;  # ptr[0]
  path @1 :Text;  # ptr[1]
  createdAt @2 :Int64;  # bits[0, 64)
  current @3 :RecordVersion;  # ptr[2]
  previous @4 :List(RecordVersion);  # ptr[3]
}
struct RecordVersion @0xb498e9811ce1d9a4 {  # 0 bytes, 2 ptrs
  version @0 :Text;  # ptr[0]
  announce @1 :Announce;  # ptr[1]
}
struct Announce @0x9845802f51b21bb9 {  # 16 bytes, 4 ptrs
  id @0 :Text;  # ptr[0]
  nodeID @1 :Text;  # ptr[1]
  signature @2 :Text;  # ptr[2]
  timestamp @3 :Int64;  # bits[0, 64)
  type @4 :AnnounceType;  # bits[64, 80)
  envelope @5 :Data;  # ptr[3]
}
enum AnnounceType @0xaabdfb0036d151b5 {
  unknown @0;
  beatTick @1;
  beatInfo @2;
  recordUpdate @3;
}
struct EnvelopeBeatTick @0x9771146df041e6c1 {  # 0 bytes, 2 ptrs
  id @0 :Text;  # ptr[0]
  session @1 :Text;  # ptr[1]
}
struct EnvelopeBeatInfo @0x9ec9af9924d4017f {  # 24 bytes, 3 ptrs
  id @0 :Text;  # ptr[0]
  session @1 :Text;  # ptr[1]
  ethereumAddr @2 :Text;  # ptr[2]
  uptimeUnix @3 :Int64;  # bits[0, 64)
  inboundWork @4 :UInt64;  # bits[64, 128)
  outboundWork @5 :UInt64;  # bits[128, 192)
}
struct EnvelopeRecordUpdate @0xa55a0b5df4b58f97 {  # 0 bytes, 3 ptrs
  id @0 :Text;  # ptr[0]
  version @1 :Text;  # ptr[1]
  versionPrev @2 :Text;  # ptr[2]
}
