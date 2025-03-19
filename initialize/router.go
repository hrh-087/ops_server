package initialize

import (
	"github.com/gin-gonic/gin"
	"ops-server/global"
	"ops-server/middleware"
	"ops-server/router"
)

func Routers() *gin.Engine {
	Router := gin.New()
	Router.Use(gin.Recovery())

	if gin.Mode() == gin.DebugMode {
		// 调试模式则使用gin自带的日志记录器
		Router.Use(gin.Logger())
	}

	systemRouter := router.RouterGroupApp.System

	PublicGroup := Router.Group(global.OPS_CONFIG.System.RouterPrefix)

	PrivateGroup := Router.Group(global.OPS_CONFIG.System.RouterPrefix)

	ProjectGroup := Router.Group(global.OPS_CONFIG.System.RouterPrefix)

	PrivateGroup.Use(middleware.JwtAuth()).Use(middleware.CasbinHandler())
	//PrivateGroup.Use(middleware.JwtAuth())

	ProjectGroup.Use(middleware.JwtAuth()).Use(middleware.ProjectAuth()).Use(middleware.CasbinHandler())

	{
		systemRouter.BaseRouter.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
	}

	{
		systemRouter.UserRouter.InitUserRouter(PrivateGroup)
		systemRouter.AuthorityRouter.InitAuthorityRouter(PrivateGroup)
		systemRouter.MenuRouter.InitMenuRouter(PrivateGroup)
		systemRouter.ApiRouter.InitApiRouter(PrivateGroup, PublicGroup)
		systemRouter.CasbinRouter.InitCasbinRouter(PrivateGroup)
		systemRouter.ProjectRouter.InitProjectRouter(PrivateGroup)
		systemRouter.OperationRecordRouter.InitSysOperationRecordRouter(ProjectGroup)

		systemRouter.SshAuthRouter.InitSshAuthRouter(ProjectGroup)
		systemRouter.CloudProduceRouter.InitCloudProduceRouter(ProjectGroup)
		systemRouter.GamePlatformRouter.InitGamePlatformRouter(ProjectGroup)

		systemRouter.AssetsServerRouter.InitAssetsServerRouter(ProjectGroup)
		systemRouter.AssetsMysqlRouter.InitAssetsMysqlRouter(ProjectGroup)
		systemRouter.AssetsMongoRouter.InitAssetsMongoRouter(ProjectGroup)
		systemRouter.AssetsRedisRouter.InitAssetsRedisRouter(ProjectGroup)
		systemRouter.AssetsKafkaRouter.InitAssetsKafkaRouter(ProjectGroup)
		systemRouter.AssetsLbRouter.InitAssetsLbRouter(ProjectGroup)

		systemRouter.GameTypeRouter.InitGameTypeRouter(ProjectGroup)
		systemRouter.GameServerRouter.InitGameServerRouter(ProjectGroup)

		systemRouter.JobRouter.InitJobRouter(ProjectGroup)
		systemRouter.JobTaskRouter.InitJobTaskRouter(PrivateGroup)
		systemRouter.JobCommandRouter.InitJobCommandRouter(ProjectGroup)

		systemRouter.GameUpdateRouter.InitGameUpdateRouter(ProjectGroup)
		systemRouter.SysTaskRouter.InitSysTaskRouter(ProjectGroup)

		systemRouter.GmRouter.InitGmRouter(ProjectGroup)

		systemRouter.CronTaskRouter.InitCronTaskRouter(ProjectGroup)

	}

	global.OPS_ROUTERS = Router.Routes()

	return Router
}
