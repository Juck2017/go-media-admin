package model

import "gorm.io/gorm"

// 客户端-通道映射表
type PusherChannel struct {
	gorm.Model
	ChannelId uint `gorm:"size:50;column:channel_id" json:"channelId"` //通道Id,多个
	PusherId uint `gorm:"size:50;column:pusher_id" json:"pusherId"` // pusher的id,可重复
}

// 设置表名
func (PusherChannel) TableName() string {
	return "pusher_channel"
}