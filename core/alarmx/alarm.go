package alarmx

type AlarmData struct {
	Title       string                 //报警标题
	Description string                 //报警描述
	Content     map[string]interface{} //kv数据
}