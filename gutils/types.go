package gutils

import (
	jsonIter "github.com/json-iterator/go"
)

const (
	XForwardedFor = "X-Forwarded-For" // 获取真实ip
	XRealIP       = "X-Real-IP"       // 获取真实ip
)

var (
	//CJson 全局json对象
	CJson = jsonIter.ConfigCompatibleWithStandardLibrary
)
