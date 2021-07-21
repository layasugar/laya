package gpprof

import (
	_ "github.com/mkevac/debugcharts"
	"log"
	"net/http"
	_ "net/http/pprof"
)

// 30085作为监控接口
func InitPprof() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:30085", nil))
	}()
}
