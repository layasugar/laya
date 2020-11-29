package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

// UserParserFn 用户自定义的解析配置文件的函数
type UserParserFn func(confStr string) (string, error)

var userParsers []UserParserFn

func init() {
	RegisterUserParser(envVarsUserParser)
}

// ReadFile 读取toml格式的配置文件
func ReadFile(confPath string, config interface{}) (err error) {

	confPath = cleanPath(confPath)

	cacheKey := fmt.Sprintf("%s%s", cacheKeyPrefix(confPath), reflect.TypeOf(config).String())

	vs, ok := confCache.Load(cacheKey)

	if !ok {
		var bs []byte
		bs, err = ioutil.ReadFile(FileAbsPath(confPath))
		if err != nil {
			return
		}
		var confStrRaw = string(bs)
		var confStr string
		confStr, err = parserConfWithUserParserFns(confStrRaw)

		//若配置文件有变化，则异步写入data/var/config/dump 目录下去，以方便查看
		if confStrRaw != confStr {
			go dumpConf(confPath, confStr)
		}
		_, err = toml.Decode(confStr, config)
		confCache.Store(cacheKey, reflect.ValueOf(config).Elem().Interface())
		return
	}

	vp := reflect.ValueOf(config)
	vp.Elem().Set(reflect.ValueOf(vs))

	return
}

func cacheKeyPrefix(confPath string) string {
	return fmt.Sprintf("%s?!", confPath)
}

// ListFiles 返回当前配置目录的toml配置文件名 []string{"aa.toml","bb.toml"}
func ListFiles(confDir string) []string {
	pattern := filepath.Join(RootPath(), confDir, "*.toml")

	fl, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}

	var files []string
	for _, fPath := range fl {
		name, _ := filepath.Rel(RootPath(), fPath)
		if strings.HasPrefix(filepath.Base(name), `.`) {
			continue
		}
		files = append(files, name)
	}
	return files
}

// RegisterUserParser 注册一个自定义配置文件解析方法
func RegisterUserParser(fn UserParserFn) {
	if userParsers == nil {
		userParsers = make([]UserParserFn, 0, 1)
	}
	userParsers = append(userParsers, fn)
}

// parserConfWithUserParserFns 使用用户自定义还是解析配置文件的内容
func parserConfWithUserParserFns(confStr string) (str string, err error) {
	str = confStr
	for _, fn := range userParsers {
		str, err = fn(str)
		if err != nil {
			return "", err
		}
	}
	return
}

//模板变量格式：{env.变量名} 或者 {env.变量名|默认值}
var envVarReg = regexp.MustCompile(`\{env\.([A-Za-z0-9_]+)(\|[^}]+)?\}`)

// defaultUserParserEnvVars 解析配置文件中的环境变量
// 将配置文件中的 {env.xxx} 的内容，从环境变量中读取并替换
func envVarsUserParser(confStr string) (string, error) {
	newStr := envVarReg.ReplaceAllStringFunc(confStr, func(subStr string) string {
		// 将 {env.xxx} 中的 xxx 部分取出
		// 或者 将 {env.xxx|val} 中的 xxx|val 部分取出
		keyWithDefaultVal := subStr[5 : len(subStr)-1]
		idx := strings.Index(keyWithDefaultVal, "|")
		if idx > 0 {
			//{env.变量名|默认值} 有默认值的格式
			key := keyWithDefaultVal[:idx]
			defaultVal := keyWithDefaultVal[idx+1:]
			envVal := os.Getenv(key)
			if envVal == "" {
				return defaultVal
			}
			return envVal
		}

		//{env.变量名} 无默认值的部分
		return os.Getenv(keyWithDefaultVal)
	})
	return newStr, nil
}
