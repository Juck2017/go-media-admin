package model
import "github.com/jinzhu/gorm"

// 系统用户表
type SysUserPrivilege struct {
	gorm.Model
	Deny uint `gorm:"size:2"` // 是否存在:0是存在, 1是剔除
	PrivilegeId uint `gorm:"size:255"` // 权限id
	UserId uint `gorm:"size:255"` // 用户id
}

// 设置表名
func (SysUserPrivilege) TableName() string {
	return "sys_user_privilege"
}
