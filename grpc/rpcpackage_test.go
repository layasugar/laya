package grpc_test

import (
	"bytes"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"github.com/layasugar/laya/grpc"
	"strings"
	"testing"
)

//手工定义pb生成的代码, tag 格式 = protobuf:"type,order,req|opt|rep|packed,name=fieldname"
type DataMessage struct {
	Name *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
}

func (m *DataMessage) Reset()         { *m = DataMessage{} }
func (m *DataMessage) String() string { return proto.CompactTextString(m) }
func (*DataMessage) ProtoMessage()    {}

func (m *DataMessage) GetName() string {
	if m.Name != nil {
		return *m.Name
	}
	return ""
}

var sericeName = "thisIsAServiceName"
var methodName = "thisIsAMethodName"
var magicCode = "PRPC"
var logId int64 = 1001
var correlationId int64 = 20001
var data []byte = []byte{1, 2, 3, 1, 2, 3, 1, 1, 2, 2, 20}
var attachment []byte = []byte{2, 2, 2, 2, 2, 1, 1, 1, 1}

func initRpcDataPackage() *grpc.Package {

	rpcDataPackage := grpc.NewPackage()

	rpcDataPackage.SetMagicCode([]byte(magicCode))
	rpcDataPackage.SetData(data)
	rpcDataPackage.SetServiceName(sericeName)
	rpcDataPackage.SetMethodName(methodName)

	rpcDataPackage.SetLogId(logId)
	rpcDataPackage.SetCorrelationId(correlationId)

	rpcDataPackage.SetAttachment(attachment)

	return rpcDataPackage
}

func equalRpcDataPackage(r grpc.Package) error {

	if !strings.EqualFold(sericeName, *r.Meta.Request.ServiceName) {
		return errors.New(fmt.Sprintf("expect serice name '%s' but actual is '%s'", sericeName, *r.Meta.Request.ServiceName))
	}

	if !strings.EqualFold(methodName, *r.Meta.Request.MethodName) {
		return errors.New(fmt.Sprintf("expect method name '%s' but actual is '%s'", methodName, *r.Meta.Request.MethodName))
	}

	if !strings.EqualFold(magicCode, r.GetMagicCode()) {
		return errors.New(fmt.Sprintf("expect magic code '%s' but actual is '%s'", magicCode, r.GetMagicCode()))
	}

	if *r.Meta.Request.LogId != logId {
		return errors.New(fmt.Sprintf("expect logId is '%d' but actual is '%d'", logId, *r.Meta.Request.LogId))
	}

	if *r.Meta.CorrelationId != correlationId {
		return errors.New(fmt.Sprintf("expect CorrelationId is '%d' but actual is '%d'", correlationId, *r.Meta.CorrelationId))
	}

	if !bytes.EqualFold(data, r.Data) {
		return errors.New(fmt.Sprintf("expect data is '%b' but actual is '%b'", data, r.Data))
	}

	if !bytes.EqualFold(attachment, r.Attachment) {
		return errors.New(fmt.Sprintf("expect attachment is '%b' but actual is '%b'", attachment, r.Attachment))
	}

	return nil
}

func validateRpcDataPackage(t *testing.T, r2 grpc.Package) {

	if !strings.EqualFold(magicCode, r2.GetMagicCode()) {
		t.Errorf("expect magic code '%s' but actual is '%s'", magicCode, r2.GetMagicCode())
	}

	if !strings.EqualFold(sericeName, r2.Meta.GetRequest().GetServiceName()) {
		t.Errorf("expect serice name '%s' but actual is '%s'", sericeName, r2.Meta.GetRequest().GetServiceName())
	}

	if !strings.EqualFold(methodName, r2.Meta.GetRequest().GetMethodName()) {
		t.Errorf("expect serice name '%s' but actual is '%s'", sericeName, r2.Meta.GetRequest().GetMethodName())
	}

}

func TestWriteReaderWithMockData(t *testing.T) {

	rpcDataPackage := initRpcDataPackage()

	b, err := rpcDataPackage.Bytes()
	if err != nil {
		t.Error(err.Error())
	}

	r2 := grpc.Package{}

	err = r2.Load(b)
	if err != nil {
		t.Error(err.Error())
	}

	validateRpcDataPackage(t, r2)

}

func WriteReaderWithRealData(rpcDataPackage *grpc.Package,
	compressType int32, t *testing.T) {
	dataMessage := DataMessage{}
	name := "hello, this is repeated string aaaaaaaaaaaaaaaaaaaaaa"
	dataMessage.Name = &name

	data, err := proto.Marshal(&dataMessage)
	if err != nil {
		t.Error(err.Error())
	}
	rpcDataPackage.SetData(data)

	b, err := rpcDataPackage.Bytes()
	if err != nil {
		t.Error(err.Error())
	}

	r2 := grpc.Package{}
	r2.SetCompressType(compressType)

	err = r2.Load(b)
	if err != nil {
		t.Error(err.Error())
	}

	validateRpcDataPackage(t, r2)

	dataMessage2 := DataMessage{}
	proto.Unmarshal(r2.Data, &dataMessage2)

	if !strings.EqualFold(name, *dataMessage2.Name) {
		t.Errorf("expect name '%s' but actual is '%s'", name, *dataMessage2.Name)
	}
}

func TestWriteReaderWithRealData(t *testing.T) {

	rpcDataPackage := initRpcDataPackage()
	WriteReaderWithRealData(rpcDataPackage, grpc.COMPRESS_NO, t)
}

func TestWriteReaderWithGZIP(t *testing.T) {

	rpcDataPackage := initRpcDataPackage()

	rpcDataPackage.SetCompressType(grpc.COMPRESS_GZIP)

	WriteReaderWithRealData(rpcDataPackage, grpc.COMPRESS_GZIP, t)

}

func TestWriteReaderWithSNAPPY(t *testing.T) {

	rpcDataPackage := initRpcDataPackage()

	rpcDataPackage.SetCompressType(grpc.COMPRESS_SNAPPY)

	WriteReaderWithRealData(rpcDataPackage, grpc.COMPRESS_SNAPPY, t)

}
