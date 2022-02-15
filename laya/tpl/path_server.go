package tpl

import (
	"github.com/layasugar/laya/laya/tpl/common"
	"github.com/layasugar/laya/laya/tpl/server_tpl"
)

var PS = common.P{
	Name: "/",
	Files: []common.F{
		{Name: ".gitignore", Content: common.GitignoreTpl},
		{Name: "go.mod", Content: server_tpl.GoModTpl},
		{Name: "main.go", Content: server_tpl.MainTpl},
		{Name: "README.md", Content: common.ReadmeTpl},
	},
	Child: []common.P{
		{
			Name: "config",
			Files: []common.F{
				{Name: "app.toml", Content: server_tpl.ConfigAppTomlTpl},
			},
		},
		{
			Name: "controllers",
			Child: []common.P{
				{
					Name: "test",
					Files: []common.F{
						{Name: "base.go", Content: server_tpl.ControllersTestBaseTpl},
						{Name: "test.go", Content: server_tpl.ControllersTestTestTpl},
					},
				},
			},
		},
		{
			Name: "models",
			Child: []common.P{
				{
					Name: "dao",
					Files: []common.F{
						{Name: "base.go", Content: server_tpl.ModelsDaoBaseTpl},
					},
					Child: []common.P{
						{
							Name: "cal",
							Child: []common.P{
								{
									Name: "task_test",
									Files: []common.F{
										{Name: "trace.go", Content: server_tpl.ModelsDaoCalHttpTestTraceTpl},
									},
								},
							},
							Files: []common.F{
								{Name: "service_test.go", Content: server_tpl.ModelsDaoCalHttpTestServiceTestTpl},
							},
						},
						{
							Name: "db",
							Files: []common.F{
								{Name: "user.go", Content: server_tpl.ModelsDaoDbUserTpl},
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
								{Name: "task_trace.go", Content: server_tpl.ModelsDataTestTestTraceTpl},
								{Name: "types.go", Content: server_tpl.ModelsDataTestTypesTpl},
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
								{Name: "trace.go", Content: server_tpl.ModelsPageTestTraceTpl},
							},
						},
					},
				},
			},
		},
		{
			Name: "routes",
			Files: []common.F{
				{Name: "test.go", Content: server_tpl.RoutesTestTpl},
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
