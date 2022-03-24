package env

import "github.com/layasugar/laya/gcf"

const (
	defaultAlarmType = ""
	defaultAlarmKey  = ""
	defaultAlarmHost = ""
)

func AlarmType() string {
	if gcf.IsSet("app.alarm.type") {
		return gcf.GetString("app.alarm.type")
	}
	return defaultAlarmType
}

func AlarmKey() string {
	if gcf.IsSet("app.alarm.key") {
		return gcf.GetString("app.alarm.key")
	}
	return defaultAlarmKey
}

func AlarmHost() string {
	if gcf.IsSet("app.alarm.addr") {
		return gcf.GetString("app.alarm.addr")
	}
	return defaultAlarmHost
}
