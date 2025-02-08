package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type AssetsRedisRouter struct {
}

func (i *AssetsRedisRouter) InitAssetsRedisRouter(Router *gin.RouterGroup) {
	router := Router.Group("assets").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("assets")

	routerApi := v1.ApiGroupApp.SystemApiGroup.AssetsRedisApi

	{
		router.GET("redis/:id/", routerApi.GetAssetsRedisById)
		router.POST("redis/", routerApi.CreateRedis)
		router.PUT("redis/", routerApi.UpdateRedis)
		router.DELETE("redis/", routerApi.DeleteRedis)
	}
	{
		routerWithoutRecord.GET("redis/", routerApi.GetAssetsRedisList)
		routerWithoutRecord.GET("redis/all/", routerApi.GetAssetsRedisAll)
	}
}
