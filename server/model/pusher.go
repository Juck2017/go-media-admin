package model

import "github.com/jinzhu/gorm"

// 推流客户端表
type Pusher struct {
	gorm.Model
	Name string `gorm:"size:50;column:name" json:"name"` //客户端名称
	ClientId string `gorm:"size:255;column:client_id" json:"clientId"` // 客户端id,唯一
	Status int `gorm:"size:2;column:status" json:"status"` // 客户端状态,0离线,1在线
}

// 设置表名
func (Pusher) TableName() string {
	return "pusher"
}