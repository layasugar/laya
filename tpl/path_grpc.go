package tpl

import (
	"github.com/layasugar/laya/tpl/common"
	"github.com/layasugar/laya/tpl/grpc_tpl"
)

var PG = common.P{
	Name: "/",
	Files: []common.F{
		{Name: ".gitignore", Content: common.GitignoreTpl},
		{Name: "go.mod", Content: grpc_tpl.GoModTpl},
		{Name: "main.go", Content: grpc_tpl.MainTpl},
		{Name: "README.md", Content: common.ReadmeTpl},
	},
	Child: []common.P{
		{
			Name: "config",
			Files: []common.F{
				{Name: "app.toml", Content: grpc_tpl.ConfigAppTomlTpl},
			},
		},
		{
			Name: "controllers",
			Child: []common.P{
				{
					Name: "test",
					Files: []common.F{
						{Name: "base.go", Content: grpc_tpl.ControllersTestBaseTpl},
						{Name: "test.go", Content: grpc_tpl.ControllersTestTestTpl},
					},
				},
			},
		},
		{
			Name: "global",
			Child: []common.P{
				{
					Name: "page",
					Files: []common.F{
						{Name: "page.go", Content: grpc_tpl.GlobalPagePageTpl},
					},
				},
			},
		},
		{
			Name: "middlewares",
			Files: []common.F{
				{Name: "grpc_test_interceptor.go", Content: grpc_tpl.MiddlewaresGrpcTestTpl},
			},
		},
		{
			Name: "models",
			Child: []common.P{
				{
					Name: "dao",
					Files: []common.F{
						{Name: "base.go", Content: grpc_tpl.ModelsDaoBaseTpl},
					},
					Child: []common.P{
						{
							Name: "cal",
							Child: []common.P{
								{
									Name: "rpc_test",
									Files: []common.F{
										{Name: "trace.go", Content: grpc_tpl.ModelsDaoCalRpcTestTraceTpl},
									},
								},
							},
							Files: []common.F{
								{Name: "service_test.go", Content: grpc_tpl.ModelsDaoCalRpcTestServiceTestTpl},
							},
						},
						{
							Name: "db",
							Files: []common.F{
								{Name: "user.go", Content: grpc_tpl.ModelsDaoDbUserTpl},
							},
						},
					},
				},
				{
					Name: "data",
					Child: []common.P{
						{
							Name: "test",
							Files: []common.F{
								{Name: "trace.go", Content: grpc_tpl.ModelsDataTestTestTraceTpl},
								{Name: "types.go", Content: grpc_tpl.ModelsDataTestTypesTpl},
							},
						},
					},
				},
				{
					Name: "page",
					Child: []common.P{
						{
							Name: "test",
							Files: []common.F{
								{Name: "trace.go", Content: grpc_tpl.ModelsPageTestTraceTpl},
							},
						},
					},
				},
			},
		},
		{
			Name: "pb",
			Files: []common.F{
				{Name: "trace.proto", Content: grpc_tpl.PbTraceTpl},
				{Name: "trace.pb.go", Content: grpc_tpl.PbTracePbTpl},
			},
		},
		{
			Name: "routes",
			Files: []common.F{
				{Name: "test.go", Content: grpc_tpl.RoutesTestTpl},
			},
		},
		{
			Name: "utils",
			Files: []common.F{
				{Name: "functions.go", Content: common.UtilsFunctionsTpl},
				{Name: "types.go", Content: common.UtilsTypesTpl},
			},
		},
	},
}
