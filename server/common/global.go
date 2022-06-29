package common

import (
	"gorm.io/gorm"
)

var (
	GlobalDB  *gorm.DB       // 全局DB
	GlobalConf *Config // 全局配置信息
)