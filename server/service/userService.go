package service

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"suitbim.com/go-media-admin/common"
	"suitbim.com/go-media-admin/model"
	"time"
)

type UserService struct {
	BaseService
}

// 自定义时间格式
type FormatTime struct {
	time.Time
}

// 在此方法中实现自定义格式的转换
func (t *FormatTime) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(output), nil
}

// 写入数据库时会调用该方法将自定义时间类型转换并写入数据库
func (t *FormatTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// 读取数据库时会调用该方法将时间数据转换成自定义时间类型
func (t *FormatTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = FormatTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}


type UserDto struct {
	ID        uint `json:"id"`
	CreatedAt *FormatTime `json:"createdAt"`
	Name string `json:"name"` // 姓名
	Username string `json:"username"` // 用户名
	OrgId uint `json:"orgId"` // 组织id
	OrgName string `json:"orgName"` // 组织名称
	Mobile string `json:"mobile"` // 手机号
	Enabled bool `json:"enabled"` // 禁用状态1,可以,0,禁用
	Photo string `json:"photo"` // 用户头像地址
}

type UserRequestDto struct {
	PageSize int `json:"pageSize" binding:"required"`
	CurrPage int `json:"currPage" binding:"required"`
	Name     string `json:"name"`
	Mobile  string `json:"mobile"`
	OrgId int	`json:"orgId"`
	UserId uint	`json:"userId"`
}

// 时间转时间戳
//func timeToTimeStamp(param time.Time) string{
//	c,_ := time.ParseInLocation("2006-01-02 15:04:05",param.Format("2006-01-02 15:04:05"),time.Local)
//	return string(c.Unix() * 1000)
//}

// 通过用户名获取用户信息
func (u *UserService) GetUserInfoByUsername(username string)(UserDto,error)  {
	var user model.SysUser
	var userDto UserDto
	var count int64
	common.GlobalDB.Table("sys_user").Where("username = ?", username).Scan(&user).Count(&count)
	if count <= 0 {
		return userDto, errors.New("暂无查询到数据")
	}
	u.StructToStruct(&userDto, &user)
	userDto.ID = user.ID
	return userDto, nil
}

// 用户登录
func (u *UserService) Login(username string, password string) error {
	var count int64
	dbCommon.Debug().Table("sys_user").Where("username = ? and password = ?",username, password).Count(&count)
	fmt.Println("查询出用户名和密码条数:", count)
	if count == 0 {
		return errors.New("用户名或密码错误")
	}
	return nil
}

// 获取用户权限列表,合并用户多角色权限和用户权限的权限,第一个参数为角色列表
func (u *UserService) GetUserPrivileges(roles []RoleDto) []PrivilegeDto {
	var privileges []PrivilegeDto
	privilegeIdKeyMap := make(map[int]bool,5)
	// 用户有多个角色,需要合并相同的权限
	for _,role := range roles {
		var userPrivileges []model.SysRolePrivilege
		common.GlobalDB.Table("sys_role_privilege").Where("role_id = ?",role.Id).Scan(&userPrivileges)
		for _,userPrivilege := range userPrivileges {
			if ok := privilegeIdKeyMap[userPrivilege.PrivilegeId]; !ok {
				privilegeIdKeyMap[userPrivilege.PrivilegeId] = true
				var privilege model.SysPrivilege
				var privilegeDto PrivilegeDto
				common.GlobalDB.Table("sys_privilege").Where("id = ?",userPrivilege.PrivilegeId).Scan(&privilege)
				privilegeDto.Id = privilege.ID
				privilegeDto.Code = privilege.Code
				privilegeDto.Type =  privilege.Type
				privilegeDto.TreeId = privilege.TreeId
				privilegeDto.TreeLevel = privilege.TreeLevel
				//privilegeDto.ParentId = privilege.Parent
				privilegeDto.Scope = privilege.Scope
				privilegeDto.Description = privilege.Description
				privilegeDto.Name  = privilege.Name
				privilegeDto.Leaf = privilege.Leaf
				privileges =  append(privileges, privilegeDto)
			}

		}
	}
	return privileges
}

// 通过用户id获取用户角色列表
func (u *UserService) GetUserRolesByUserId(userId uint) []RoleDto {
	var roles []RoleDto
	var userRoles []model.SysUserRole
	common.GlobalDB.Table("sys_user_role").Where("user_id = ?", userId).Scan(&userRoles)
	for _,userRole := range userRoles{
		var role RoleDto
		common.GlobalDB.Table("sys_role").Where("id = ?", userRole.RoleId).Scan(&role)
		roles = append(roles, role)
	}
	return roles
}

// 通过组织id获取用户组织
func (u *UserService) GetUserOrgById(id uint) OrgDto {
	var sysOrg model.SysOrg
	var orgDto OrgDto
	common.GlobalDB.Table("sys_org").Where("id = ?", id).Scan(&sysOrg)
	orgDto.Id = sysOrg.ID
	orgDto.Name = sysOrg.Name
	orgDto.ParentId = sysOrg.ParentId
	orgDto.TreeLevel = sysOrg.TreeLevel
	orgDto.TreeId = sysOrg.TreeId
	orgDto.Leaf = sysOrg.Leaf
	orgDto.Type = sysOrg.Type
	orgDto.Code = sysOrg.Code
	return orgDto
}

// 获取用户列表,获取其及其子组织下用户
func (u *UserService) GetUserList(requestDto UserRequestDto) ([]UserDto, int64, error) {
	var userDtoList []UserDto
	var countNum int64
	dbData := common.GlobalDB.Debug().Table("sys_user").Select("sys_user.*, sys_org.name as org_name").
		Joins("left join sys_org on sys_user.org_id = sys_org.id")
	if requestDto.Name != ""{
		dbData = dbData.Where("sys_user.name = ?", requestDto.Name)
	}
	if requestDto.Mobile != ""{
		dbData = dbData.Where("sys_user.mobile = ?", requestDto.Mobile)
	}
	if requestDto.OrgId != 0{
		dbData = dbData.Where("sys_user.org_id = ?", requestDto.OrgId)
	}
	dbData.Count(&countNum)
	dbData.Offset((requestDto.CurrPage-1)*requestDto.PageSize).Limit(requestDto.PageSize).Scan(&userDtoList)
	if len(userDtoList) == 0 {
		return userDtoList, countNum, errors.New("暂无数据")
	}
	return userDtoList, countNum, nil
}

// 保存用户
func (u *UserService) SaveUser(userData model.SysUser)  error {
	// 获取组织名称
	var orgName string
	common.GlobalDB.Debug().Table("sys_org").Select("name").Where("id=?", userData.OrgId).Scan(&orgName)
	if orgName == ""{
		return errors.New("添加失败")
	}
	// 说明是新增
	if userData.ID == 0 {
		// 判断用户名是否存在
		var usernameNum int64
		common.GlobalDB.Debug().Table("sys_user").Where("username=?", userData.Username).Count(&usernameNum)
		if usernameNum > 0{
			return errors.New("添加失败,用户名已经存在")
		}
	}
	_,err := u.Save(&userData)
	if err != nil {
		return err
	}
	return nil
}

// 重置密码
func (u *UserService) ResetPassword(userData *model.SysUser) error {
	rowsAffected := common.GlobalDB.Debug().Model(userData).Updates(userData).RowsAffected
	if rowsAffected == 0 {
		return errors.New("重置密码失败")
	}
	return nil
}

// 用户状态切换
func (u *UserService) Enabled(userData *model.SysUser) error {
	rowsAffected := common.GlobalDB.Debug().Model(userData).Update("enabled", userData.Enabled).RowsAffected
	if rowsAffected == 0 {
		return errors.New("重置密码失败")
	}
	return nil
}