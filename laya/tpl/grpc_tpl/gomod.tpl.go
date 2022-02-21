package grpc_tpl

const GoModTpl = `module {{.goModName}}

go 1.17

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/layasugar/laya {{.versionName}}
	github.com/opentracing/opentracing-go v1.2.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
	gorm.io/gorm v1.22.5
)
`
