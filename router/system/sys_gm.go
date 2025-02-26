package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type GmRouter struct {
}

func (g *GmRouter) InitGmRouter(Router *gin.RouterGroup) {
	router := Router.Group("gm").Use(middleware.OperationRecord())

	routerApi := v1.ApiGroupApp.SystemApiGroup.GmApi

	{
		router.POST("setSwitch/", routerApi.SetSwitch)         // 设置开关
		router.POST("getSwitchList/", routerApi.GetSwitchList) // 获取开关列表
	}
}
