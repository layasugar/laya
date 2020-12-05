package response

import "strconv"

// RspError 和响应关联的错误，在models中使用
type RspError struct {
	RespData
}

func (re *RspError) Error() string {
	return re.Msg
}

func register(errStr string) error {
	err := &RspError{}
	err.Code, _ = strconv.Atoi(errStr)
	return err
}

// Render 展示是返回的数据
func (re *RspError) Render() (no int, msg string) {
	return re.Code, re.Msg
}

// Exchange 用err置换展示的数据
func Exchange(err error) (code int, msg string) {
	switch err.(type) {
	case *RspError:
		return err.(*RspError).Render()
	default:
		return Exchange(ErrSysErr)
	}
}
