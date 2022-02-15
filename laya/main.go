package main

import (
	"github.com/layasugar/laya/laya/model"
	"github.com/layasugar/laya/laya/template"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var commands = []*cli.Command{
	{
		Name:  "model",
		Usage: "生成模型代码",
		Subcommands: []*cli.Command{
			{
				Name:  "model",
				Usage: `generate model model`,
				Subcommands: []*cli.Command{
					{
						Name:  "ddl",
						Usage: `generate model model from ddl`,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "src, s",
								Usage: "the path or path globbing patterns of the ddl",
							},
						},
						Action: model.DDL,
					},
					{
						Name:  "datasource",
						Usage: `generate model from datasource`,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "url",
								Usage: `the data source of database,like "root:password@tcp(127.0.0.1:3306)/database"`,
							},
						},
						Action: model.DataSource,
					},
				},
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
