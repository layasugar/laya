package gutils

import (
	"errors"
	"github.com/google/uuid"
	"net"
	"os"
)

// 创建LogID
func NewLogID() string {
	return uuid.New().String()
}

// 获取本地IP
func LocalIP() (string, error) {
	adds, err := net.InterfaceAddrs()
	if err != nil {
		return os.Hostname()
	}

	for _, a := range adds {
		if inet, ok := a.(*net.IPNet); ok && !inet.IP.IsLoopback() {
			if inet.IP.To4() != nil {
				return inet.IP.String(), nil
			}
		}
	}

	return "", errors.New("fail to get local ip")
}
