package model

import "github.com/jinzhu/gorm"

// 系统用户角色表
type SysUserRole struct {
	gorm.Model
	RoleId int `gorm:"size:20" json:"role_id"` // 角色id
	UserId int `gorm:"size:11" json:"user_id"` // 用户id
}

// 设置表名
func (SysUserRole) TableName() string {
	return "sys_user_role"
}