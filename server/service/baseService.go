package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"strings"
	"suitbim.com/go-media-admin/common"
)

var dbCommon  = *common.GlobalDB
// 基本service结构体
type BaseService struct {
}

// 更新父节点状态
func (b *BaseService) UpdateLeafStatus(db *gorm.DB ,mapData map[string]interface{}) error{
	tableName := mapData["tableName"].(string)
	parentId := mapData["parentId"].(int32)
	operation := mapData["operation"]
	// parentId不为0,则说明需要更新
	if parentId != 0 {
		// 1,当前节点从无子节点变为有子节点,则leaf由1变为0(新增)
		if operation == "save" {
			rowsAffected := db.Debug().Table(tableName).Where("id=?", parentId).Update("leaf", false).RowsAffected
			if rowsAffected == 0 {
				fmt.Println("rowsAffected:0")
				return errors.New("影响数据为0")
			}
			fmt.Println("更新父节点状态成功")
			return nil
		} else if operation == "delete" {
			var num int64
			err := db.Debug().Table(tableName).Where("parent_id=?", parentId).Count(&num).Error
			if err != nil {
				fmt.Println("1:",err)
				return err
			}
			if num == 0 {
				// 2,当前节点从有子节点变为无子节点,则leaf由0变为1(删除)
				err = db.Debug().Table(tableName).Where("id=?", parentId).Update("leaf", true).Error
				if err != nil {
					fmt.Println("3:",err)
					return err
				}
				fmt.Println("更新父节点状态成功")
			}
		}
	}
	return nil
}

// 动态生成树数据,参数mapData中包含(tableName表名,parentId)
func (b *BaseService) AsyncGenerateTreeData(db *gorm.DB ,mapData map[string]interface{}) (map[string]interface{},error) {
	type temp struct {
		ParentId int `json:"parent_id,omitempty"`
		TreeId string `json:"tree_id"`
		TreeLevel int `json:"tree_level"`
		Leaf bool `json:"leaf"`
		MaxTreeId string `json:"max_tree_id"`
	}
	var tempData temp
	parentId := mapData["parentId"].(int32)
	tableName := mapData["tableName"]
	parentSql := fmt.Sprintf("select parent_id,tree_id,tree_level,SUBSTR('00000' || (SUBSTR(MAX(tree_id),-5)+1), -5) as max_tree_id from %s where parent_id ", tableName)
	if parentId == 0 {
		parentSql += "is null;"
	} else {
		parentSql += fmt.Sprintf("=%v",parentId)
	}
	db.Debug().Raw(parentSql).Scan(&tempData)
	// 1,parentId为null,则说明是一级分类,则去查询数据库parentId为null的数据,没有则是无数据,有则取最大treeId+1
	resultMap := make(map[string]interface{}, 5)
	if parentId == 0 {
		resultMap["treeLevel"] = 1
		resultMap["leaf"] = true
		if tempData.TreeId == "" {
			// 说明数据库一条数据没有
			resultMap["treeId"] =  "00001"
			return resultMap, nil
		} else {
			// 说明有数据,treeId = treeId+1
			resultMap["treeId"] =  tempData.MaxTreeId
			return resultMap, nil
		}
	}
	// 2,parentId不为null,则先判断是否存在该parentId对应id的数据,不存在则返回数据不存在,存在则继续
	// 判断parentId对应的数据是否存在,没有则新增,有则取treeId最大值
	var idData temp
	rawSql := fmt.Sprintf("select parent_id,tree_level,leaf,tree_id from %s where id = %d;", tableName, parentId)
	db.Debug().Raw(rawSql).Scan(&idData)
	if idData.TreeId == "" {
		fmt.Println("暂无数据")
		return nil, errors.New("暂无数据")
	}
	resultMap["leaf"] = true
	resultMap["treeLevel"] = idData.TreeLevel + 1
	if tempData.TreeId == "" {
		resultMap["treeId"] =  idData.TreeId+".00001"
		return resultMap, nil
	}
	treeArr := strings.Split(tempData.TreeId, ".")
	newTreeId := fmt.Sprintf("%s.%s", strings.Join(treeArr[0:len(treeArr)-1], "."), tempData.MaxTreeId)
	resultMap["treeId"] =  newTreeId
	return resultMap, nil
}

/**
	结构体相同字段赋值
	binding: 需要赋值的结构体指针
	data: 已经有数据的结构体指针
*/
func (b *BaseService) StructToStruct(binding interface{}, data interface{}) {
	bValue := reflect.ValueOf(binding).Elem()
	dValue := reflect.ValueOf(data).Elem()
	for i:=0;i<dValue.NumField();i++ {
		name := dValue.Type().Field(i).Name
		if ok := bValue.FieldByName(name).IsValid(); ok {
			bValue.FieldByName(name).Set(reflect.ValueOf(dValue.Field(i).Interface()))
		}
	}
}

