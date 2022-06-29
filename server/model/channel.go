package model

import "gorm.io/gorm"

// 通道表
type Channel struct {
	gorm.Model
	Name string `gorm:"size:255;column:name" json:"name"`                      // 通道名称
	Code string `gorm:"size:50;column:code" json:"code"` // 通道编码,一个流对应一个通道,推流标识,新建时就分配好
	StreamAddress string `gorm:"size:255;column:stream_address" json:"streamAddress"` // 取流地址(一般为内网地址)
	IsNvrAddress int `gorm:"size:2;column:is_nvr_address" json:"isNvrAddress"` // 是否为nvr地址,0,否,1,是
	StreamType int `gorm:"size:2;column:stream_type" json:"streamType"`        // 流的类型,rtsp或其他类型
	Status int `gorm:"size:2;column:status" json:"status"`                     // 通道状态,0,未开启,1,使用中
	PusherId int `gorm:"size:50;column:pusher_id" json:"pusherId"`              // 对应的客户端id
}

// 设置表名
func (Channel) TableName() string {
	return "channel"
}