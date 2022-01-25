package gcal

import (
	"github.com/layasugar/laya/gcal/converter"
	"github.com/layasugar/laya/gcal/protocol"
	"github.com/layasugar/laya/gcal/service"
)

// PbRPCRequest 别名
type PbRPCRequest = protocol.PbRPCRequest

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
