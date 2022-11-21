package gcnf

const (
	defaultAlarmType = ""
	defaultAlarmKey  = ""
	defaultAlarmHost = ""

	_appAlarmType = "app.alarm.type"
	_appAlarmKey  = "app.alarm.key"
	_appAlarmAddr = "app.alarm.addr"
)

func AlarmType() string {
	if gcf.IsSet(_appAlarmType) {
		return gcf.GetString(_appAlarmType)
	}
	return defaultAlarmType
}

func AlarmKey() string {
	if gcf.IsSet(_appAlarmKey) {
		return gcf.GetString(_appAlarmKey)
	}
	return defaultAlarmKey
}

func AlarmHost() string {
	if gcf.IsSet(_appAlarmAddr) {
		return gcf.GetString(_appAlarmAddr)
	}
	return defaultAlarmHost
}
