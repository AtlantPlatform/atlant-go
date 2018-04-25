# fs.capnp
@0xe07347b5287484b4;
$import "/go.capnp".package("proto");
$import "/go.capnp".import("proto");
struct ObjectMeta @0xb2b188dc2f537652 {  # 24 bytes, 5 ptrs
  id @0 :Text;  # ptr[0]
  path @1 :Text;  # ptr[1]
  createdAt @2 :Int64;  # bits[0, 64)
  version @3 :Text;  # ptr[2]
  versionPrevious @4 :Text;  # ptr[3]
  isDeleted @5 :Bool;  # bits[64, 65)
  size @6 :Int64;  # bits[128, 192)
  userMeta @7 :Text;  # ptr[4]
}
