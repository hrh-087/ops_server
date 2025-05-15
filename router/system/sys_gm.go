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

		router.POST("getActivityExtra/", routerApi.GetActivityExtra) // 获取活动额外配置
		router.POST("setGameCron/", routerApi.SetGameCron)           // 设置游戏定时任务
		router.POST("getGameCronList/", routerApi.GetGameCronList)   // 获取游戏定时任务列表

		router.POST("kickLivePlayer/", routerApi.KickLivePlayer)       // 踢出直播服玩家
		router.POST("getLiveServerList/", routerApi.GetLiveServerList) // 获取直播服列表
		router.POST("clearLiveServer/", routerApi.ClearLiveServer)     // 清除直播服数据
		router.POST("setLiveServer/", routerApi.SetLiveActivityServer) // 设置直播服数据
	}
}
