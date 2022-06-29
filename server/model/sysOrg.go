package model

import "gorm.io/gorm"

// 系统组织表
type SysOrg struct {
	gorm.Model
	Code string `gorm:"size:255;column:code" json:"code"` // 组织编码
	Name string `gorm:"size:50;column:name" json:"name"` // 组织名称
	Leaf bool `gorm:"size:2;column:leaf" json:"leaf"` // leaf: 0或者1
	ParentId *int32 `gorm:"size:20;column:parent_id" json:"parentId"` // 父id
	TreeId string `gorm:"size:255;column:tree_id" json:"treeId"` // 树id
	TreeLevel int `gorm:"size:5;column:tree_level" json:"treeLevel"` // 树级别
	Type string `gorm:"size:255;column:type" json:"type"` // 组织类型,company代表公司,dept代表部门
}

// 设置表名
func (SysOrg) TableName() string {
	return "sys_org"
}