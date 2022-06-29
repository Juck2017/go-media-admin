package model

import (
	"gorm.io/gorm"
)

// 系统权限表
type SysPrivilege struct {
	gorm.Model
	Code string `gorm:"size:255;column:code" json:"code"` // 权限编码
	Name string `gorm:"size:50;column:name" json:"name"` // 权限名称
	Leaf bool `gorm:"size:2;column:leaf" json:"leaf"`
	ParentId *int32 `gorm:"size:20;column:parent_id" json:"parentId,omitempty"` // 父id
	TreeId string `gorm:"size:255;column:tree_id" json:"treeId"` // 树id
	TreeLevel int `gorm:"size:5;column:tree_level" json:"treeLevel"` // 树级别
	Description string `gorm:"size:255;column:description" json:"description"` // 权限描述
	Scope string `gorm:"size:255;column:scope" json:"scope"` // system
	Type string `gorm:"size:255;column:type" json:"type"` // 类型:menu代表菜单,button代表按钮
}

// 设置表名
func (SysPrivilege) TableName() string {
	return "sys_privilege"
}