package model

import "github.com/jinzhu/gorm"

// 系统用户表
type SysUser struct {
	gorm.Model
	Name string `gorm:"size:255"` // 姓名
	Username string `gorm:"size:255"` // 用户名
	Password string `gorm:"size:255"` // 密码,应该是md5加密后存储
	OrgId uint `gorm:"size:20"` // 组织id
	Mobile string `gorm:"size:11"` // 手机号
	Enabled bool `gorm:"size:2;default:0"` // 禁用状态1,可以,0,禁用
	Photo string `gorm:"size:255"` // 用户头像地址
}

// 设置表名
func (SysUser) TableName() string {
	return "sys_user"
}