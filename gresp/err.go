package gresp


var Errno = map[uint32]string{
	400:    "系统错误",
}

var (
	SystemErr       = Err(400)
)

// RspError
type RspError struct {
	Code uint32
	Msg  string
}

func (re *RspError) Error() string {
	return re.Msg
}

func Err(code uint32) (err error) {
	err = &RspError{
		Code: code,
		Msg:  Errno[code],
	}
	return err
}

// Render
func (re *RspError) Render() (code uint32, msg string) {
	return re.Code, re.Msg
}
