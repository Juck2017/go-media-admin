package model

import "gorm.io/gorm"

// 系统组织权限表
type SysOrgPrivilege struct {
	gorm.Model
	OrgId int `gorm:"size:255;column:org_id" json:"orgId"` // 组织id
	PrivilegeId int `gorm:"size:50;column:privilege_id" json:"privilegeId"` // 权限id
}

// 设置表名
func (SysOrgPrivilege) TableName() string {
	return "sys_org_privilege"
}