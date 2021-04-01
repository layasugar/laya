package gconf

import (
	"github.com/layatips/laya/gconf/fileutil"
	"sync"
)

var watcher *fileutil.Watcher

var confCache sync.Map

var initLock sync.Mutex

func initWatcher() {
	initLock.Lock()
	defer initLock.Unlock()

	if watcher == nil {
		watcher = fileutil.NewWatcher(path, 65535)
		//默认需要监听所有的event
		_ = RegisterFileWatcher("/*", defaultConfChangeHandler)
	}
}

// 注册一个配置文件变化的事件handler
func RegisterFileWatcher(pattern string, handler fileutil.WatcherEventHandler) error {
	if watcher == nil {
		initWatcher()
	}
	return watcher.RegisterFileWatcher(pattern, handler)
}

// 配置文件变化默认回调
func defaultConfChangeHandler(e *fileutil.WatcherEvent) error {
	return nil
}
