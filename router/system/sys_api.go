package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type ApiRouter struct{}

func (s *ApiRouter) InitApiRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {
	apiRouter := Router.Group("api").Use(middleware.OperationRecord())
	//apiRouter := Router.Group("api")
	apiRouterWithoutRecord := Router.Group("api")

	apiPublicRouterWithoutRecord := RouterPub.Group("api")
	apiRouterApi := v1.ApiGroupApp.SystemApiGroup.SystemApiApi
	{
		apiRouter.GET("getApiGroups/", apiRouterApi.GetApiGroups)          // 获取路由组
		apiRouter.GET("syncApi/", apiRouterApi.SyncApi)                    // 同步Api
		apiRouter.POST("ignoreApi/", apiRouterApi.IgnoreApi)               // 忽略Api
		apiRouter.POST("enterSyncApi/", apiRouterApi.EnterSyncApi)         // 确认同步Api
		apiRouter.POST("createApi/", apiRouterApi.CreateApi)               // 创建Api
		apiRouter.POST("deleteApi/", apiRouterApi.DeleteApi)               // 删除Api
		apiRouter.POST("getApiById/", apiRouterApi.GetApiById)             // 获取单条Api消息
		apiRouter.POST("updateApi/", apiRouterApi.UpdateApi)               // 更新api
		apiRouter.DELETE("deleteApisByIds/", apiRouterApi.DeleteApisByIds) // 删除选中api
	}
	{
		apiRouterWithoutRecord.POST("getAllApis/", apiRouterApi.GetAllApis) // 获取所有api
		apiRouterWithoutRecord.POST("getApiList/", apiRouterApi.GetApiList) // 获取Api列表
	}
	{
		apiPublicRouterWithoutRecord.GET("freshCasbin/", apiRouterApi.FreshCasbin) // 刷新casbin权限
	}
}
