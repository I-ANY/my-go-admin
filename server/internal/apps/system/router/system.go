package router

import (
	"biz-auto-api/internal/apps/system/apis"
	"biz-auto-api/internal/common"
	"github.com/gin-gonic/gin"
)

func SetSystemRouter(v1 *gin.RouterGroup) {
	systemGroup := v1.Group("/system")
	{
		systemGroup.GET(common.HealthzApi, common.GetHealthzFunc())
	}
	baseApi := apis.SysBase{}
	{
		systemGroup.POST("/login", baseApi.UserLogin)
		systemGroup.POST("/refreshToken", baseApi.RefreshToken)
		systemGroup.POST("/logout", baseApi.UserLogout)
		systemGroup.GET("/userInfo", baseApi.GetUserInfo)
		systemGroup.GET("/userPermCode", baseApi.GetUserPermCode)
		systemGroup.GET("/userMenu", baseApi.GetUserMenu)
		systemGroup.GET("/star-portal/login-url", baseApi.GetStarPortalLoginUrl)
		systemGroup.POST("/star-portal/login", baseApi.StarPortalLogin)

	}

	menuApi := apis.SysMenu{}
	{
		systemGroup.GET("/menu/tree", menuApi.GetMenuTree)
		systemGroup.PUT("/menu/:id", menuApi.UpdateMenu)
		systemGroup.POST("/menu", menuApi.AddMenu)
		systemGroup.DELETE("/menu/:id", menuApi.DeleteMenu)

	}

	apiApi := apis.SysApi{}
	{
		systemGroup.GET("/api/list", apiApi.GetApiList)
	}

	roleApi := apis.SysRole{}
	{
		systemGroup.GET("/role/list", roleApi.GetRoleList)
		systemGroup.PUT("/role/:id", roleApi.UpdateRole)
		systemGroup.POST("/role", roleApi.AddRole)
		systemGroup.DELETE("/role/:id", roleApi.DeleteRole)
	}

	deptApi := apis.SysDept{}
	{
		systemGroup.GET("/dept/tree", deptApi.GetDeptTree)
	}

	userApi := apis.SysUser{}
	{
		systemGroup.GET("/user/list", userApi.GetUserList)
		systemGroup.PUT("/user/:id", userApi.UpdateUser)
		systemGroup.POST("/user", userApi.AddUser)
		systemGroup.DELETE("/user/:id", userApi.DeleteUser)
	}
	operaLog := apis.SysOperaLog{}
	{
		systemGroup.GET("/operaLog/list", operaLog.GetOperaLogList)
	}

	dictAPi := apis.SysDict{}
	{
		systemGroup.GET("/dict/all", dictAPi.GetAllDictData)

		//systemGroup.POST("/dict/type")
		//systemGroup.DELETE("/dict/type/:id")
		//systemGroup.PUT("/dict/type/:id")
		//systemGroup.GET("/dict/type/list")
		//
		//systemGroup.POST("/dict/data")
		//systemGroup.DELETE("/dict/data/:id")
		//systemGroup.PUT("/dict/data/:id")
		//systemGroup.GET("/dict/:id/data/list")
	}
	resourceApi := apis.SysResource{}
	{
		systemGroup.GET("/resource/list", resourceApi.GetResourceList)
		systemGroup.GET("/resource/table/list", resourceApi.GetResourceTableList)
		systemGroup.GET("/resource/table/field", resourceApi.GetResourceTableField)
		systemGroup.POST("/resource", resourceApi.AddResource)
		systemGroup.DELETE("/resource/:id", resourceApi.DeleteResource)
		systemGroup.PUT("/resource/:id", resourceApi.UpdateResource)
		systemGroup.GET("/resource/view/search-form-schemas", resourceApi.GetResourceViewFormSchemas)
		systemGroup.GET("/resource/view/table-columns", resourceApi.GetResourceViewTableColumns)
		systemGroup.GET("/resource/detail/list", resourceApi.GetResourceDetailList)
		systemGroup.GET("/role/resource/info", resourceApi.GetRoleResourceInfo)
		systemGroup.GET("/role/:roleId/resource/detail/list", resourceApi.GetRoleResourceDetailList)
		systemGroup.POST("/role/:roleId/resource", resourceApi.UpdateRoleResource)
		systemGroup.GET("/role/authed/resource", resourceApi.GetRoleAuthedResource)
		systemGroup.GET("/resource/subcategory", resourceApi.GetBusinessResource)
		systemGroup.POST("/role/:roleId/resource/auth", resourceApi.RoleResourceAuth)
	}
}
