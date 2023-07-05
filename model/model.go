package model

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const defaultConfigFile = "./config/app.toml"

func Init(cli *cli.Context) error {
	tableName := cli.String("table")
	if tableName == "" {
		return errors.New("请传入表名")
	}
	packageName := cli.String("package")
	databaseName := cli.String("database")
	configFile := cli.String("config")
	out := cli.String("out")

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// 加载配置
	dbName, addr := readConfig(configFile, databaseName)
	if packageName == "" {
		packageName = dbName
	}

	// 计算输出目录
	var outPath string
	if out != "" {
		outPath = pwd + "/" + out + "/" + packageName
	} else {
		outPath = pwd + "/models/dao/db/" + packageName
	}

	return generate(packageName, dbName, tableName, addr, outPath)
}

func generate(packageName, database, table, addr, out string) error {
	fmt.Println("Connecting to mysql server " + addr)
	// 创建一个目录
	dPath := out
	err := os.MkdirAll(dPath, os.ModeDir)
	if err != nil {
		return errors.New(fmt.Sprintf("创建目录失败：%s" + dPath))
	}

	columnDataTypes, columnsSorted, err := GetColumnsFromMysqlTable(addr, database, table)
	if err != nil {
		return errors.New("error in selecting column data information from mysql information schema")
	}
	structName := ToCamelCase(table)

	// Generate struct string based on columnDataTypes
	structs, err := Generate(*columnDataTypes, columnsSorted, table, structName, packageName, true, true, false)
	if err != nil {
		return errors.New("Error in creating struct from json: " + err.Error())
	}

	fileName := dPath + "/" + table + ".go"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("Open File fail: " + err.Error())
	}
	_, err = file.WriteString(string(structs))
	if err != nil {
		return errors.New("Save File fail: " + err.Error())
	}
	fmt.Printf("success database: %s,table: %s\r\n", database, table)
	return nil
}

func readConfig(configFile string, databaseName string) (string, string) {
	// 加载配置
	if configFile == "" {
		viper.SetConfigFile(defaultConfigFile)
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	} else {
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	sDb := viper.Get("mysql")
	switch sDb.(type) {
	case []interface{}:
		si := sDb.([]interface{})
		for _, item := range si {
			if sim, ok := item.(map[string]interface{}); ok {
				var nameKey string
				var addr string
				if name, ok1 := sim["name"]; ok1 {
					if nameStr, okInterface := name.(string); okInterface {
						nameKey = nameStr
					}
				}
				if dsn, ok2 := sim["dsn"]; ok2 {
					if dsnStr, okInterface := dsn.(string); okInterface {
						addr = dsnStr
					}
				}

				if databaseName == "" && nameKey == "default" {
					return parseAddr(addr), addr
				}
				if databaseName != "" && nameKey == databaseName {
					return parseAddr(addr), addr
				}
			}
		}
	default:
		panic(fmt.Errorf("Fatal error config file: \n"))
	}
	return "", ""
}

// root:123456@tcp(127.0.0.1:3306)/laya-template?charset=utf8&parseTime=True&loc=Local
func parseAddr(addr string) string {
	return strings.Split(strings.Split(addr, "/")[1], "?")[0]
}
