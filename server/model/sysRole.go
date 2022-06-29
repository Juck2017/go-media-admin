package model

import "gorm.io/gorm"

// 系统角色表
type SysRole struct {
	gorm.Model
	Code string `gorm:"size:255;column:code" json:"code"` // 角色编码
	Name string `gorm:"size:50;column:name" json:"name"` // 角色名称
	Description string `gorm:"size:255;column:description" json:"description"` // 角色描述
	OrgId int `gorm:"size:255;column:org_id" json:"orgId"` // 组织id
}

// 设置表名
func (SysRole) TableName() string {
	return "sys_role"
}