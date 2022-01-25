package cal

import (
	"gitlab.xthktech.cn/bs/gxe/env"
	"gitlab.xthktech.cn/bs/gxe/cal/converter"
	"gitlab.xthktech.cn/bs/gxe/cal/log"
	"gitlab.xthktech.cn/bs/gxe/cal/protocol"
	"gitlab.xthktech.cn/bs/gxe/cal/service"
	"gitlab.xthktech.cn/bs/gxe/utils/fileutil"
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

// var MCPACK1Converter = converter.MCPACK1

// RAWConverter 别名
var RAWConverter = converter.RAW

// LoadService load one service from struct
// Deprecated
func LoadService(serc *service.Config, idc string) error {
	return service.LoadService(serc, idc)
}

// LoadServices load one service from struct
// Deprecated
func LoadServices(sercs []*service.Config, idc string) error {
	for _, serc := range sercs {
		if err := service.LoadService(serc, idc); err != nil {
			return err
		}
	}

	return nil
}

// LoadServiceFromFile load one service from config file
// Deprecated
func LoadServiceFromFile(file string, idc string) error {
	return service.LoadServiceFromTOMLFile(file, idc)
}

// LoadServicesFromFolder load multiple services from config file folder
// Deprecated
func LoadServicesFromFolder(folder string, idc string) error {
	return service.LoadServicesFromTOMLDir(folder, idc)
}

// ConfigFileWatcher 监听配置问题
func ConfigFileWatcher(event *fileutil.WatcherEvent) error {
	switch event.Type {
	case fileutil.WatcherEventCreate, fileutil.WatcherEventChange:
		return service.LoadServiceFromTOMLFile(event.Path, env.IDC())
	case fileutil.WatcherEventDelete:
		return service.RemoveServiceByTOMLFile(event.Path)
	}

	return nil
}

func SetLogPath(path string) bool {
	return log.SetLogPath(path)
}
