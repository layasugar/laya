package template

import (
	"fmt"
	"github.com/layasugar/laya/tpl"
	"github.com/layasugar/laya/tpl/common"
	"github.com/layasugar/laya/version"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var projectName = ""
var goModName = ""
var tagName = "`"
var versionName = version.VERSION

func GenHttpTemplates(ctx *cli.Context) error {
	name := ctx.String("name")
	projectName, goModName = parseName(name)
	log.Printf("start, project_name: %s, go_mod_name: %s", projectName, goModName)

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	absPwd, err := filepath.Abs(pwd + "/" + projectName)
	if err != nil {
		return err
	}

	err = os.Mkdir(absPwd, 0664)
	if err != nil {
		return err
	}

	recursion(absPwd, tpl.PH)

	return nil
}

func GenGrpcTemplates(ctx *cli.Context) error {
	name := ctx.String("name")
	projectName, goModName = parseName(name)
	log.Printf("start, project_name: %s, go_mod_name: %s", projectName, goModName)

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	absPwd, err := filepath.Abs(pwd + "/" + projectName)
	if err != nil {
		return err
	}

	err = os.Mkdir(absPwd, 0664)
	if err != nil {
		return err
	}

	recursion(absPwd, tpl.PG)

	return nil
}

func GenServerTemplates(ctx *cli.Context) error {
	name := ctx.String("name")
	projectName, goModName = parseName(name)
	log.Printf("start, project_name: %s, go_mod_name: %s", projectName, goModName)

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	absPwd, err := filepath.Abs(pwd + "/" + projectName)
	if err != nil {
		return err
	}

	err = os.Mkdir(absPwd, 0664)
	if err != nil {
		return err
	}

	recursion(absPwd, tpl.PS)

	return nil
}

func recursion(cp string, p common.P) {
	var currentPath string
	if p.Name == "/" {
		currentPath = cp
	} else {
		currentPath = cp + "/" + p.Name
		err := os.Mkdir(currentPath, 0664)
		if err != nil {
			log.Print(err.Error())
		}
	}

	for _, f := range p.Files {
		var (
			err  error
			file *os.File
		)
		tt := template.Must(template.New("queue").Parse(f.Content))
		fn := currentPath + "/" + f.Name
		if file, err = os.Create(fn); err != nil {
			if !os.IsExist(err) {
				fmt.Printf("Could not create %s: %s (skip)\n", fn, err)
				continue
			}
			_ = os.Remove(fn)
		}

		v := map[string]string{
			"projectName": projectName,
			"goModName":   goModName,
			"tagName":     tagName,
			"versionName": versionName,
		}
		_ = tt.Execute(file, v)
		_ = file.Close()
	}

	for _, pt := range p.Child {
		recursion(currentPath, pt)
	}
}
