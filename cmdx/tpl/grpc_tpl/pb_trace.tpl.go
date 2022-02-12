package grpc_tpl

const PbTraceTpl = `syntax = "proto3";

package pb;

option go_package = ".;pb";

service Greeter {
  rpc SayHello(HelloRequest) returns (HelloReply) {}
  rpc GrpcTraceTest(GrpcTraceTestReq) returns (HelloReply) {}
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}

message GrpcTraceTestReq {
  int32 kind = 1;
}
`
