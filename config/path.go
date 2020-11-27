package conf

import (
	"fmt"
	"github.com/LaYa-op/laya/env"
	"github.com/LaYa-op/laya/utils/fileutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// RootPath 返回配置文件的根目录
func RootPath() string {
	return env.ConfRootPath()
}

// FileAbsPath 返回配置文件的完整路径
// confRelPath 为相对路径 eg：db/cluster/db_a.toml
func FileAbsPath(confRelPath string) string {
	return filepath.Join(RootPath(), confRelPath)
}

// cleanPath 将文件路径归一
// 如 输入的是 /app.toml 归一化为 app.toml
// 如 输入的是 db/../app.toml 归一化为 app.toml
func cleanPath(confPath string) string {
	if strings.HasPrefix(confPath, "/") {
		confPath = strings.TrimLeft(confPath, "/")
	}
	return filepath.Clean(confPath)
}

// FileRelPath 返回配置文件的相对路径
// confAbsPath 为配置文件的绝对路径
func FileRelPath(confAbsPath string) (string, error) {
	return filepath.Rel(RootPath(), confAbsPath)
}

func validConfPath(confPath string) error {
	if strings.Contains(confPath, "../") {
		return fmt.Errorf("path='%s' contains '../' is not allow", confPath)
	}
	return nil
}

func dumpConfDir() string {
	return filepath.Join(env.DataRootPath(), "var", "conf", "dump")
}

var dumpOnce sync.Once

// dumpConf 导出配置文件内容到指定目录
func dumpConf(confPath string, confStr string) {
	dumpOnce.Do(func() {
		root := dumpConfDir()
		if fileutil.Exists(root) {
			os.RemoveAll(root)
		}
	})

	fpath := filepath.Join(dumpConfDir(), confPath)
	dir := filepath.Dir(fpath)
	if !fileutil.Exists(dir) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Println("create dump conf dir failed:", err)
			return
		}
	}
	err := fileutil.FilePutContents(fpath, []byte(confStr))
	if err != nil {
		log.Printf("wrote dump_file (%s) failed, err= %s\n", fpath, err)
	}
}
