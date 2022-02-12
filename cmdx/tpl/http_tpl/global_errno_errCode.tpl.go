package http_tpl

const GlobalErrnoErrCodeTpl = `package errno

import "{{.goModName}}/global"

var (
	ComUnauthorized = global.Err(401)
	UserNotFound    = global.Err(4001)
)
`
