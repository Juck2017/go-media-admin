package service

import (
	"suitbim.com/go-media-admin/common"
	"suitbim.com/go-media-admin/core"
)

type ChannelService struct {
	BaseService
	Dispatch *core.Dispatcher
}

type ChannelRequestDto struct {
	PageSize int `json:"pageSize" binding:"required"`
	CurrPage int `json:"currPage" binding:"required"`
	AreaId   int `json:"areaId"`
}

type ChannelDto struct {
	Id uint `json:"id"`
	PusherName string `json:"pusherName,omitempty"` // 所属客户端名称
	Code string `json:"code"`
	Name string `json:"name"`
	AccessTime *FormatTime `json:"accessTime,omitempty"`
	StreamAddress string `json:"streamAddress"`
	IsNvrAddress string `json:"isNvrAddress,omitempty"`
	Status *int `json:"status,omitempty"`     // 当前通道状态
	OnlineNum *int `json:"onlineNum,omitempty"` // 当前通道在线使用数量
}

/**
	获取通道列表
	获取当前通道在线数量,当前通道状态需要单独获取
 */
func (c *ChannelService) GetChannelList(cameraDto ChannelRequestDto) ([]ChannelDto, int64) {
	var channelList []ChannelDto
	var countNum int64
	dbData := common.GlobalDB.Table("channel").
		Select("channel.id,pusher.name AS pusher_name,channel.code,channel.name,channel.stream_address," +
			"channel.created_at AS access_time," +
			"(CASE channel.is_nvr_address WHEN 1 THEN '是' ELSE '否' END) AS is_nvr_address").
		Joins("LEFT JOIN pusher ON pusher.id=channel.pusher_id")
	//if cameraDto.AreaId != 0 {
	//	dbData = dbData.Where("area_id = ?", cameraDto.AreaId)
	//}
	dbData.Count(&countNum)
	dbData.Offset((cameraDto.CurrPage-1)*cameraDto.PageSize).Limit(cameraDto.PageSize).Scan(&channelList)
	for i := 0; i<len(channelList);i++ {
		channelList[i].Status = c.getChannelStatus(channelList[i].Code)
		channelList[i].OnlineNum = c.getChannelOnlineNum(channelList[i].Code)
	}
	return channelList, countNum
}

// 根据客户端id获取通道列表
func (c *ChannelService) GetChannelListByClientId(clientId string) []ChannelDto {
	var channelList []ChannelDto
	common.GlobalDB.Table("channel").
		Select("channel.id,channel.code,channel.name,channel.stream_address").
		Joins("LEFT JOIN pusher ON pusher.id=channel.pusher_id").
		Where("pusher.client_id=?", clientId).Scan(&channelList)
	return channelList
}

// 获取对应通道状态
func (c *ChannelService) getChannelStatus(channelCode string) *int {
	if len(c.Dispatch.ChannelMap[channelCode].Viewers) > 0 {
		online := 1
		return &online
	}else {
		offline := 0
		return &offline
	}
}

// 获取对应通道在线数量
func (c *ChannelService) getChannelOnlineNum(channelCode string) *int {
	num := len(c.Dispatch.ChannelMap[channelCode].Viewers)
	return &num
}
