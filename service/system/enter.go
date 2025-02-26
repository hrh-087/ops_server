package system

type ServiceGroup struct {
	UserService
	JwtService
	CasbinService
	MenuService
	AuthorityService
	BaseMenuService
	ApiService
	ProjectService
	SshAuthService
	// 全局配置管理
	CloudProduceService
	GamePlatformService
	// 资产管理
	AssetsServerService
	AssetsMysqlService
	AssetsMongoService
	AssetsRedisService
	AssetsKafkaService
	AssetsLbService
	// 游戏管理
	GameTypeService
	GameServerService
	// 任务
	JobService
	JobTaskService
	JobCommandService
	// 更新
	GameUpdateService
	// 操作记录
	OperationRecordService

	SysTaskService

	GmService
}
