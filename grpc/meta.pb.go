/*
Package grpc is a generated protocol buffer package.

It is generated from these files:
	pbrpc.proto

It has these top-level messages:
	RpcMeta
	RpcRequestMeta
	RpcResponseMeta
	ChunkInfo
*/
package grpc

import proto "google.golang.org/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RpcMeta struct {
	Request            *RpcRequestMeta  `protobuf:"bytes,1,opt,name=request" json:"request,omitempty"`
	Response           *RpcResponseMeta `protobuf:"bytes,2,opt,name=response" json:"response,omitempty"`
	CompressType       *int32           `protobuf:"varint,3,opt,name=compress_type,json=compressType" json:"compress_type,omitempty"`
	CorrelationId      *int64           `protobuf:"varint,4,opt,name=correlation_id,json=correlationId" json:"correlation_id,omitempty"`
	AttachmentSize     *int32           `protobuf:"varint,5,opt,name=attachment_size,json=attachmentSize" json:"attachment_size,omitempty"`
	ChunkInfo          *ChunkInfo       `protobuf:"bytes,6,opt,name=chuck_info,json=chuckInfo" json:"chuck_info,omitempty"`
	AuthenticationData []byte           `protobuf:"bytes,7,opt,name=authentication_data,json=authenticationData" json:"authentication_data,omitempty"`
	XXX_unrecognized   []byte           `json:"-"`
}

func (m *RpcMeta) Reset() {
	*m = RpcMeta{}
}

func (m *RpcMeta) String() string {
	return proto.CompactTextString(m)
}

func (*RpcMeta) ProtoMessage() {}

func (*RpcMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{0}
}

func (m *RpcMeta) GetRequest() *RpcRequestMeta {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *RpcMeta) GetResponse() *RpcResponseMeta {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *RpcMeta) GetCompressType() int32 {
	if m != nil && m.CompressType != nil {
		return *m.CompressType
	}
	return 0
}

func (m *RpcMeta) GetCorrelationId() int64 {
	if m != nil && m.CorrelationId != nil {
		return *m.CorrelationId
	}
	return 0
}

func (m *RpcMeta) GetAttachmentSize() int32 {
	if m != nil && m.AttachmentSize != nil {
		return *m.AttachmentSize
	}
	return 0
}

func (m *RpcMeta) GetChuckInfo() *ChunkInfo {
	if m != nil {
		return m.ChunkInfo
	}
	return nil
}

func (m *RpcMeta) GetAuthenticationData() []byte {
	if m != nil {
		return m.AuthenticationData
	}
	return nil
}

type RpcRequestMeta struct {
	ServiceName      *string `protobuf:"bytes,1,req,name=service_name,json=serviceName" json:"service_name,omitempty"`
	MethodName       *string `protobuf:"bytes,2,req,name=method_name,json=methodName" json:"method_name,omitempty"`
	LogId            *int64  `protobuf:"varint,3,opt,name=log_id,json=logId" json:"log_id,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *RpcRequestMeta) Reset() {
	*m = RpcRequestMeta{}
}

func (m *RpcRequestMeta) String() string {
	return proto.CompactTextString(m)
}
func (*RpcRequestMeta) ProtoMessage() {}

func (*RpcRequestMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{1}
}

func (m *RpcRequestMeta) GetServiceName() string {
	if m != nil && m.ServiceName != nil {
		return *m.ServiceName
	}
	return ""
}

func (m *RpcRequestMeta) GetMethodName() string {
	if m != nil && m.MethodName != nil {
		return *m.MethodName
	}
	return ""
}

func (m *RpcRequestMeta) GetLogId() int64 {
	if m != nil && m.LogId != nil {
		return *m.LogId
	}
	return 0
}

type RpcResponseMeta struct {
	ErrorCode        *int32  `protobuf:"varint,1,opt,name=error_code,json=errorCode" json:"error_code,omitempty"`
	ErrorText        *string `protobuf:"bytes,2,opt,name=error_text,json=errorText" json:"error_text,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *RpcResponseMeta) Reset() {
	*m = RpcResponseMeta{}
}

func (m *RpcResponseMeta) String() string {
	return proto.CompactTextString(m)
}

func (*RpcResponseMeta) ProtoMessage() {}

func (*RpcResponseMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{2}
}

func (m *RpcResponseMeta) GetErrorCode() int32 {
	if m != nil && m.ErrorCode != nil {
		return *m.ErrorCode
	}
	return 0
}

func (m *RpcResponseMeta) GetErrorText() string {
	if m != nil && m.ErrorText != nil {
		return *m.ErrorText
	}
	return ""
}

