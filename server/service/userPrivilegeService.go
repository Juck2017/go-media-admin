package service

import (
	"errors"
	"suitbim.com/go-media-admin/common"
	"suitbim.com/go-media-admin/model"
	"sort"
)

// 配置用户权限请求体
type UserPrivilegeDto struct {
	UserId uint `json:"userId" binding:"required"` // 用户id
	PrivilegeIds []uint `json:"privilegeIds" binding:"required"` // 权限ids
}

type UserPrivilegeService struct {
	BaseService
}

// 获取用户权限列表
func (o *UserPrivilegeService) GetUserPrivilegeList(userId uint64)  ([]PrivilegeDto,error) {
	/**
		1、首先从用户权限表中查询用户对应的权限id列表,如果有则再查询对应权限数据返回直接返回,否则执行下面的操作
		2、根据用户id查询用户对应的角色列表
		3、根据用户角色列表查询角色对应的权限并且合并
		4、根据权限列表的id查询对应的权限列表返回
	 */
	var privileges []PrivilegeDto
	var roleIds []uint
	var privilegeIds []int
	// 1、首先从用户权限表中查询用户对应的权限id列表,如果没有则执行下面的操作,有则直接跳过返回
	common.GlobalDB.Debug().Table("sys_user_privilege").Select("privilege_id").Where("user_id=?", userId).Scan(&privilegeIds)
	if len(privilegeIds) == 0 {
		// 2、根据用户id查询用户对应的角色列表
		common.GlobalDB.Debug().Table("sys_user_role").Select("role_id").Where("user_id=?", userId).Scan(&roleIds)

		// 3、根据用户角色列表查询角色对应的权限并且合并
		common.GlobalDB.Debug().Table("sys_role_privilege").Select("privilege_id").Where("role_id in (?)", roleIds).Scan(&privilegeIds)
		tMap := make(map[int]int)
		for _, v := range privilegeIds {
			tMap[v] = 1
		}
		privilegeIds = []int{}
		for key,_ := range tMap{
			privilegeIds = append(privilegeIds, key)
		}
		sort.Ints(privilegeIds) // 排序
	}

	// 4、根据权限列表的id查询对应的权限列表返回
	common.GlobalDB.Debug().Table("sys_privilege").Where("id in (?)", privilegeIds).Scan(&privileges)
	return privileges, nil
}

// 配置用户权限列表
func (o *UserPrivilegeService) ConfigUserPrivileges(userPrivilegeDto UserPrivilegeDto) error{
	/**
		1、查询并判断用户权限表中是否有对应的用户的权限
		2、有的话则先根据用户id删除对应用户的权限数据
		3、根据前端传过来的权限列表构造新的数据重新插入数据
	 */
	tx := common.GlobalDB.Begin()
	defer func() {
		if err := recover(); err != nil{
			tx.Callback()
		}
	}()
	var userPrivilegeNum int64
	// 1、查询并判断用户权限表中是否有对应的用户的权限id列表
	tx.Debug().Table("sys_user_privilege").Where("user_id=?",userPrivilegeDto.UserId).Count(&userPrivilegeNum)
	// 2、有的话则先根据用户id删除对应用户的权限数据
	if userPrivilegeNum > 0 {
		rowsAffected := tx.Debug().Table("sys_user_privilege").Where("user_id=?",userPrivilegeDto.UserId).Unscoped().
			Delete(&model.SysUserPrivilege{}).RowsAffected
		if rowsAffected == 0 {
			tx.Callback()
			return errors.New("添加用户权限失败")
		}
	}
	// 3、根据前端传过来的权限列表构造新的数据重新插入数据
	var userPrivileges []model.SysUserPrivilege
	for _, userPrivilegeId := range userPrivilegeDto.PrivilegeIds{
		userPrivilege := model.SysUserPrivilege{
			Deny: 0,
			PrivilegeId: userPrivilegeId,
			UserId: userPrivilegeDto.UserId,
		}
		userPrivileges = append(userPrivileges, userPrivilege)
	}
	rowsAffected := tx.Debug().Create(&userPrivileges).RowsAffected
	if rowsAffected == 0 {
		tx.Callback()
		return errors.New("添加用户权限失败")
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}