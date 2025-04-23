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
		router.POST("getRankList/", routerApi.GetRankList)
		router.POST("getRankRewardList/", routerApi.GetRankRewardList)
		router.POST("setRankReward/", routerApi.SetRankConfig)
		router.POST("uploadGameConfig/", routerApi.UploadGameConfig)
		router.POST("getItemConfigInfo/", routerApi.GetItemConfigInfo)
		router.POST("dimensionPush/", routerApi.DimensionPush)       // 维度推送
		router.POST("getAccountInfo/", routerApi.GetPlayerId)        // 获取玩家账号信息
		router.POST("importPlayerData/", routerApi.ImportPlayerData) // 导入玩家数据
	}
}
