package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type GameTypeRouter struct {
}

func (*GameTypeRouter) InitGameTypeRouter(Router *gin.RouterGroup) {
	router := Router.Group("game").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("game")

	routerApi := v1.ApiGroupApp.SystemApiGroup.GameTypeApi

	{
		router.POST("type/", routerApi.CreateGameType)
		router.PUT("type/", routerApi.UpdateGameType)
		router.DELETE("type/", routerApi.DeleteGameType)
		router.GET("type/:id/", routerApi.GetGameTypeById)
		router.POST("type/copyAll/", routerApi.CopyGameType)
	}
	{
		routerWithoutRecord.GET("type/", routerApi.GetGameTypeList)
		routerWithoutRecord.GET("type/all/", routerApi.GetGameTypeAll)
	}
}
