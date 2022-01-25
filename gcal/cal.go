// Package cal 提供了一个支持多种交互协议和打包格式的扩展包。RAL规定了
// 一套高度抽象的交互过程规范，将整个后端交互过程分成了交互协议和数据打
// 包/解包两大块，可以支持一些常用的后端交互协议，标准化协议扩充的开发过程，
// 促进代码复用。
//
package cal

import (
	"fmt"
	"reflect"
	"strings"

	"gitlab.xthktech.cn/bs/gxe/cal/balance"
	"gitlab.xthktech.cn/bs/gxe/cal/context"
	"gitlab.xthktech.cn/bs/gxe/cal/converter"
	"gitlab.xthktech.cn/bs/gxe/cal/protocol"
	"gitlab.xthktech.cn/bs/gxe/cal/service"
)

// Cal 发送网络请求，并对象化返回数据
// request 为interface{} 是一个不好的设计，但又难以抽象出接口，只能做类型强转
// response 本身是个指针，会尝试根据 convertType 去赋值
// Cal 是个快捷函数，完成如下工作：
// 1 找到服务发现对象 S
// 2 找到负载均衡策略B 找到要访问的主机
// 3 将请求对象序列化
// 4 发送请求
// 5 将返回数据对象化
func Cal(serviceName string, request interface{}, response interface{}, converterType ConverterType) (err error) {
	ctx := context.NewContext()
	ctx.ServiceName = serviceName
	ctx.Caller = "CAL"
	serv, _ := service.GetService(serviceName)
	if serv == nil {
		err := fmt.Errorf("can't find service %s, may be load failed", serviceName)
		ctx.CurRecord().Error = err
		ctx.FlushLog()
		return err
	}

	return calWithService(ctx, serv, request, response, converterType)
}

// calWithService 跳过service查找过程
func calWithService(ctx *context.Context, serv service.Service, request interface{}, response interface{}, converterType ConverterType) (err error) {
	defer ctx.FlushLog()
	ctx.TimeStatisStart("cost")
	ctx.ServiceName = serv.GetName()
	bala, _ := balance.GetBalance(serv.GetStrategy())
	if bala == nil {
		err = fmt.Errorf("can't find balance %s, use default strategy", serv.GetStrategy())
		ctx.CurRecord().Error = err
		bala = balance.GetDefaultBalance()
	}

	ctx.BalanceName = bala.Name()

	retry := serv.GetRetry()
	if retry < 0 {
		retry = 0
	}

	ctx.MaxTry = retry

	var addr *service.Addr
	var rsp *protocol.Response
	for i := 0; i < retry+1; i++ {
		proto, err := protocol.NewProtocol(ctx, serv, request)
		if err != nil {
			ctx.CurRecord().Error = err
			return err
		}
		ctx.Protocol = proto.Protocol()

		if i > 1 {
			ctx.TimeStatisStart("cost")
		}
		ctx.CurRecord().RecordTimePoint("req_start_time")
		addr, err = bala.FetchServer(serv)
		if err != nil {
			ctx.CurRecord().Error = err
			ctx.TimeStatisStop("cost")
			return err
		}

		ctx.CurRecord().RecordTimePoint("talk_start_time")
		rsp, err = proto.Do(ctx, addr)
		ctx.TimeStatisStop("cost")
		if err == nil {
			break
		}
		ctx.CurRecord().Error = err
		ctx.NextRecord()
	}

	// 请求失败了
	if err != nil {
		return
	}
	err = valueRsp(ctx, rsp, converterType, response)

	return
}

func valueRsp(ctx *context.Context, calResult *protocol.Response, converterType converter.ConverterType, rsp interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	rval := reflect.ValueOf(rsp).Elem()
	rtype := reflect.TypeOf(rsp).Elem()
	num := rtype.NumField()

	headField := rval.FieldByName("Head")
	if headField.CanSet() {
		headField.Set(reflect.ValueOf(calResult.Head))
	}
	for i := 0; i < num; i++ {
		if strings.ToLower(rtype.Field(i).Tag.Get("cal")) == "head" {
			headField = rval.Field(i)
			if headField.CanSet() {
				headField.Set(reflect.ValueOf(calResult.Head))
			}
			break
		}
	}

	requestField := rval.FieldByName("Request")
	if requestField.CanSet() {
		requestField.Set(reflect.ValueOf(calResult.Request))
	}
	for i := 0; i < num; i++ {
		if strings.ToLower(rtype.Field(i).Tag.Get("cal")) == "request" {
			requestField = rval.Field(i)
			if requestField.CanSet() {
				requestField.Set(reflect.ValueOf(calResult.Request))
			}
			break
		}
	}

	bodyField := rval.FieldByName("Body")
	if bodyField.CanSet() {
		conv, _ := converter.GetConverter(converterType)
		if conv == nil {
			return fmt.Errorf("bad rsp convert type: %s", converterType)
		}

		b := bodyField.Addr().Interface()
		ctx.TimeStatisStart("unpack")
		err = conv.UnPack(calResult.Body, b)
		ctx.TimeStatisStop("unpack")
		if err != nil {
			return err
		}
	}
	for i := 0; i < num; i++ {
		if strings.ToLower(rtype.Field(i).Tag.Get("cal")) == "body" {
			bodyField = rval.Field(i)
			if bodyField.CanSet() {
				conv, _ := converter.GetConverter(converterType)
				if conv == nil {
					return fmt.Errorf("bad rsp convert type: %s", converterType)
				}

				b := bodyField.Addr().Interface()
				ctx.TimeStatisStart("unpack")
				err = conv.UnPack(calResult.Body, b)
				ctx.TimeStatisStop("unpack")
				if err != nil {
					return err
				}
			}
			break
		}
	}

	rawField := rval.FieldByName("Raw")
	if rawField.CanSet() {
		rawField.Set(reflect.ValueOf(calResult.Body))
	}
	for i := 0; i < num; i++ {
		if strings.ToLower(rtype.Field(i).Tag.Get("cal")) == "raw" {
			rawField = rval.Field(i)
			if rawField.CanSet() {
				rawField.Set(reflect.ValueOf(calResult.Body))
			}
			break
		}
	}

	originRspField := rval.FieldByName("OriginRsp")
	if originRspField.CanSet() {
		originRspField.Set(reflect.ValueOf(calResult.OriginRsp))
	}
	for i := 0; i < num; i++ {
		if strings.ToLower(rtype.Field(i).Tag.Get("cal")) == "origin_rsp" {
			originRspField = rval.Field(i)
			if originRspField.CanSet() {
				originRspField.Set(reflect.ValueOf(calResult.OriginRsp))
			}
			break
		}
	}

	return nil
}
