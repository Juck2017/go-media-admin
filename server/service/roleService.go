package service

import (
	"errors"
	"suitbim.com/go-media-admin/common"
	"suitbim.com/go-media-admin/model"
)

type RoleService struct{
	BaseService
}

// 返回数据体
type RoleDto struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`        // 角色名称
	OrgId       int `json:"orgId"`       // 组织id
	OrgName     string `json:"orgName"`     // 组织名称
	Code        string `json:"code"`        // 编码
	Description string `json:"description"` // 描述
}

// 请求数据体
type RoleRequestDto struct {
	PageSize int `json:"pageSize" binding:"required"`
	CurrPage int `json:"currPage" binding:"required"`
	Name     string
	Code     string
}

// 获取角色列表
func (s *RoleService) GetRoleList(requestDto RoleRequestDto) ([]RoleDto, int64, error) {
	var roleList []RoleDto
	var countNum int64
	dbData := common.GlobalDB.Debug().Table("sys_role").Select("sys_role.*, sys_org.name as org_name").
		Joins("left join sys_org on sys_role.org_id = sys_org.id")
	if requestDto.Name != ""{
		dbData = dbData.Where("sys_role.name = ?", requestDto.Name)
	}
	if requestDto.Code != ""{
		dbData = dbData.Where("sys_role.code = ?", requestDto.Code)
	}
	dbData.Count(&countNum)
	dbData.Offset((requestDto.CurrPage-1)*requestDto.PageSize).Limit(requestDto.PageSize).Scan(&roleList)
	if len(roleList) == 0 {
		return roleList, countNum, errors.New("暂无数据")
	}
	return roleList, countNum, nil
}

/**
	保存角色
 */
func (s *RoleService) SaveRole(roleData model.SysRole) (RoleDto,  error) {
	var roleDto RoleDto
	tx := common.GlobalDB.Begin()
	defer func() {
		if r := recover(); r != nil{
			tx.Callback()
		}
	}()
	var orgName string
	common.GlobalDB.Debug().Table("sys_org").Select("name").Where("id=?", roleData.OrgId).Scan(&orgName)
	if orgName == ""{
		tx.Callback()
		return roleDto, errors.New("添加失败")
	}
	roleModel,err := s.Save(&roleData)
	if err != nil {
		return roleDto, err
	}
	s.StructToStruct(&roleDto, roleModel.(*model.SysRole))
	roleDto.Id = roleData.ID
	roleDto.OrgName = orgName
	//if roleData.ID == 0 {
	//	// id为空为新增
	//	rowsAffected := tx.Debug().Create(&roleData).RowsAffected
	//	if rowsAffected == 0 {
	//		tx.Callback()
	//		return roleDto, errors.New("添加失败")
	//	}
	//	roleDto.Code = roleData.Code
	//	roleDto.Name = roleData.Name
	//	roleDto.OrgId = roleData.OrgId
	//	roleDto.Description = roleData.Description
	//	roleDto.Id = roleData.ID
	//	roleDto.OrgName = orgName
	//} else {
	//	// id不为空为修改
	//		rowsAffected :=tx.Debug().Model(&model.SysRole{}).Where("id=?", roleData.ID).
	//		Updates(model.SysRole{Name: roleData.Name, Code: roleData.Code,
	//			Description: roleData.Description, OrgId: roleData.OrgId}).RowsAffected
	//	if rowsAffected == 0 {
	//		return roleDto, errors.New("修改失败")
	//	}
	//	roleDto.Code = roleData.Code
	//	roleDto.Name = roleData.Name
	//	roleDto.OrgId = roleData.OrgId
	//	roleDto.Description = roleData.Description
	//	roleDto.Id = roleData.ID
	//	roleDto.OrgName = orgName
	//}
	//if tx.Commit().Error != nil{
	//	return roleDto, errors.New("添加失败")
	//}
	return roleDto, nil
}
/**
*/
func (s *RoleService) Delete(id uint) error {
	tx := common.GlobalDB.Begin()
	defer func() {
		if r := recover(); r != nil{
			tx.Callback()
		}
	}()
	// 删除时将角色对应权限一并删除
	if tx.Debug().Unscoped().Where("id=?",id).Delete(&model.SysRole{}).Error != nil {
		tx.Callback()
		return errors.New("删除失败")
	}
	if tx.Debug().Unscoped().Where("role_id=?",id).Delete(&model.SysRolePrivilege{}).Error != nil {
		tx.Callback()
		return errors.New("删除失败")
	}
	if tx.Commit().Error != nil {
		tx.Callback()
		return errors.New("删除失败")
	}
	return nil
}
