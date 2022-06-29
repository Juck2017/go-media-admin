package service

import (
	"errors"
	"suitbim.com/go-media-admin/common"
	"suitbim.com/go-media-admin/model"
)

// 配置用户角色请求体
type ConfigRoleDto struct {
	UserId int `json:"userId" binding:"required"`
	RoleIds []int `json:"roleIds" binding:"required"`
}

// 用户角色列表
type UserRoleService struct {
	BaseService
}

// 获取用户角色列表
func (o *UserRoleService) GetUserRoleList(userId uint)  ([]RoleDataDto,error) {
	var roles []RoleDataDto
	common.GlobalDB.Table("sys_role").Debug().Select("sys_role.*").
		Joins("LEFT JOIN sys_user_role on sys_user_role.role_id = sys_role.id").
		Where("sys_user_role.user_id =?", userId).Scan(&roles)
	return roles, nil
}

func (o *UserRoleService) ConfigRoles(configRole ConfigRoleDto) error {
	// 先查询用户角色表中是否有当前用户id的数据,根据用户id删除用户角色表中数据,再重新插入新的数据
	tx := common.GlobalDB.Begin()
	defer func() {
		if r := recover(); r != nil{
			tx.Callback()
		}
	}()
	var rowNum int64
	tx.Debug().Table("sys_user_role").Where("user_id = ?", configRole.UserId).Count(&rowNum)
	if rowNum > 0 {
		if err := tx.Debug().Where("user_id = ?", configRole.UserId).Unscoped().Delete(&model.SysUserRole{}).Error; err!= nil{
			tx.Callback()
			return err
		}
	}
	var userRoles []model.SysUserRole
	for _, roleId := range configRole.RoleIds{
		userRole := model.SysUserRole{
			UserId: configRole.UserId,
			RoleId: roleId,
		}
		userRoles = append(userRoles, userRole)
	}
	if len(userRoles) > 0 {
		var rowAffected = tx.Debug().Create(&userRoles).RowsAffected
		if rowAffected == 0 {
			tx.Callback()
			return errors.New("配置用户角色失败")
		}
	}
	if tx.Commit().Error != nil{
		tx.Callback()
		return errors.New("配置用户角色失败")
	}
	return nil
}