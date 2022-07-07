package alarm

import "github.com/layasugar/laya/env"

// AlarmContext alarm
type AlarmContext struct {
	alarmType string
	alarmHost string
	alarmKey  string
}

var alarm *AlarmContext

func Alarm(title string, content string, data map[string]interface{}) {
	al := getAlarm()
	if nil == al {
		return
	}
	data["project"] = env.AppName()
	data["env"] = env.AppMode()
	switch al.alarmType {
	case "dingding":
		var d = AlarmData{
			Title:       title,
			Description: content,
			Content:     data,
		}
		go DingAlarm(al.alarmKey, al.alarmHost, d)
	case "http":
	}
}

func getAlarm() *AlarmContext {
	if nil == alarm {
		if env.AlarmType() == "" {
			return nil
		}
		alarm = &AlarmContext{
			alarmType: env.AlarmType(),
			alarmHost: env.AlarmHost(),
			alarmKey:  env.AlarmKey(),
		}
	}
	return alarm
}

func ReloadAlarm() {
	if nil == alarm {
		if env.AlarmType() == "" {
			return
		}
		alarm = &AlarmContext{
			alarmType: env.AlarmType(),
			alarmHost: env.AlarmHost(),
			alarmKey:  env.AlarmKey(),
		}
	}
}
