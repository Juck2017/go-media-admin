package common

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"suitbim.com/go-media-admin/model"
)

// 初始化db
func InitDb()  {
	db,err := gorm.Open(sqlite.Open(GlobalConf.App.Db.Path),&gorm.Config{})
	if err != nil {
		GlobalConf.Logger.Errorf("Open Database [%s] Failed: %v", GlobalConf.App.Db.Path, err)
		return
	}
	GlobalDB = db
	if GlobalConf.App.Db.DdlAuto == "update"{
		db.AutoMigrate(
			&model.SysUser{},
			&model.SysPrivilege{},
			&model.SysDict{},
			&model.SysOrg{},
			&model.SysRolePrivilege{},
			&model.SysRole{},
			&model.SysUserRole{},
			&model.SysUserPrivilege{},
			&model.SysOrgPrivilege{},
			&model.Channel{},
			&model.Pusher{},
		)
		GlobalConf.Logger.Infof("==== Global update database table struct success ====")
	}
	GlobalConf.Logger.Infof("==== Global init database success ====")
}
