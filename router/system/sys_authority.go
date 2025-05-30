package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type AuthorityRouter struct {
}

// InitAuthorityRouter 初始化角色路由
func (u *AuthorityRouter) InitAuthorityRouter(Router *gin.RouterGroup) {
	authorityRouter := Router.Group("authority").Use(middleware.OperationRecord())
	//authorityRouter := Router.Group("authority")
	authorityRouterWithoutRecord := Router.Group("authority")
	authorityApi := v1.ApiGroupApp.SystemApiGroup.AuthorityApi
	{
		authorityRouter.POST("createAuthority/", authorityApi.CreateAuthority) // 创建角色
		authorityRouter.POST("deleteAuthority/", authorityApi.DeleteAuthority) // 删除角色
		authorityRouter.PUT("updateAuthority/", authorityApi.UpdateAuthority)  // 更新角色
		authorityRouter.POST("copyAuthority/", authorityApi.CopyAuthority)     // 拷贝角色
		//authorityRouter.POST("setDataAuthority", authorityApi.SetDataAuthority) // 设置角色资源权限
	}
	{
		authorityRouterWithoutRecord.POST("getAuthorityList/", authorityApi.GetAuthorityList) // 获取角色列表
	}
}
