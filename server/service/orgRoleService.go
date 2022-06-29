package service

import (
	"fmt"
	"suitbim.com/go-media-admin/common"
	"suitbim.com/go-media-admin/model"
)

type RoleDataDto struct {
	Id uint `json:"id"`
	OrgId int `json:"orgId"` // 组织id
	Name string `json:"name"` // 角色名称
	Code string `json:"code"` // 角色编码
}

// 组织角色列表
type OrgRoleService struct {
	BaseService
}

// 获取组织角色列表
func (o *OrgRoleService) GetOrgRoleList(orgId uint64)  ([]RoleDataDto,error) {
	var roles []RoleDataDto
	// 在组织表中根据当前组织id递归去寻找父组织id,然后去角色表根据组织id寻找对应的角色列表
	var orgList []model.SysOrg
	var countOrgNum int64
	common.GlobalDB.Table("sys_org").Debug().Scan(&orgList).Count(&countOrgNum)
	if countOrgNum > 0{
		arrList := make([]interface{}, len(orgList))
		for i := 0; i <  len(orgList); i++ {
			arrList[i] = &orgList[i]
		}
		parentOrgId, err := o.GetTreeParentId(arrList, []uint64{},orgId)
		fmt.Println(parentOrgId)
		if err != nil {
			return roles, err
		}
		common.GlobalDB.Table("sys_role").Debug().Where("org_id=?", parentOrgId).Scan(&roles)
	}
	return roles, nil
}

//func treeParentId(orgList []model.SysOrg, ids []uint, orgId uint) (uint,error) {
//	if len(orgList) > 0{
//		for _, org := range orgList{
//			if org.ParentId != nil && org.ID == orgId {
//				ids = append(ids, orgId)
//				return treeParentId(orgList,ids, uint(*org.ParentId))
//			}
//			if org.ParentId == nil && org.ID == orgId{
//				ids = append(ids, org.ID)
//				return treeParentId(orgList,ids, 0)
//			}
//		}
//	} else {
//		return 0, errors.New("组织树数据不能为空")
//	}
//	return ids[len(ids)-1], nil
//}