package gcnf

import "time"

const (
	defaultNullString = ""
)

var configChargeHandleFunc []func()
var t *time.Timer
