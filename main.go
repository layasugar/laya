package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/layasugar/laya/model"
	"github.com/layasugar/laya/template"
)

var commands = []*cli.Command{
	{
		Name:  "model",
		Usage: "生成模型代码",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: `generate model model`,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "table",
						Usage: `指定表名, 必传`,
					},
					&cli.StringFlag{
						Name:  "package",
						Usage: `指定包名, 默认库名`,
					},
					&cli.StringFlag{
						Name:  "config",
						Usage: `指定配置文件, 默认当前路径config/app.toml`,
					},
					&cli.StringFlag{
						Name:  "out",
						Usage: `指定输出路径, 默认是models/dao/db/package`,
					},
					&cli.StringFlag{
						Name:  "database",
						Usage: `指定数据库名, 默认是default配置的连接数据库`,
					},
				},
				Action: model.Init,
			},
		},
	},
	{
		Name:  "template",
		Usage: "生成框架模板",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "初始化一个默认程序模板",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "指定gomod名称",
					},
				},
				Action: template.GenServerTemplates,
			},
			{
				Name:  "init-http",
				Usage: "初始化一个http模板",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "指定gomod名称",
					},
				},
				Action: template.GenHttpTemplates,
			},
			{
				Name:  "init-grpc",
				Usage: "初始化一个grpc模板",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "指定gomod名称",
					},
				},
				Action: template.GenGrpcTemplates,
			},
		},
	},
}

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = commands

	// cli already print error messages
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
