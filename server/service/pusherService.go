package service

import (
	"suitbim.com/go-media-admin/common"
)

type PusherService struct {
	BaseService
}

type PusherRequestDto struct {
	PageSize int `json:"pageSize" binding:"required"`
	CurrPage int `json:"currPage" binding:"required"`
	Name string `json:"name"`
}

type PusherDto struct {
	Id uint `json:"id"`
	Name string `json:"name"`
	AccessTime *FormatTime `json:"accessTime"`
	ClientId string `json:"clientId"`
	Status string `json:"status"`     // 当前通道状态
}

/**
	获取客户端列表
*/
func (p *PusherService) GetPusherList(pusherDto PusherRequestDto) ([]PusherDto, int64, error) {
	var pusherList []PusherDto
	var countNum int64
	dbData := common.GlobalDB.Debug().Table("pusher").
		Select("id," +
			"created_at AS access_time," +
			"name," +
			"client_id as client_id," +
			"CASE status WHEN 1 THEN '在线' ELSE '离线' END as status")
	if pusherDto.Name != "" {
		dbData = dbData.Where("name like ?", "%"+pusherDto.Name+"%")
	}
	dbData.Count(&countNum)
	dbData.Offset((pusherDto.CurrPage-1)*pusherDto.PageSize).Limit(pusherDto.PageSize).Scan(&pusherList)
	return pusherList, countNum, nil
}
