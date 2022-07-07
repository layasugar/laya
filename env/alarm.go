package env

import "github.com/layasugar/laya/gcnf"

const (
	defaultAlarmType = ""
	defaultAlarmKey  = ""
	defaultAlarmHost = ""

	_appAlarmType = "app.alarm.type"
	_appAlarmKey  = "app.alarm.key"
	_appAlarmAddr = "app.alarm.addr"
)

func AlarmType() string {
	if gcnf.IsSet(_appAlarmType) {
		return gcnf.GetString(_appAlarmType)
	}
	return defaultAlarmType
}

func AlarmKey() string {
	if gcnf.IsSet(_appAlarmKey) {
		return gcnf.GetString(_appAlarmKey)
	}
	return defaultAlarmKey
}

func AlarmHost() string {
	if gcnf.IsSet(_appAlarmAddr) {
		return gcnf.GetString(_appAlarmAddr)
	}
	return defaultAlarmHost
}
