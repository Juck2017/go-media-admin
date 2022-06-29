package service

import (
	"errors"
	"suitbim.com/go-media-admin/common"
	"suitbim.com/go-media-admin/model"
)

type RolePrivilegeDto struct {
	Id uint `json:"id"`
	RoleId int `json:"roleId"` // 角色id
	PrivilegeId int `json:"PrivilegeId"` // 权限id
}

type ConfigRolePrivilegeDto struct {
	RoleId int `json:"roleId" binding:"required"`
	PrivilegeIds []int `json:"privilegeIds" binding:"required"`
}

type RolePrivilegeService struct {}

// 获取角色权限列表
func (o *RolePrivilegeService) GetRolePrivilegeList(roleId int)  ([]PrivilegeDto,error) {
	var rolePrivileges []RolePrivilegeDto
	var privileges []PrivilegeDto
	var countRNum int64
	common.GlobalDB.Table("sys_role_privilege").Where("role_id=?", roleId).Scan(&rolePrivileges).Count(&countRNum)
	if countRNum == 0 {
		return privileges, errors.New("暂无数据")
	}
	for _,rolePrivilege := range rolePrivileges{
		var privilege PrivilegeDto
		common.GlobalDB.Table("sys_privilege").Where("id=?", rolePrivilege.PrivilegeId).Scan(&privilege)
		privileges = append(privileges, privilege)
	}
	if len(privileges) == 0 {
		return privileges, errors.New("暂无数据")
	}
	return privileges, nil
}

// 配置角色权限
func (o *RolePrivilegeService) ConfigRolePrivileges(configRolePrivilegeDto ConfigRolePrivilegeDto) error{
	tx := common.GlobalDB.Begin()
	defer func() {
		if err := recover(); err != nil{
			tx.Callback()
		}
	}()
	var countNum int64
	var countRoleNum int64
	// 1,查询并判断是否对应角色数据,有则先删除,否则直接插入
	tx.Debug().Table("sys_role").Where("id=?",configRolePrivilegeDto.RoleId).Count(&countRoleNum)
	if countRoleNum == 0 {
		tx.Callback()
		return errors.New("添加权限失败")
	}
	tx.Debug().Table("sys_role_privilege").Where("role_id=?",configRolePrivilegeDto.RoleId).Count(&countNum)
	if countNum > 0 {
		// 2,删除对应roleId的数据
		rowsAffected := tx.Debug().Unscoped().Where("role_id=?",configRolePrivilegeDto.RoleId).Delete(&model.SysRolePrivilege{}).RowsAffected
		if rowsAffected == 0 {
			tx.Callback()
			return errors.New("添加权限失败")
		}
	}
	// 3,插入新权限数据
	for _,privilegeId := range configRolePrivilegeDto.PrivilegeIds{
		rowsAffected := tx.Debug().Create(&model.SysRolePrivilege{RoleId: configRolePrivilegeDto.RoleId, PrivilegeId: privilegeId}).RowsAffected
		if rowsAffected == 0 {
			tx.Callback()
			return errors.New("添加权限失败")
		}
	}
	if tx.Commit().Error != nil{
		tx.Callback()
		return errors.New("添加权限失败")
	}
	return nil
}