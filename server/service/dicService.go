package service

import (
	"errors"
	"suitbim.com/go-media-admin/common"
	"suitbim.com/go-media-admin/model"
)

type DicService struct {
	BaseService
}

type DicDto struct {
	Id uint  `json:"id"`
	Code string `json:"code"`
	Leaf  bool `json:"leaf"`
	Name string `json:"name"`
	ParentId *int32 `json:"parentId,omitempty"`
	TreeId string `json:"treeId"`
	TreeLevel int `json:"treeLevel"`
}

// 获取map类型的字典
func (d *DicService) GetDicMap() (map[string][]DicDto,error) {
	var dicMap = make(map[string][]DicDto,2)
	var dicLevel1 []model.SysDict
	var countNum int64
	// 获取一级字典名称
	common.GlobalDB.Table("sys_dict").Where("parent_id is null").Scan(&dicLevel1).Count(&countNum)
	if countNum <= 0 {
		return nil, errors.New("暂无查询到字典")
	}
	for _,dicLevel :=  range dicLevel1 {
		var dicLevel2 []DicDto
		// 获取二级字典名称
		common.GlobalDB.Table("sys_dict").Where("parent_id = ?",  dicLevel.ID).Scan(&dicLevel2)
		dicMap[dicLevel.Code]  = dicLevel2
	}
	return dicMap, nil
}

// 获取字典列表
func (d *DicService) List() ([]DicDto,error) {
	var dicList []DicDto
	var countNum int64
	// 获取一级字典名称
	common.GlobalDB.Table("sys_dict").Scan(&dicList).Count(&countNum)
	if countNum == 0 {
		return nil,errors.New("暂未查到数据")
	}
	return dicList,nil
}

//func generateTreeData(dicData model.SysDict) (model.SysDict,error) {
//	type temp struct {
//		ParentId string `json:"parent_id"`
//		TreeId string `json:"tree_id"`
//		TreeLevel int `json:"tree_level"`
//		Leaf int `json:"leaf"`
//		MaxTreeId string `json:"max_tree_id"`
//	}
//	var tempData temp
//	common.GlobalDB.Raw("select parent_id,tree_id,tree_level,SUBSTR('00000' || (SUBSTR(MAX(tree_id),-5)+1), -5) as max_tree_id from sys_dict where parent_id = ?;",
//		dicData.ParentId).Scan(&tempData)
//	// parentId为0,则说明是一级分类,则去查询数据库parentId为0的数据,没有则是无数据,有则取最大treeId+1
//	if dicData.ParentId == 0 {
//		if tempData.TreeId == "" {
//			// 说明数据库一条数据没有
//			dicData.TreeId = "00001"
//			return dicData, nil
//		} else {
//			// 说明有数据,treeId = treeId+1
//			dicData.TreeId = tempData.MaxTreeId
//			return dicData, nil
//		}
//	}
//	// parentId不为0,则先判断是否存在该parentId对应id的数据,不存在则返回数据不存在,存在则继续
//	// 判断parentId对应的数据是否存在,没有则新增,有则取treeId最大值
//	var idData temp
//	common.GlobalDB.Raw("select parent_id,tree_level,leaf,tree_id from sys_dict where id = ?;", dicData.ParentId).Scan(&idData)
//	if idData.TreeId == "" {
//		fmt.Println("暂无数据")
//		return dicData, errors.New("暂无数据")
//	}
//	dicData.Leaf = 1
//	dicData.TreeLevel =  idData.TreeLevel + 1
//	if tempData.TreeId == "" {
//		dicData.TreeId = idData.TreeId+".00001"
//		return dicData, nil
//	}
//	treeArr := strings.Split(tempData.TreeId, ".")
//	newTreeId := fmt.Sprintf("%s.%s", strings.Join(treeArr[0:len(treeArr)-1], "."), tempData.MaxTreeId)
//	dicData.TreeId = newTreeId
//	return dicData, nil
//}

