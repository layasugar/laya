package response

const (
	SysErr       = "-1"
	InvalidInput = "400"
	Unauthorized = "401"
	ServerErr    = "500"
)

var (
	ErrSysErr       = register(SysErr)
	ErrInvalidInput = register(InvalidInput)
	ErrUnauthorized = register(Unauthorized)
	ErrServerErr    = register(ServerErr)
)
