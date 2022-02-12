package tpl

import (
	"github.com/layasugar/laya/laya/tpl/common"
	"github.com/layasugar/laya/laya/tpl/http_tpl"
)

var PH = common.P{
	Name: "/",
	Files: []common.F{
		{Name: ".gitignore", Content: common.GitignoreTpl},
		{Name: "go.mod", Content: http_tpl.GoModTpl},
		{Name: "main.go", Content: http_tpl.MainTpl},
		{Name: "README.md", Content: common.ReadmeTpl},
	},
	Child: []common.P{
		{
			Name: "config",
			Files: []common.F{
				{Name: "app.toml", Content: http_tpl.ConfigAppTomlTpl},
			},
		},
		{
			Name: "controllers",
			Child: []common.P{
				{
					Name: "test",
					Files: []common.F{
						{Name: "base.go", Content: http_tpl.ControllersTestBaseTpl},
						{Name: "test.go", Content: http_tpl.ControllersTestTestTpl},
					},
				},
			},
			Files: []common.F{
				{Name: "base.go", Content: http_tpl.ControllersBaseTpl},
			},
		},
		{
			Name: "global",
			Child: []common.P{
				{
					Name: "errno",
					Files: []common.F{
						{Name: "err_code.go", Content: http_tpl.GlobalErrnoErrCodeTpl},
					},
				},
				{
					Name: "page",
					Files: []common.F{
						{Name: "page.go", Content: http_tpl.GlobalPagePageTpl},
					},
				},
			},
			Files: []common.F{
				{Name: "http_response.go", Content: http_tpl.GlobalHttpResponseTpl},
			},
		},
		{
			Name: "middlewares",
			Files: []common.F{
				{Name: "http_test_middleware.go", Content: http_tpl.MiddlewaresHttpTestTpl},
			},
		},
		{
			Name: "models",
			Child: []common.P{
				{
					Name: "dao",
					Files: []common.F{
						{Name: "base.go", Content: http_tpl.ModelsDaoBaseTpl},
					},
					Child: []common.P{
						{
							Name: "cal",
							Child: []common.P{
								{
									Name: "http_test",
									Files: []common.F{
										{Name: "trace.go", Content: http_tpl.ModelsDaoCalHttpTestTraceTpl},
									},
								},
							},
							Files: []common.F{
								{Name: "service_test.go", Content: http_tpl.ModelsDaoCalHttpTestServiceTestTpl},
							},
						},
						{
							Name: "db",
							Files: []common.F{
								{Name: "user.go", Content: http_tpl.ModelsDaoDbUserTpl},
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
								{Name: "trace.go", Content: http_tpl.ModelsDataTestTestTraceTpl},
								{Name: "types.go", Content: http_tpl.ModelsDataTestTypesTpl},
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
								{Name: "trace.go", Content: http_tpl.ModelsPageTestTraceTpl},
							},
						},
					},
				},
			},
		},
		{
			Name: "routes",
			Files: []common.F{
				{Name: "test.go", Content: http_tpl.RoutesTestTpl},
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
