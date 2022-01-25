package grpcx_test

import (
	"github.com/golang/protobuf/proto"
	pbrpc "github.com/layasugar/laya/grpcx/pbrpc"
	"strings"
	"testing"
)

func TestPropertySetAndGet(t *testing.T) {
	var serviceName string = "ThisAServiceName"
	var methodName string = "ThisAMethodName"
	var logId string = "ksadfhksjadhfkjsdhf"

	request := pbrpc.Request{
		ServiceName: serviceName,
		MethodName:  methodName,
		TraceId:     logId,
	}

	if !strings.EqualFold(serviceName, request.GetServiceName()) {
		t.Errorf("set ServiceName value is '%s', but get value is '%s' ", serviceName, request.GetServiceName())
	}

	if !strings.EqualFold(methodName, request.GetMethodName()) {
		t.Errorf("set methodName value is '%s', but get value is '%s' ", methodName, request.GetMethodName())
	}

	if logId != request.GetTraceId() {
		t.Errorf("set logId value is '%s', but get value is '%s' ", logId, request.GetTraceId())
	}

	data, err := proto.Marshal(&request)
	if err != nil {
		t.Errorf("marshaling error: %s", err.Error())
	}

	request2 := new(pbrpc.Request)
	err = proto.Unmarshal(data, request2)
	if err != nil {
		t.Errorf("marshaling error: %s", err.Error())
	}

	if !strings.EqualFold(request.GetServiceName(), request2.GetServiceName()) {
		t.Errorf("set ServiceName value is '%s', but get value is '%s' ", request.GetServiceName(), request2.GetServiceName())
	}
}
