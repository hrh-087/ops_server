package system

type RouterGroup struct {
	UserRouter
	BaseRouter
	AuthorityRouter
	MenuRouter
	ApiRouter
	CasbinRouter
	ProjectRouter
	SshAuthRouter
	OperationRecordRouter

	CloudProduceRouter
	GamePlatformRouter

	AssetsServerRouter
	AssetsMysqlRouter
	AssetsMongoRouter
	AssetsRedisRouter
	AssetsKafkaRouter
	AssetsLbRouter

	GameTypeRouter
	GameServerRouter

	JobRouter
	JobTaskRouter
	JobCommandRouter

	GameUpdateRouter

	SysTaskRouter

	GmRouter

	CronTaskRouter
}
