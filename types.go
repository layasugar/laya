package laya

import (
	"github.com/LaYa-op/laya/i18n"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
)

// 定义redis链接池,mysql连接池,语言包bundle
var Redis *redis.Client
var DB *gorm.DB
var I18n *i18n.I18n
