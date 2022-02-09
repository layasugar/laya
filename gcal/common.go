package gcal

import (
	"github.com/layasugar/laya/gcal/converter"
	"github.com/layasugar/laya/gcal/protocol"
	"github.com/layasugar/laya/gcal/service"
	"github.com/layasugar/laya/gcal/pool"
)

// grpc连接池
var pbTc = &pool.Pool{}

// HTTPRequest 别名
type HTTPRequest = protocol.HTTPRequest

// HTTPHead 别名
type HTTPHead = protocol.HTTPHead

// ConverterType 别名
type ConverterType = converter.ConverterType

// JSONConverter 别名
var JSONConverter = converter.JSON

// FORMConverter 别名
var FORMConverter = converter.FORM

// RAWConverter 别名
var RAWConverter = converter.RAW

// LoadService load one service from struct
func LoadService(configs []map[string]interface{}) error {
	return service.LoadService(configs)
}

func GetRpcConn(serverName string) {

}