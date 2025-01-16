package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
)

type GameServerRouter struct {
}

func (*GameServerRouter) InitGameServerRouter(Router *gin.RouterGroup) {
	router := Router.Group("game")
	routerWithoutRecord := Router.Group("game")

	routerApi := v1.ApiGroupApp.SystemApiGroup.GameServerApi

	{
		router.POST("server/", routerApi.CreateGameServer)
		router.PUT("server/", routerApi.UpdateGameServer)
		router.DELETE("server/", routerApi.DeleteGameServer)
		router.GET("server/:id/", routerApi.GetGameServerById)
		router.POST("server/install/", routerApi.InstallGameServer)
	}
	{
		routerWithoutRecord.GET("server/", routerApi.GetGameServerList)
		routerWithoutRecord.GET("server/all/", routerApi.GetGameServerAll)
	}
}
