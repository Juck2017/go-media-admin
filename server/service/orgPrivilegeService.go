package service

import (
	"errors"
	"suitbim.com/go-media-admin/common"
	"suitbim.com/go-media-admin/model"
)

type OrgPrivilegeDto struct {
	Id uint `json:"id"`
	OrgId int `json:"orgId"` // 组织id
	PrivilegeId int `json:"PrivilegeId"` // 权限id
}

type ConfigPrivilegeDto struct {
	OrgId int `json:"orgId" binding:"required"`
	PrivilegeIds []int `json:"privilegeIds" binding:"required"`
}

type OrgPrivilegeService struct {
	BaseService
}

// 获取组织权限列表
func (o *OrgPrivilegeService) GetOrgPrivilegeList(orgId uint64)  ([]PrivilegeDto,error) {
	var opIds []uint
	var privileges []PrivilegeDto

	var orgList []model.SysOrg
	common.GlobalDB.Debug().Table("sys_org").Scan(&orgList)
	arrList := make([]interface{}, len(orgList))
	for i := 0; i <  len(orgList); i++ {
		arrList[i] = &orgList[i]
	}
	parentOrgId, err := o.GetTreeParentId(arrList, []uint64{}, orgId)
	if err != nil {
		return privileges, err
	}
	common.GlobalDB.Debug().Table("sys_org_privilege").Select("privilege_id").Where("org_id=?", parentOrgId).Scan(&opIds)
	common.GlobalDB.Debug().Table("sys_privilege").Where("id in (?)", opIds).Scan(&privileges)
	return privileges, nil
}

func (o *OrgPrivilegeService) ConfigPrivileges(configPrivilegeDto ConfigPrivilegeDto) error{
	tx := common.GlobalDB.Begin()
	defer func() {
		if err := recover(); err != nil{
			tx.Callback()
		}
	}()
	var countNum int64
	var countOrgNum int64
	// 1,查询并判断是否对应组织数据,有则先删除,否则直接插入
	tx.Debug().Table("sys_org").Where("id=?",configPrivilegeDto.OrgId).Count(&countOrgNum)
	if countOrgNum == 0 {
		tx.Callback()
		return errors.New("添加权限失败")
	}
	tx.Debug().Table("sys_org_privilege").Where("org_id=?",configPrivilegeDto.OrgId).Count(&countNum)
	if countNum > 0 {
		// 2,删除对应orgId的数据
		rowsAffected := tx.Debug().Unscoped().Where("org_id=?",configPrivilegeDto.OrgId).Delete(&model.SysOrgPrivilege{}).RowsAffected
		if rowsAffected == 0 {
			tx.Callback()
			return errors.New("添加权限失败1")
		}
	}
	// 3,插入新权限数据
	for _,privilegeId := range configPrivilegeDto.PrivilegeIds{
		rowsAffected := tx.Debug().Create(&model.SysOrgPrivilege{OrgId: configPrivilegeDto.OrgId, PrivilegeId: privilegeId}).RowsAffected
		if rowsAffected == 0 {
			tx.Callback()
			return errors.New("添加权限失败2")
		}
	}
	tx.Commit()
	return nil
}