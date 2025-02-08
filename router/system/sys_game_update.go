package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type GameUpdateRouter struct{}

func (g *GameUpdateRouter) InitGameUpdateRouter(Router *gin.RouterGroup) {
	router := Router.Group("job").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("job")

	routerApi := v1.ApiGroupApp.SystemApiGroup.GameUpdateApi

	{
		router.POST("gameUpdate/", routerApi.CreateGameUpdate)
		router.PUT("gameUpdate/", routerApi.UpdateGameUpdate)
		router.DELETE("gameUpdate/", routerApi.DeleteGameUpdate)
		router.POST("gameUpdate/exec/", routerApi.ExecUpdateTask)
	}

	{
		routerWithoutRecord.GET("gameUpdate/:id/", routerApi.GetGameUpdateById)
		routerWithoutRecord.GET("gameUpdate/", routerApi.GetGameUpdateList)
	}
}
