package http_tpl

const GoModTpl = `module {{.goModName}}

go 1.17

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/layasugar/laya v0.4.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	gorm.io/gorm v1.22.5
)
`
