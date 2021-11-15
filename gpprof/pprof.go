package gpprof

import (
	_ "github.com/mkevac/debugcharts"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
)

var p = &http.Server{Addr: "0.0.0.0:30085", Handler: nil}
var lock sync.Mutex

// StartPprof 开启pprof监控分析
func StartPprof() {
	lock.Lock()
	defer lock.Unlock()
	p = &http.Server{Addr: "0.0.0.0:30085", Handler: nil}
	go func() {
		log.Printf("[gpprof] http listen: 0.0.0.0:30085, %s", p.ListenAndServe())
	}()
	log.Printf("[gpprof] http listen: 0.0.0.0:30085, %s", "http: Server started")
}

// StopPprof 关闭pprof监控分析
func StopPprof() {
	lock.Lock()
	defer lock.Unlock()
	if p != nil {
		_ = p.Close()
	}
}
