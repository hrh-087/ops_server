package system

import "ops-server/service"

type ApiGroup struct {
	BaseApi
	AuthorityApi
	AuthorityMenuApi
	SystemApiApi
	CasbinApi
	ProjectApi
	SshAuthApi
	CloudProduceApi
	GamePlatformApi

	// 资产管理
	AssetsServerApi
	AssetsMysqlApi
	AssetsRedisApi
	AssetsMongoApi
	AssetsKafkaApi
	AssetsLbApi

	// 游戏服管理
	GameTypeApi
	GameServerApi

	// 作业管理
	JobApi
	JobTaskApi
	JobCommandApi

	// 更新
	GameUpdateApi
	// 操作记录
	OperationRecordApi

	SysTaskApi
	GmApi
}

var (
	userService      = service.ServiceGroupApp.SystemServiceGroup.UserService
	authorityService = service.ServiceGroupApp.SystemServiceGroup.AuthorityService
	casbinService    = service.ServiceGroupApp.SystemServiceGroup.CasbinService
	menuService      = service.ServiceGroupApp.SystemServiceGroup.MenuService
	baseMenuService  = service.ServiceGroupApp.SystemServiceGroup.BaseMenuService
	apiService       = service.ServiceGroupApp.SystemServiceGroup.ApiService
	projectService   = service.ServiceGroupApp.SystemServiceGroup.ProjectService
	sshAuthService   = service.ServiceGroupApp.SystemServiceGroup.SshAuthService

	cloudProduceService = service.ServiceGroupApp.SystemServiceGroup.CloudProduceService
	platformService     = service.ServiceGroupApp.SystemServiceGroup.GamePlatformService

	assetsServerService = service.ServiceGroupApp.SystemServiceGroup.AssetsServerService
	assetsMysqlService  = service.ServiceGroupApp.SystemServiceGroup.AssetsMysqlService
	assetsRedisService  = service.ServiceGroupApp.SystemServiceGroup.AssetsRedisService
	assetsMongoService  = service.ServiceGroupApp.SystemServiceGroup.AssetsMongoService
	assetsKafkaService  = service.ServiceGroupApp.SystemServiceGroup.AssetsKafkaService
	assetsLbService     = service.ServiceGroupApp.SystemServiceGroup.AssetsLbService

	gameTypeService   = service.ServiceGroupApp.SystemServiceGroup.GameTypeService
	gameServerService = service.ServiceGroupApp.SystemServiceGroup.GameServerService

	jobService        = service.ServiceGroupApp.SystemServiceGroup.JobService
	jobTaskService    = service.ServiceGroupApp.SystemServiceGroup.JobTaskService
	jobCommandService = service.ServiceGroupApp.SystemServiceGroup.JobCommandService

	gameUpdateService      = service.ServiceGroupApp.SystemServiceGroup.GameUpdateService
	operationRecordService = service.ServiceGroupApp.SystemServiceGroup.OperationRecordService

	sysTaskService = service.ServiceGroupApp.SystemServiceGroup.SysTaskService
	gmService      = service.ServiceGroupApp.SystemServiceGroup.GmService
)
