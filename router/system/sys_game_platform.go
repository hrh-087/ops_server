package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type GamePlatformRouter struct {
}

func (p *GamePlatformRouter) InitGamePlatformRouter(Router *gin.RouterGroup) {
	router := Router.Group("platform").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("platform")

	routerApi := v1.ApiGroupApp.SystemApiGroup.GamePlatformApi

	{
		router.GET("gamePlatform/:id/", routerApi.GetPlatformById)
		router.POST("gamePlatform/", routerApi.CreatePlatform)
		router.PUT("gamePlatform/", routerApi.UpdatePlatform)
		router.DELETE("gamePlatform/", routerApi.DeletePlatform)
		router.POST("gamePlatform/kick/", routerApi.KickGameServer)

	}
	{
		routerWithoutRecord.GET("gamePlatform/", routerApi.GetPlatformList)
		routerWithoutRecord.GET("gamePlatform/all/", routerApi.GetPlatformAll)
	}
}
