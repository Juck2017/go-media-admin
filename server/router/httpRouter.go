package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"net/http"
	"suitbim.com/go-media-admin/common"
	"suitbim.com/go-media-admin/core"
	_ "suitbim.com/go-media-admin/docs"
	"suitbim.com/go-media-admin/middleware"
	"suitbim.com/go-media-admin/model"
	"suitbim.com/go-media-admin/service"
	"suitbim.com/go-media-admin/utils"
)

func InitRouterGroup(dispatcher *core.Dispatcher) {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// 不需要认证的组
	noAuthGroup := r.Group("")
	{
		noAuthGroup.POST("/login", func(c *gin.Context) {
			var loginInfo struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			c.ShouldBindJSON(&loginInfo)
			var userService service.UserService
			err := userService.Login(loginInfo.Username, loginInfo.Password)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": err.Error(),
					"result":  []string{},
					"success": false,
				})
				return
			}
			token, _ := utils.GenToken(loginInfo.Username)
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": "登录成功",
				"result":  map[string]string{"token": token},
				"success": true,
			})
		})
		noAuthGroup.GET("/getChannelList", func(c *gin.Context) {
			// 根据客户端id获取通道列表(无需授权)
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
				Success bool		`json:"success"`
			}
			msg.Result = map[string]interface{}{}
			msg.Code = http.StatusOK
			clientId := c.Query("clientId")
			if clientId == "" {
				msg.Message = "客户端id不能为空"
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Message = "获取通道列表成功"
			var channelService service.ChannelService
			channelDtoList := channelService.GetChannelListByClientId(clientId)
			if len(channelDtoList) == 0 {
				msg.Result = []string{}
			}else {
				msg.Result = channelDtoList
			}
			msg.Success =  true
			c.JSON(http.StatusOK, msg)
		})
		noAuthGroup.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	}
	// token认证中间件
	jwtAuth := middleware.JWTAuth{}
	// 需要认证的组
	authGroup := r.Group("/rest", jwtAuth.AuthJwtToken())
	{
		authGroup.POST("/auth/getUserInfo", func(c *gin.Context) {
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			currentUser, exist := c.Get("username")
			if exist {
				msg.Code = http.StatusOK
				msg.Message = "获取用户信息成功"
				var userService service.UserService
				user, err := userService.GetUserInfoByUsername(currentUser.(string))
				if err != nil {
					msg.Message = err.Error()
				}
				roles := userService.GetUserRolesByUserId(user.ID)
				org := userService.GetUserOrgById(user.OrgId)
				privileges := userService.GetUserPrivileges(roles)
				msg.Result = map[string]interface{}{
					"org":        org,
					"privileges": privileges,
					"roles":      roles,
					"user":       user,
				}
				c.JSON(http.StatusOK, msg)
			}
		})
		authGroup.POST("/dic/getDicMap", func(c *gin.Context) {
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			var dicService service.DicService
			msg.Code = http.StatusOK
			msg.Message = "获取系统字典成功"
			dicMap, err := dicService.GetDicMap()
			msg.Result = dicMap
			if err != nil {
				msg.Message = err.Error()
				msg.Result = []string{}
				c.JSON(http.StatusOK, msg)
				return
			}
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/dic/list", func(c *gin.Context) {
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			var dicService service.DicService
			dicList, err := dicService.List()
			if err != nil {
				msg.Code = http.StatusCreated
				msg.Message = err.Error()
				msg.Result = []string{}
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Code = http.StatusOK
			msg.Message = "获取字典列表成功"
			msg.Result = dicList
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/dic/save", func(c *gin.Context) {
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			var dicService service.DicService
			dictModel := model.SysDict{}
			c.BindJSON(&dictModel)
			result, err := dicService.SaveTreeData(&dictModel)
			if err != nil {
				msg.Code = http.StatusCreated
				msg.Message = err.Error()
				msg.Result = []string{}
				c.JSON(http.StatusCreated, msg)
				return
			}
			dictData := result.(*model.SysDict)
			var dto service.DicDto
			dicService.StructToStruct(&dto, dictData)
			dto.Id = dictData.ID
			msg.Code = http.StatusOK
			msg.Message = "创建或更新数据成功"
			msg.Result = dto
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/dic/delete", func(c *gin.Context) {
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			dictData := model.SysDict{}
			c.BindJSON(&dictData)
			var dicService service.DicService
			err := dicService.DeleteTreeData(&dictData)
			msg.Result = []string{}
			if err != nil {
				msg.Code = http.StatusCreated
				msg.Message = err.Error()
				c.JSON(http.StatusCreated, msg)
				return
			}
			msg.Code = http.StatusOK
			msg.Message = "删除数据成功"
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/privilege/list", func(c *gin.Context) {
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			var privilegeService service.PrivilegeService
			privileges, err := privilegeService.GetPrivilegeList()
			if err != nil {
				msg.Code = http.StatusCreated
				msg.Message = err.Error()
				msg.Result = []string{}
				c.JSON(http.StatusCreated, msg)
				return
			}
			msg.Code = http.StatusOK
			msg.Message = "获取权限列表成功"
			msg.Result = privileges
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/privilege/save", func(c *gin.Context) {
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			privilegeData := model.SysPrivilege{}
			c.BindJSON(&privilegeData)
			var privilegeService service.PrivilegeService
			result, err := privilegeService.SaveTreeData(&privilegeData)
			if err != nil {
				msg.Code = http.StatusCreated
				msg.Message = err.Error()
				msg.Result = []string{}
				c.JSON(http.StatusCreated, msg)
				return
			}
			privilege := result.(*model.SysPrivilege)
			var dto service.PrivilegeDto
			privilegeService.StructToStruct(&dto, privilege)
			dto.Id = privilege.ID
			msg.Code = http.StatusOK
			msg.Message = "创建或更新数据成功"
			msg.Result = dto
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/privilege/delete", func(c *gin.Context) {
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			privilegeData := model.SysPrivilege{}
			c.BindJSON(&privilegeData)
			var privilegeService service.PrivilegeService
			err := privilegeService.DeleteTreeData(&privilegeData)
			msg.Result = []string{}
			if err != nil {
				msg.Code = http.StatusCreated
				msg.Message = err.Error()
				c.JSON(http.StatusCreated, msg)
				return
			}
			msg.Code = http.StatusOK
			msg.Message = "删除数据成功"
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/org/list", func(c *gin.Context) {
			// 获取组织列表
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			var orgService service.OrgService
			orgList, err := orgService.GetOrgList()
			if err != nil {
				msg.Code = http.StatusCreated
				msg.Message = err.Error()
				msg.Result = []string{}
				c.JSON(http.StatusCreated, msg)
				return
			}
			msg.Code = http.StatusOK
			msg.Message = "获取组织列表成功"
			msg.Result = orgList
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/org/rootList", func(c *gin.Context) {
			// 获取根组织列表
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			var orgService service.OrgService
			orgRootList, err := orgService.GetRootOrgList()
			if err != nil {
				msg.Code = http.StatusOK
				msg.Message = err.Error()
				msg.Result = []string{}
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Code = http.StatusOK
			msg.Message = "获取根组织列表成功"
			msg.Result = orgRootList
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/org/roles", func(c *gin.Context) {
			// 获取组织角色列表
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
				Success bool `json:"success"`
			}
			orgData := model.SysOrg{}
			c.BindJSON(&orgData)
			var orgRoleService service.OrgRoleService
			roleList, err := orgRoleService.GetOrgRoleList(uint64(orgData.ID))
			if err != nil {
				msg.Code = http.StatusOK
				msg.Message = err.Error()
				msg.Result = []string{}
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Code = http.StatusOK
			msg.Message = "获取组织角色列表成功"
			msg.Result = roleList
			msg.Success = true
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/org/save", func(c *gin.Context) {
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			orgData := model.SysOrg{}
			c.BindJSON(&orgData)
			var orgService service.OrgService
			result, err := orgService.SaveTreeData(&orgData)
			if err != nil {
				msg.Code = http.StatusCreated
				msg.Message = err.Error()
				msg.Result = []string{}
				c.JSON(http.StatusCreated, msg)
				return
			}
			sysOrg := result.(*model.SysOrg)
			var orgDto service.OrgDto
			orgService.StructToStruct(&orgDto, sysOrg)
			orgDto.Id = sysOrg.ID
			msg.Code = http.StatusOK
			msg.Message = "创建或更新数据成功"
			msg.Result = orgDto
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/org/delete", func(c *gin.Context) {
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			orgData := model.SysOrg{}
			c.BindJSON(&orgData)
			var orgService service.OrgService
			err := orgService.DeleteTreeData(&orgData)
			msg.Result = []string{}
			if err != nil {
				msg.Code = http.StatusCreated
				msg.Message = err.Error()
				c.JSON(http.StatusCreated, msg)
				return
			}
			msg.Code = http.StatusOK
			msg.Message = "删除数据成功"
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/org/privileges", func(c *gin.Context) {
			// 获取组织对应权限列表
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
				Success bool `json:"success"`
			}
			orgData := model.SysOrg{}
			c.BindJSON(&orgData)
			msg.Message = "获取组织权限列表成功"
			msg.Code = http.StatusOK
			var orgPrivilegeService service.OrgPrivilegeService
			orgPrivilegeList,err := orgPrivilegeService.GetOrgPrivilegeList(uint64(orgData.ID))
			if err != nil {
				msg.Result = []string{}
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Result = orgPrivilegeList
			msg.Success = true
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/org/configPrivileges", func(c *gin.Context) {
			// 配置权限
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			msg.Result = []string{}
			var orgPrivilegeDto service.ConfigPrivilegeDto
			if err := c.ShouldBind(&orgPrivilegeDto); err != nil {
				msg.Code = http.StatusBadRequest
				msg.Message = "json解析失败"
				c.JSON(http.StatusCreated, msg)
				return
			}
			var orgPrivilegeService service.OrgPrivilegeService
			err := orgPrivilegeService.ConfigPrivileges(orgPrivilegeDto)
			if err != nil {
				msg.Code = http.StatusCreated
				msg.Message = err.Error()
				c.JSON(http.StatusCreated, msg)
				return
			}
			msg.Code = http.StatusOK
			msg.Message = "更新权限成功"
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/role/list", func(c *gin.Context) {
			// 获取角色列表
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			msg.Result = []string{}
			msg.Code = http.StatusOK
			// 构造请求数据体
			var requestData service.RoleRequestDto
			err := c.ShouldBindJSON(&requestData)
			if err != nil {
				msg.Message = "请求数据解析失败"
				c.JSON(http.StatusOK, msg)
				return
			}
			var roleService service.RoleService
			roleList,total, err := roleService.GetRoleList(requestData)
			if err != nil {
				msg.Message = err.Error()
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Message = "获取角色列表成功"
			msg.Result = map[string]interface{}{
				"list": roleList,
				"currPage": requestData.CurrPage,
				"pageSize": requestData.PageSize,
				"total": total,
			}
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/role/privileges", func(c *gin.Context) {
			// 获取角色权限列表
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			msg.Result = []string{}
			msg.Code = http.StatusOK
			// 构造请求数据体
			var rolePrivilege model.SysRolePrivilege
			err := c.ShouldBindJSON(&rolePrivilege)
			if err != nil {
				msg.Message = "请求数据解析失败"
				c.JSON(http.StatusOK, msg)
				return
			}
			var rolePrivilegeService service.RolePrivilegeService
			rolePrivilegeList, err := rolePrivilegeService.GetRolePrivilegeList(rolePrivilege.RoleId)
			if err != nil {
				msg.Message = err.Error()
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Message = "获取角色权限成功"
			msg.Result = rolePrivilegeList
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/role/configPrivileges", func(c *gin.Context) {
			// 配置角色权限
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			msg.Result = []string{}
			msg.Code = http.StatusOK
			var rolePrivilegeDto service.ConfigRolePrivilegeDto
			if err := c.ShouldBind(&rolePrivilegeDto); err != nil {
				msg.Message = "json解析失败"
				c.JSON(http.StatusOK, msg)
				return
			}
			var rolePrivilegeService service.RolePrivilegeService
			err := rolePrivilegeService.ConfigRolePrivileges(rolePrivilegeDto)
			if err != nil {
				msg.Message = err.Error()
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Code = http.StatusOK
			msg.Message = "更新权限成功"
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/role/save", func(c *gin.Context) {
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			roleData := model.SysRole{}
			c.BindJSON(&roleData)
			var roleService service.RoleService
			result, err := roleService.SaveRole(roleData)
			if err != nil {
				msg.Code = http.StatusOK
				msg.Message = err.Error()
				msg.Result = []string{}
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Code = http.StatusOK
			msg.Message = "创建或更新数据成功"
			msg.Result = result
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/role/delete", func(c *gin.Context) {
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			roleData := model.SysRole{}
			c.BindJSON(&roleData)
			var roleService service.RoleService
			err := roleService.Delete(roleData.ID)
			msg.Result = []string{}
			if err != nil {
				msg.Code = http.StatusCreated
				msg.Message = err.Error()
				c.JSON(http.StatusCreated, msg)
				return
			}
			msg.Code = http.StatusOK
			msg.Message = "删除数据成功"
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/user/list", func(c *gin.Context) {
			// 获取用户列表
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
				Success bool		`json:"success"`
			}
			msg.Result = map[string]interface{}{}
			msg.Code = http.StatusOK
			// 构造请求数据体
			var requestData service.UserRequestDto
			err := c.ShouldBindJSON(&requestData)
			if err != nil {
				msg.Message = "请求数据解析失败"
				c.JSON(http.StatusOK, msg)
				return
			}
			var userService service.UserService
			userDtoList,total, err := userService.GetUserList(requestData)
			if err != nil {
				msg.Message = err.Error()
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Message = "获取用户列表成功"
			msg.Result = map[string]interface{}{
				"list": userDtoList,
				"currPage": requestData.CurrPage,
				"pageSize": requestData.PageSize,
				"total": total,
			}
			msg.Success =  true
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/user/roles", func(c *gin.Context) {
			// 获取用户角色列表
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
				Success bool		`json:"success"`
			}
			msg.Result = []string{}
			msg.Code = http.StatusOK
			// 构造请求数据体
			var user model.SysUser
			err := c.ShouldBindJSON(&user)
			if err != nil {
				msg.Message = "请求数据解析失败"
				c.JSON(http.StatusOK, msg)
				return
			}
			var userRoleService service.UserRoleService
			userRoleList, err := userRoleService.GetUserRoleList(user.ID)
			if err != nil {
				msg.Message = err.Error()
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Result = userRoleList
			if len(userRoleList) == 0 {
				msg.Result = []string{}
			}
			msg.Message = "获取用户角色列表成功"
			msg.Success =  true
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/user/privileges", func(c *gin.Context) {
			// 获取用户对应权限列表
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
				Success bool `json:"success"`
			}
			userData := model.SysUser{}
			c.BindJSON(&userData)
			msg.Message = "获取用户权限列表成功"
			msg.Code = http.StatusOK
			var userPrivilegeService service.UserPrivilegeService
			userPrivilegeList, err := userPrivilegeService.GetUserPrivilegeList(uint64(userData.ID))
			if err != nil {
				msg.Result = []string{}
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Result = userPrivilegeList
			msg.Success = true
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/user/configRoles", func(c *gin.Context) {
			// 配置用角色
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
				Success bool		`json:"success"`
			}
			msg.Result = []string{}
			msg.Code = http.StatusOK
			// 构造请求数据体
			var configRole service.ConfigRoleDto
			err := c.ShouldBindJSON(&configRole)
			if err != nil {
				msg.Message = "请求数据解析失败"
				c.JSON(http.StatusOK, msg)
				return
			}
			var userRoleService service.UserRoleService
			err = userRoleService.ConfigRoles(configRole)
			if err != nil {
				msg.Message = err.Error()
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Result = []string{}
			msg.Message = "配置用户角色成功"
			msg.Success =  true
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/user/configPrivileges", func(c *gin.Context) {
			// 配置用户权限
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
				Success bool		`json:"success"`
			}
			msg.Result = []string{}
			msg.Code = http.StatusOK
			// 构造请求数据体
			var configUProDto service.UserPrivilegeDto
			err := c.ShouldBindJSON(&configUProDto)
			if err != nil {
				msg.Message = "请求数据解析失败"
				c.JSON(http.StatusOK, msg)
				return
			}
			var userPrivilegeService service.UserPrivilegeService
			err = userPrivilegeService.ConfigUserPrivileges(configUProDto)
			if err != nil {
				msg.Message = err.Error()
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Result = []string{}
			msg.Message = "配置用户权限成功"
			msg.Success =  true
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/user/resetPassword", func(c *gin.Context) {
			// 重置密码
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
				Success bool `json:"success"`
			}
			userData := model.SysUser{}
			c.BindJSON(&userData)
			msg.Message = "重置密码成功"
			msg.Code = http.StatusOK
			msg.Result = []string{}
			var userService service.UserService
			err := userService.ResetPassword(&userData)
			if err != nil {
				msg.Message = "重置密码失败"
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Success = true
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/user/enabled", func(c *gin.Context) {
			// 用户状态切换
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
				Success bool `json:"success"`
			}
			userData := model.SysUser{}
			c.BindJSON(&userData)
			msg.Message = "切换用户状态成功"
			msg.Code = http.StatusOK
			msg.Result = []string{}
			var userService service.UserService
			err := userService.Enabled(&userData)
			if err != nil {
				msg.Message = "切换用户状态失败"
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Success = true
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/user/save", func(c *gin.Context) {
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
			}
			userData := model.SysUser{}
			c.BindJSON(&userData)
			var userService service.UserService
			err := userService.SaveUser(userData)
			msg.Result = []string{}
			if err != nil {
				msg.Code = http.StatusCreated
				msg.Message = err.Error()
				c.JSON(http.StatusOK, msg)
				return
			}
			msg.Code = http.StatusOK
			msg.Message = "创建或更新数据成功"
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/channel/list", func(c *gin.Context) {
			// 获取通道列表
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
				Success bool		`json:"success"`
			}
			msg.Result = map[string]interface{}{}
			msg.Code = http.StatusOK
			// 构造请求数据体
			var requestData service.ChannelRequestDto
			err := c.ShouldBindJSON(&requestData)
			if err != nil {
				msg.Message = "请求数据解析失败"
				c.JSON(http.StatusOK, msg)
				return
			}
			cameraService := &service.ChannelService{Dispatch: dispatcher}
			cameraDtoList,total := cameraService.GetChannelList(requestData)
			msg.Message = "获取通道列表成功"
			if len(cameraDtoList) == 0 {
				msg.Result = map[string]interface{}{
					"list": []string{},
					"currPage": requestData.CurrPage,
					"pageSize": requestData.PageSize,
					"total": total,
				}
			}else {
				msg.Result = map[string]interface{}{
					"list": cameraDtoList,
					"currPage": requestData.CurrPage,
					"pageSize": requestData.PageSize,
					"total": total,
				}
			}
			msg.Success =  true
			c.JSON(http.StatusOK, msg)
		})
		authGroup.POST("/pusher/list", func(c *gin.Context) {
			// 获取推流客户端列表
			var msg struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Result  interface{} `json:"result"`
				Success bool		`json:"success"`
			}
			msg.Result = map[string]interface{}{}
			msg.Code = http.StatusOK
			// 构造请求数据体
			var requestData service.PusherRequestDto
			err := c.ShouldBindJSON(&requestData)
			if err != nil {
				msg.Message = "请求数据解析失败"
				c.JSON(http.StatusOK, msg)
				return
			}
			var pusherService service.PusherService
			pusherDtoList,total, err := pusherService.GetPusherList(requestData)
			msg.Message = "获取客户端列表成功"
			if len(pusherDtoList) == 0 {
				msg.Result = map[string]interface{}{
					"list": []string{},
					"currPage": requestData.CurrPage,
					"pageSize": requestData.PageSize,
					"total": total,
				}
			} else {
				msg.Result = map[string]interface{}{
					"list": pusherDtoList,
					"currPage": requestData.CurrPage,
					"pageSize": requestData.PageSize,
					"total": total,
				}
			}
			msg.Success =  true
			c.JSON(http.StatusOK, msg)
		})
	}
	r.Run(fmt.Sprintf("0.0.0.0:%d", common.GlobalConf.App.Port))
}
