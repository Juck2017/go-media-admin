package model

import (
	"gorm.io/gorm"
)

// 系统字典表
type SysDict struct {
	gorm.Model
	Code string `gorm:"size:255;column:code" json:"code"` // 字典编码
	Name string `gorm:"size:50;column:name" json:"name"` // 字典名称
	Leaf bool `gorm:"size:2;column:leaf" json:"leaf"`
	ParentId *int32 `gorm:"size:column:parent_id" json:"parentId,omitempty"` // 父id
	TreeId string `gorm:"size:255;column:tree_id" json:"treeId"` // 树id
	TreeLevel int `gorm:"size:5;column:tree_level;default:1;" json:"treeLevel"` // 树级别
}

// 设置表名
func (SysDict) TableName() string {
	return "sys_dict"
}