/**
	保存和更新树结构,参数为结构体接口指针
*/
func (b *BaseService) SaveTreeData(dataModel interface{}) (interface{},error) {
	tx := common.GlobalDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("err:", r)
		}
	}()
	// 通过反射获取id, parentId, tableName
	id := reflect.ValueOf(dataModel).Elem().FieldByName("ID").Uint()
	parentId := reflect.ValueOf(dataModel).Elem().FieldByName("ParentId").Interface().(*int32)
	tableName := reflect.ValueOf(dataModel).MethodByName("TableName").Call(make([]reflect.Value, 0))[0].Interface().(string)
	if id == 0 {
		var n int32 = 0
		if parentId == nil {
			parentId = &n
		}
		paramMap := map[string]interface{}{
			"tableName": tableName,
			"parentId": *parentId,
			"operation": "save",
		}
		// id为空为新增
		resultMap, err := b.AsyncGenerateTreeData(tx, paramMap)
		if err != nil {
			tx.Callback()
			return dataModel, err
		}
		treeLevel := resultMap["treeLevel"].(int)
		leaf := resultMap["leaf"].(bool)
		treeId := resultMap["treeId"].(string)
		// 通过反射设置treeLevel,leaf,treeId的值
		reflect.ValueOf(dataModel).Elem().FieldByName("TreeLevel").Set(reflect.ValueOf(treeLevel))
		reflect.ValueOf(dataModel).Elem().FieldByName("Leaf").Set(reflect.ValueOf(leaf))
		reflect.ValueOf(dataModel).Elem().FieldByName("TreeId").Set(reflect.ValueOf(treeId))
		rowsAffected := tx.Debug().Create(dataModel).RowsAffected
		if rowsAffected == 0 {
			return dataModel,errors.New("添加失败")
		}
		// 更新leaf状态
		if *parentId != 0{
			err = b.UpdateLeafStatus(tx, paramMap)
			if err != nil {
				tx.Callback()
				fmt.Println(err)
				return dataModel, errors.New("添加或更新数据失败")
			}
		}
	} else {
		sqlStr := "id="+strconv.Itoa(int(id))+" and parent_id "
		if parentId == nil {
			sqlStr += "is null"
		} else {
			sqlStr += "="+ strconv.Itoa(int(*parentId))
		}
		// id不为空为修改
		rowsAffected := tx.Debug().Model(dataModel).Where(sqlStr).Updates(dataModel).RowsAffected
		if rowsAffected == 0 {
			return dataModel, errors.New("修改失败")
		}
		if err := tx.Debug().Model(dataModel).Scan(dataModel).Error; err != nil {
			return dataModel, err
		}
	}
	if err := tx.Commit().Error;err != nil{
		fmt.Println(err)
		return dataModel, errors.New("添加或更新数据失败")
	}
	return dataModel, nil
}

//	删除树操作,默认根据id删除,参数为模型指针
func (b *BaseService) DeleteTreeData(dataModel interface{}) error {
	// 删除前判断是否有parentId为当前id的,有的话则不让删除,没有的话则可以删除
	tx := common.GlobalDB.Begin()
	var countNum int64
	// 通过反射获取Id
	id := reflect.ValueOf(dataModel).Elem().FieldByName("ID").Uint()
	tx.Debug().Model(dataModel).Where("parent_id=?",id).Count(&countNum)
	if countNum > 0{
		return errors.New("当前节点存在子节点")
	}
	// 查询当前id的parentId
	if err := tx.Debug().Model(dataModel).Scan(dataModel).Error; err != nil{
		return errors.New("删除失败")
	}
	rowsAffected := tx.Debug().Unscoped().Delete(dataModel).RowsAffected
	if rowsAffected == 0 {
		return errors.New("删除失败")
	}
	// 通过反射获取id, parentId, tableName
	parentId := reflect.ValueOf(dataModel).Elem().FieldByName("ParentId").Interface().(*int32)
	tableName := reflect.ValueOf(dataModel).MethodByName("TableName").Call(make([]reflect.Value, 0))[0].Interface().(string)
	// 当前id对应的parentId数据下无子节点,则更新leaf状态为1, 如果leaf是true,则不进行更新操作
	if parentId != nil{
		paramMap := map[string]interface{}{
			"tableName": tableName,
			"parentId": *parentId,
			"operation": "delete",
		}
		err := b.UpdateLeafStatus(tx, paramMap)
		if err != nil {
			tx.Callback()
			fmt.Println(err)
			return errors.New("删除失败")
		}
	}

	if err := tx.Commit().Error;err != nil{
		fmt.Println(err)
		return errors.New("删除失败")
	}
	return nil
}

// 添加或更新方法
func (b *BaseService) Save(dataModel interface{}) (interface{},  error) {
	tx := common.GlobalDB.Begin()
	defer func() {
		if r := recover(); r != nil{
			tx.Callback()
		}
	}()
	id := reflect.ValueOf(dataModel).Elem().FieldByName("ID").Uint()
	if id == 0 {
		// id为空为新增
		rowsAffected := tx.Debug().Create(dataModel).RowsAffected
		if rowsAffected == 0 {
			tx.Callback()
			return dataModel, errors.New("添加失败")
		}
	} else {
		// id不为空为修改
		rowsAffected :=tx.Debug().Model(dataModel).Where("id=?", id).
			Updates(dataModel).RowsAffected
		if rowsAffected == 0 {
			return dataModel, errors.New("修改失败")
		}
	}
	if tx.Commit().Error != nil{
		return dataModel, errors.New("添加失败")
	}
	return dataModel, nil
}

/**
	根据当前id递归获取当前树的父id
*/
func  (b *BaseService) GetTreeParentId(modelList []interface{}, ids []uint64, orgId uint64) (uint64,error) {
	if len(modelList) > 0{
		for _, dataModel := range modelList{
			id := reflect.ValueOf(dataModel).Elem().FieldByName("ID").Uint()
			parentId := reflect.ValueOf(dataModel).Elem().FieldByName("ParentId").Interface().(*int32)
			if parentId != nil && id == orgId {
				ids = append(ids, orgId)
				return b.GetTreeParentId(modelList, ids, uint64(*parentId))
			}
			if parentId == nil && id == orgId{
				ids = append(ids, id)
				return b.GetTreeParentId(modelList,ids, 0)
			}
		}
	} else {
		return 0, errors.New("组织树数据不能为空")
	}
	return ids[len(ids)-1], nil
}