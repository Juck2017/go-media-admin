package service

import (
	"errors"
	"suitbim.com/go-media-admin/common"
)

type PrivilegeService struct {
	BaseService
}

type PrivilegeDto struct {
	Id uint  `json:"id"`
	Code string `json:"code"` // 权限编码
	Name string `json:"name"` // 权限名称
	Leaf bool `json:"leaf"`
	ParentId *int32 `json:"parentId,omitempty"` // 父id,omitempty忽略空值
	TreeId string `json:"treeId"` // 树id
	TreeLevel int `json:"treeLevel"` // 树级别
	Description string `json:"description"` // 权限描述
	Scope string `json:"scope"` // system权限范围,system系统,project项目
	Type string `json:"type"` // 类型:menu代表菜单,button代表按钮
}

// 获取权限列表
func (p *PrivilegeService) GetPrivilegeList() ([]PrivilegeDto,error) {
	var privilegeList []PrivilegeDto
	var countNum int64
	common.GlobalDB.Table("sys_privilege").Scan(&privilegeList).Count(&countNum)
	if countNum == 0 {
		return privilegeList, errors.New("暂无数据")
	}
	return privilegeList, nil
}

// 新增或者更新
//func (p *PrivilegeService) Save(privilegeDto PrivilegeDto) (PrivilegeDto,error) {
//	//var privilegeDto PrivilegeDto
//	tx := common.GlobalDB.Begin()
//	defer func() {
//		if r := recover(); r != nil {
//			fmt.Println("err:", r)
//		}
//	}()
//	paramMap := map[string]interface{}{
//		"tableName": "sys_privilege",
//		"parentId": privilegeDto.ParentId,
//		"operation": "save",
//	}
//	if privilegeDto.Id == 0 {
//		// id为空为新增
//		resultMap, err := p.AsyncGenerateTreeData(tx,paramMap)
//		if err != nil {
//			tx.Callback()
//			return privilegeDto,err
//		}
//		privilegeDto.TreeLevel = resultMap["treeLevel"].(int)
//		privilegeDto.Leaf = resultMap["leaf"].(bool)
//		privilegeDto.TreeId = resultMap["treeId"].(string)
//
//		savePrivilege := model.SysPrivilege{}
//		p.StructToStruct(&savePrivilege, &privilegeDto)
//		savePrivilege.Parent.Int32 = privilegeDto.ParentId
//		savePrivilege.Parent.Valid = true
//		if savePrivilege.Parent.Int32 == 0 {
//			savePrivilege.Parent.Valid = false
//		}
//		rowsAffected := tx.Debug().Create(&savePrivilege).RowsAffected
//		if rowsAffected == 0 {
//			tx.Callback()
//			return privilegeDto, errors.New("添加失败")
//		}
//		privilegeDto.Id = savePrivilege.ID
//		// 更新leaf状态
//		err = p.UpdateLeafStatus(tx, paramMap)
//		if err != nil {
//			tx.Callback()
//			fmt.Println(err)
//			return privilegeDto, errors.New("添加或更新数据失败")
//		}
//	} else {
//		// id不为空为修改
//		rowsAffected := tx.Debug().Model(&model.SysPrivilege{}).Where("id=? and parent_id=?", privilegeDto.Id, privilegeDto.ParentId).
//			Updates(model.SysPrivilege{Name: privilegeDto.Name, Code: privilegeDto.Code,
//				Type: privilegeDto.Type,Scope: privilegeDto.Scope, Description: privilegeDto.Description}).RowsAffected
//		if rowsAffected == 0 {
//			tx.Callback()
//			return privilegeDto, errors.New("修改失败")
//		}
//	}
//
//	if err := tx.Commit().Error;err != nil{
//		fmt.Println(err)
//		return privilegeDto, errors.New("添加或更新数据失败")
//	}
//	return privilegeDto, nil
//}

// 删除
//func (p *PrivilegeService) Delete(id uint) error {
//	// 删除前判断是否有parentId为当前id的,有的话则不让删除,没有的话则可以删除
//	tx := common.GlobalDB.Begin()
//	var countNum int64
//	tx.Debug().Model(&model.SysPrivilege{}).Where("parent_id=?",id).Count(&countNum)
//	if countNum > 0{
//		return errors.New("当前节点存在子节点")
//	}
//	// 查询当前id的parentId
//	var sysPrivilege model.SysPrivilege
//	if err := tx.Debug().Model(&model.SysPrivilege{}).Where("id=?", id).Scan(&sysPrivilege).Error; err != nil{
//		return errors.New("删除失败")
//	}
//	rowsAffected := tx.Debug().Unscoped().Where("id=?",id).Delete(&model.SysPrivilege{}).RowsAffected
//	if rowsAffected == 0 {
//		return errors.New("删除失败")
//	}
//	// 当前id对应的parentId数据下无子节点,则更新leaf状态为1, 如果leaf是true,则不进行更新操作
//	if sysPrivilege.ParentId != nil{
//		paramMap := map[string]interface{}{
//			"tableName": "sys_privilege",
//			"parentId": *sysPrivilege.ParentId,
//			"operation": "delete",
//		}
//		err := p.UpdateLeafStatus(tx, paramMap)
//		if err != nil {
//			tx.Callback()
//			fmt.Println(err)
//			return errors.New("删除失败")
//		}
//	}
//
//	if err := tx.Commit().Error;err != nil{
//		fmt.Println(err)
//		return errors.New("删除失败")
//	}
//	return nil
//}