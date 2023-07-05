package common

const UtilsTypesTpl = `package utils

type HeaderParam struct {
	AppName    string {{.tagName}}header:"app-name"{{.tagName}}
	RequestID  string {{.tagName}}header:"request-id"{{.tagName}}
	XTEmployee string {{.tagName}}header:"xt-employee"{{.tagName}}
}

type XTEmployee struct {
	RealName string {{.tagName}}json:"real_name"{{.tagName}}
}

type AuthList struct {
	Name       string
	Sign       string
	HttpMethod string
	HttpPath   string
}
`
