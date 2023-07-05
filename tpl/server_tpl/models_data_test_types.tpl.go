package server_tpl

const ModelsDataTestTypesTpl = `package test

type Rsp struct {
	Code string {{.tagName}}json:"code"{{.tagName}}
}
`
