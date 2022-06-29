package service

import (
	"errors"
	"suitbim.com/go-media-admin/common"
)

type OrgService struct {
	BaseService
}

type OrgDto struct {
	Id uint `json:"id"`
	Code string `json:"code"` // 组织编码
	Name string `json:"name"` // 组织名称
	Leaf bool `json:"leaf"` // leaf: 0或者1
	ParentId *int32 `json:"parentId,omitempty"` // 父id
	TreeId string `json:"treeId"` // 树id
	TreeLevel int `json:"treeLevel"` // 树级别
	Type string `json:"type"` // 组织类型,company代表公司,dept代表部门
}

func (o *OrgService) GetOrgList()  ([]OrgDto,error) {
	var orgList []OrgDto
	var countNum int64
	common.GlobalDB.Table("sys_org").Scan(&orgList).Count(&countNum)
	if countNum == 0 {
		return orgList, errors.New("暂无数据")
	}
	return orgList, nil
}

// 获取根组织列表
func (o *OrgService) GetRootOrgList() ([]OrgDto,error) {
	var rootOrgList []OrgDto
	common.GlobalDB.Table("sys_org").Where("parent_id is null").Scan(&rootOrgList)
	if len(rootOrgList) == 0 {
		return rootOrgList,errors.New("暂无数据")
	}
	return rootOrgList, nil
}

//func (o *OrgService) Save(sysOrgData model.SysOrg) (OrgDto,error) {
//	tx := common.GlobalDB.Begin()
//	defer func() {
//		if r := recover(); r != nil {
//			tx.Callback()
//		}
//	}()
//	var orgDto OrgDto
//	paramMap := map[string]interface{}{
//		"tableName": "sys_org",
//		"parentId": *sysOrgData.ParentId,
//	}
//	if sysOrgData.ID == 0 {
//		// id为空为新增
//		resultMap, err := o.AsyncGenerateTreeData(tx,paramMap)
//		if err != nil {
//			return orgDto,err
//		}
//		orgDto.Name = sysOrgData.Name
//		orgDto.Code = sysOrgData.Code
//		orgDto.ParentId = sysOrgData.ParentId
//		orgDto.TreeLevel = resultMap["treeLevel"].(int)
//		orgDto.Leaf = resultMap["leaf"].(bool)
//		orgDto.TreeId = resultMap["treeId"].(string)
//		orgDto.Type = sysOrgData.Type
//		sysOrgData.TreeId = resultMap["treeId"].(string)
//		sysOrgData.TreeLevel = resultMap["treeLevel"].(int)
//		sysOrgData.Leaf = resultMap["leaf"].(bool)
//		rowsAffected := common.GlobalDB.Debug().Create(&sysOrgData).RowsAffected
//		if rowsAffected == 0 {
//			return orgDto, errors.New("添加失败")
//		}
//		orgDto.Id = sysOrgData.ID
//	} else {
//		// id不为空为修改
//		rowsAffected := common.GlobalDB.Debug().Model(&model.SysOrg{}).Where("id=? and parent_id=?", sysOrgData.ID, sysOrgData.ParentId).
//			Updates(model.SysOrg{Name: sysOrgData.Name, Code: sysOrgData.Code,
//				Type: sysOrgData.Type}).RowsAffected
//		if rowsAffected == 0 {
//			return orgDto, errors.New("修改失败")
//		}
//		orgDto.Id = sysOrgData.ID
//		orgDto.Name = sysOrgData.Name
//		orgDto.Code = sysOrgData.Code
//		orgDto.Type = sysOrgData.Type
//		orgDto.ParentId = sysOrgData.ParentId
//	}
//	return orgDto, nil
//}

//func (o *OrgService) Delete(id uint) error {
//	// 删除前判断是否有parentId为当前id的,有的话则不让删除,没有的话则可以删除
//	var countNum int64
//	common.GlobalDB.Debug().Model(&model.SysOrg{}).Where("parent_id=?",id).Count(&countNum)
//	if countNum > 0{
//		return errors.New("当前节点存在子节点")
//	}
//	rowsAffected := common.GlobalDB.Debug().Unscoped().Where("id=?",id).Delete(&model.SysOrg{}).RowsAffected
//	if rowsAffected == 0 {
//		return errors.New("删除失败")
//	}
//	return nil
//}

