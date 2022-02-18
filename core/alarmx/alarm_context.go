package alarmx

import "github.com/layasugar/laya/genv"

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
	data["project"] = genv.AppName()
	data["env"] = genv.AppMode()
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
		if genv.AlarmType() == "" {
			return nil
		}
		alarm = &AlarmContext{
			alarmType: genv.AlarmType(),
			alarmHost: genv.AlarmHost(),
			alarmKey:  genv.AlarmKey(),
		}
	}
	return alarm
}