type ChunkInfo struct {
	StreamId         *int64 `protobuf:"varint,1,req,name=stream_id,json=streamId" json:"stream_id,omitempty"`
	ChunkId          *int64 `protobuf:"varint,2,req,name=chunk_id,json=chunkId" json:"chunk_id,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *ChunkInfo) Reset() {
	*m = ChunkInfo{}
}

func (m *ChunkInfo) String() string {
	return proto.CompactTextString(m)
}

func (*ChunkInfo) ProtoMessage() {}

func (*ChunkInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{3}
}

func (m *ChunkInfo) GetStreamId() int64 {
	if m != nil && m.StreamId != nil {
		return *m.StreamId
	}
	return 0
}

func (m *ChunkInfo) GetChunkId() int64 {
	if m != nil && m.ChunkId != nil {
		return *m.ChunkId
	}
	return 0
}

func init() {
	proto.RegisterType((*RpcMeta)(nil), "pbrpc.RpcMeta")
	proto.RegisterType((*RpcRequestMeta)(nil), "pbrpc.RpcRequestMeta")
	proto.RegisterType((*RpcResponseMeta)(nil), "pbrpc.RpcResponseMeta")
	proto.RegisterType((*ChunkInfo)(nil), "pbrpc.ChunkInfo")
}

func init() { proto.RegisterFile("pbrpc.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 388 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0x5f, 0x6f, 0xd3, 0x30,
	0x14, 0xc5, 0x95, 0x94, 0xac, 0xcd, 0x6d, 0xd7, 0x21, 0xa3, 0x21, 0x23, 0x84, 0x08, 0x41, 0x88,
	0x3c, 0xad, 0xd2, 0xbe, 0x42, 0x79, 0xc9, 0x03, 0x20, 0x99, 0xbd, 0x47, 0xc6, 0xbe, 0x5b, 0xa2,
	0x36, 0xb6, 0xb1, 0x6f, 0xd0, 0xb6, 0xef, 0xc7, 0xf7, 0x42, 0x71, 0xba, 0x3f, 0xdd, 0x5b, 0x7c,
	0xce, 0xef, 0xc4, 0xc7, 0xf7, 0xc2, 0xd2, 0xfd, 0xf6, 0x4e, 0x5d, 0x38, 0x6f, 0xc9, 0xb2, 0x2c,
	0x1e, 0xca, 0x7f, 0x29, 0xcc, 0x85, 0x53, 0xdf, 0x91, 0x24, 0xdb, 0xc0, 0xdc, 0xe3, 0x9f, 0x01,
	0x03, 0xf1, 0xa4, 0x48, 0xaa, 0xe5, 0xe5, 0xf9, 0xc5, 0x94, 0x10, 0x4e, 0x89, 0xc9, 0x18, 0x39,
	0xf1, 0x40, 0xb1, 0x4b, 0x58, 0x78, 0x0c, 0xce, 0x9a, 0x80, 0x3c, 0x8d, 0x89, 0xb7, 0xcf, 0x13,
	0x93, 0x13, 0x23, 0x8f, 0x1c, 0xfb, 0x0c, 0xa7, 0xca, 0xf6, 0xce, 0x63, 0x08, 0x0d, 0xdd, 0x39,
	0xe4, 0xb3, 0x22, 0xa9, 0x32, 0xb1, 0x7a, 0x10, 0xaf, 0xee, 0x1c, 0xb2, 0x2f, 0xb0, 0x56, 0xd6,
	0x7b, 0xdc, 0x4b, 0xea, 0xac, 0x69, 0x3a, 0xcd, 0x5f, 0x15, 0x49, 0x35, 0x13, 0xa7, 0xcf, 0xd4,
	0x5a, 0xb3, 0xaf, 0x70, 0x26, 0x89, 0xa4, 0x6a, 0x7b, 0x34, 0xd4, 0x84, 0xee, 0x1e, 0x79, 0x16,
	0xff, 0xb6, 0x7e, 0x92, 0x7f, 0x75, 0xf7, 0xc8, 0x36, 0x00, 0xaa, 0x1d, 0xd4, 0xae, 0xe9, 0xcc,
	0xb5, 0xe5, 0x27, 0xb1, 0xea, 0xeb, 0x43, 0xd5, 0x6d, 0x3b, 0x98, 0x5d, 0x6d, 0xae, 0xad, 0xc8,
	0x23, 0x33, 0x7e, 0xb2, 0x0d, 0xbc, 0x91, 0x03, 0xb5, 0x68, 0xa8, 0x53, 0x53, 0x07, 0x2d, 0x49,
	0xf2, 0x79, 0x91, 0x54, 0x2b, 0xc1, 0x8e, 0xad, 0x6f, 0x92, 0x64, 0xb9, 0x83, 0xf5, 0xf1, 0x94,
	0xd8, 0x27, 0x58, 0x05, 0xf4, 0x7f, 0x3b, 0x85, 0x8d, 0x91, 0x3d, 0xf2, 0xa4, 0x48, 0xab, 0x5c,
	0x2c, 0x0f, 0xda, 0x0f, 0xd9, 0x23, 0xfb, 0x08, 0xcb, 0x1e, 0xa9, 0xb5, 0x7a, 0x22, 0xd2, 0x48,
	0xc0, 0x24, 0x45, 0xe0, 0x1c, 0x4e, 0xf6, 0xf6, 0x66, 0x7c, 0xff, 0x2c, 0xbe, 0x3f, 0xdb, 0xdb,
	0x9b, 0x5a, 0x97, 0x3f, 0xe1, 0xec, 0xc5, 0x80, 0xd9, 0x07, 0x00, 0xf4, 0xde, 0xfa, 0x46, 0x59,
	0x8d, 0x71, 0x7d, 0x99, 0xc8, 0xa3, 0xb2, 0xb5, 0x1a, 0x9f, 0x6c, 0xc2, 0x5b, 0x8a, 0xbb, 0xca,
	0x0f, 0xf6, 0x15, 0xde, 0x52, 0xb9, 0x85, 0xfc, 0x71, 0x0c, 0xec, 0x3d, 0xe4, 0x81, 0x3c, 0xca,
	0x7e, 0xbc, 0x77, 0x6c, 0x3d, 0x13, 0x8b, 0x49, 0xa8, 0x35, 0x7b, 0x07, 0x0b, 0x35, 0x92, 0xa3,
	0x97, 0x46, 0x6f, 0x1e, 0xcf, 0xb5, 0xfe, 0x1f, 0x00, 0x00, 0xff, 0xff, 0xf4, 0x13, 0x42, 0x5e,
	0x5f, 0x02, 0x00, 0x00,
}
