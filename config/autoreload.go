package config

import (
	"fmt"
	"github.com/LaYa-op/laya/utils/fileutil"
	"log"
	"strings"
	"sync"
)

var watcher *fileutil.Watcher

var confCache sync.Map

var initLock sync.Mutex

func initWatcher() {
	initLock.Lock()
	defer initLock.Unlock()

	if watcher == nil {
		watcher = fileutil.NewWatcher(RootPath(), 65535)
		//默认需要监听所有的event
		_ = RegisterFileWatcher("/*", defaultConfChangeHandler)
	}
}

// RegisterFileWatcher 注册一个配置文件变化的事件handler
func RegisterFileWatcher(pattern string, handler fileutil.WatcherEventHandler) error {
	if watcher == nil {
		initWatcher()
	}
	return watcher.RegisterFileWatcher(pattern, handler)
}

// defaultConfChangeHandler 配置文件变化默认回调
func defaultConfChangeHandler(e *fileutil.WatcherEvent) error {
	log.Println("[conf_autoreload] config file changed:", e)
	cacheKeyPre := cacheKeyPrefix(cleanPath(e.Name))
	confCache.Range(func(key interface{}, value interface{}) bool {
		keyStr := fmt.Sprintf("%v", key)
		if strings.HasPrefix(keyStr, cacheKeyPre) {
			confCache.Delete(key)
			log.Println("[conf_autoreload] config file was changed and clean the cache:", e)
		}
		return true
	})
	return nil
}
