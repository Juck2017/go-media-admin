package model

import "gorm.io/gorm"

// 系统角色权限表
type SysRolePrivilege struct {
	gorm.Model
	PrivilegeId int `gorm:"size:255;column:privilege_id" json:"privilegeId"` // 权限id
	RoleId int `gorm:"size:255;column:role_id" json:"roleId"` // 角色id
}

// 设置表名
func (SysRolePrivilege) TableName() string {
	return "sys_role_privilege"
}