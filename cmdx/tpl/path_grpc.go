package tpl

import "github.com/layasugar/laya/cmdx/tpl/grpc_tpl"

var PG = P{
	Name: "/",
	Files: []F{
		{Name: ".gitignore", Content: gitignore},
		{Name: "go.mod", Content: grpc_tpl.GoModTpl},
		{Name: "main.go", Content: grpc_tpl.MainTpl},
		{Name: "README.md", Content: readmeTpl},
	},
	Child: []P{
		{
			Name: "config",
			Files: []F{
				{Name: "app.toml", Content: grpc_tpl.AppTomlTpl},
			},
		},
		{
			Name: "controllers",
			Files: []F{
				{Name: "base.go", Content: grpc_tpl.BaseTpl},
			},
			Child: []P{
				{
					Name: "test",
					Files: []F{
						{Name: "base.go", Content: grpc_tpl.BaseTpl},
						{Name: "test.go", Content: grpc_tpl.BaseTpl},
					},
				},
			},
		},
		{Name: "global", Child: []P{
			{Name: "errno"},
			{Name: "page"},
		}},
		{Name: "middlewares"},
		{Name: "models", Child: []P{
			{Name: "dao", Child: []P{
				{Name: "cal", Child: []P{
					{Name: "rpc_test"},
				}},
				{Name: "db"},
			}},
			{Name: "data", Child: []P{
				{Name: "test"},
			}},
			{Name: "page", Child: []P{
				{Name: "test"},
			}},
		}},
		{Name: "pb"},
		{Name: "routes"},
		{Name: "utils"},
	},
}