// 添加或更新字典
//func (d *DicService) Save(dicData DicDto) (DicDto,error) {
	//tx := common.GlobalDB.Begin()
	//defer func() {
	//	if r := recover(); r != nil {
	//		fmt.Println("err:", r)
	//	}
	//}()
	//if dicData.Id == 0 {
	//	paramMap := map[string]interface{}{
	//		"tableName": "sys_dict",
	//		"parentId": dicData.ParentId,
	//		"operation": "save",
	//	}
	//	// id为空为新增
	//	resultMap, err := d.AsyncGenerateTreeData(tx, paramMap)
	//	if err != nil {
	//		tx.Callback()
	//		return dicData,err
	//	}
	//	dicData.TreeLevel = resultMap["treeLevel"].(int)
	//	dicData.Leaf = resultMap["leaf"].(bool)
	//	dicData.TreeId = resultMap["treeId"].(string)
	//
	//	saveDict := model.SysDict{}
	//	d.StructToStruct(&saveDict, &dicData)
	//	saveDict.Leaf = resultMap["leaf"].(bool)
	//	saveDict.Parent.Int32 = dicData.ParentId
	//	saveDict.Parent.Valid = true
	//	if saveDict.Parent.Int32 == 0 {
	//		saveDict.Parent.Valid = false
	//	}
	//	rowsAffected := tx.Debug().Create(&saveDict).RowsAffected
	//	if rowsAffected == 0 {
	//		return dicData, errors.New("添加失败")
	//	}
	//	dicData.Id = saveDict.ID
	//	// 更新leaf状态
	//	err = d.UpdateLeafStatus(tx, paramMap)
	//	if err != nil {
	//		tx.Callback()
	//		fmt.Println(err)
	//		return dicData, errors.New("添加或更新数据失败")
	//	}
	//} else {
	//	// id不为空为修改
	//	rowsAffected := tx.Debug().Model(&model.SysDict{}).Where("id=? and parent_id=?", dicData.Id, dicData.ParentId).
	//		Updates(model.SysDict{Name: dicData.Name, Code: dicData.Code}).RowsAffected
	//	if rowsAffected == 0 {
	//		return dicData, errors.New("修改失败")
	//	}
	//}
	//if err := tx.Commit().Error;err != nil{
	//	fmt.Println(err)
	//	return dicData, errors.New("添加或更新数据失败")
	//}
//	return dicData, nil
//}

// 删除字典
//func (d *DicService) Delete(id uint) error {
//	// 字典下面有叶子分支不让删除
//	// 根据id查询parentId是否有数据
//	tx := common.GlobalDB.Begin()
//	var count int64
//	tx.Debug().Model(&model.SysDict{}).Where("parent_id=?",id).Count(&count)
//	if count > 0{
//		return errors.New("当前节点存在子节点")
//	}
//	// 查询当前id的parentId
//	var sysDict model.SysDict
//	if err := tx.Debug().Model(&model.SysDict{}).Where("id=?", id).Scan(&sysDict).Error; err != nil{
//		return errors.New("删除失败")
//	}
//	rowsAffected := tx.Debug().Unscoped().Where("id=?",id).Delete(&model.SysDict{}).RowsAffected
//	if rowsAffected == 0 {
//		return errors.New("删除失败")
//	}
//
//	var n int32 = 0
//	if sysDict.ParentId == nil {
//		sysDict.ParentId = &n
//		fmt.Println(*sysDict.ParentId)
//	}
//	paramMap := map[string]interface{}{
//		"tableName": "sys_dict",
//		"parentId": *sysDict.ParentId,
//		"operation": "delete",
//	}
//	// 当前id对应的parentId数据下无子节点,则更新leaf状态为1
//	err := d.UpdateLeafStatus(tx, paramMap)
//	if err != nil {
//		tx.Callback()
//		fmt.Println(err)
//		return errors.New("删除失败")
//	}
//	if err = tx.Commit().Error;err != nil{
//		fmt.Println(err)
//		return errors.New("删除失败")
//	}
//	return nil
//}