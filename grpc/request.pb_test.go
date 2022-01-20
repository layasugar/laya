package grpc_test

import (
	"google.golang.org/protobuf/proto"
	"github.com/layasugar/laya/grpc"
	"strings"
	"testing"
)

func TestPropertySetAndGet(t *testing.T) {

	if true {
		return
	}

	var serviceName string = "ThisAServiceName"
	var methodName string = "ThisAMethodName"
	var logId int64 = 1

	request := grpc.RpcRequestMeta{
		ServiceName: &serviceName,
		MethodName:  &methodName,
		LogId:       &logId,
	}

	if !strings.EqualFold(serviceName, request.GetServiceName()) {
		t.Errorf("set ServiceName value is '%s', but get value is '%s' ", serviceName, request.GetServiceName())
	}

	if !strings.EqualFold(methodName, request.GetMethodName()) {
		t.Errorf("set methodName value is '%s', but get value is '%s' ", methodName, request.GetMethodName())
	}

	if logId != request.GetLogId() {
		t.Errorf("set logId value is '%d', but get value is '%d' ", logId, request.GetLogId())
	}

	data, err := proto.Marshal(&request)
	if err != nil {
		t.Errorf("marshaling error: %s", err.Error())
	}

	request2 := new(grpc.RpcRequestMeta)
	err = proto.Unmarshal(data, request2)
	if err != nil {
		t.Errorf("marshaling error: %s", err.Error())
	}

	if !strings.EqualFold(request.GetServiceName(), request2.GetServiceName()) {
		t.Errorf("set ServiceName value is '%s', but get value is '%s' ", request.GetServiceName(), request2.GetServiceName())
	}
}
