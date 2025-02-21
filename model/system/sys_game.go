package system

type SysGameType struct {
	ProjectModel
	Name            string `json:"name" gorm:"comment:游戏类型名称"`
	Code            string `json:"code" gorm:"index;comment:游戏类型代码"`
	TcpPort         string `json:"tcpPort" gorm:"comment:tcp端口范围"`
	HttpPort        string `json:"httpPort" gorm:"comment:http端口范围"`
	GrpcPort        string `json:"grpcPort" gorm:"comment:grpc端口范围"`
	VmidRule        int64  `json:"vmidRule" gorm:"index;comment:vmid规则"`
	ConfigTemplate  string `json:"configTemplate" gorm:"type:text;comment:配置模板"`
	ComposeTemplate string `json:"composeTemplate" gorm:"type:text;comment:docker-compose模板"`
	//ServiceType     string `json:"serviceType" gorm:"index;comment:服务类型: rim服务/游戏服务"`
	Order   int64 `json:"-" gorm:"comment:排序"`
	IsFight bool  `json:"isFight" gorm:"comment:是否是战斗服"`
}

func (*SysGameType) TableName() string {
	return "sys_game_type"
}

type SysGameServer struct {
	ProjectModel
	PlatformId uint            `from:"platformId" json:"platformId" gorm:"comment:渠道id"`
	Platform   SysGamePlatform `json:"platform" gorm:"foreignKey:PlatformId;references:ID"`

	Name        string `form:"name" json:"name" gorm:"comment:名称"`
	Vmid        int64  `json:"vmid" gorm:"comment:vmid"`
	TcpPort     int64  `json:"tcpPort" gorm:"comment:tcp端口"`
	HttpPort    int64  `json:"httpPort" gorm:"comment:http端口"`
	GrpcPort    int64  `json:"grpcPort" gorm:"comment:grpc端口"`
	Status      int64  `form:"status" json:"status" gorm:"default:5;comment:状态 0: 待安装 1:安装中 2:已安装 3:已删除 4 安装失败 5待安装"`
	ConfigFile  string `json:"configFile" gorm:"type:text;comment:配置文件"`
	ComposeFile string `json:"composeFile" gorm:"type:text;comment:docker-compose文件"`

	GameTypeId uint        `form:"gameTypeId" json:"gameTypeId" gorm:"comment:游戏类型id"`
	GameType   SysGameType `json:"gameType" gorm:"foreignKey:GameTypeId;references:ID"`

	RedisId uint           `json:"redisId" gorm:"comment:redis_id"`
	Redis   SysAssetsRedis `json:"-" gorm:"foreignKey:RedisId;references:ID"`

	MongoId uint           `json:"mongoId" gorm:"comment:mongo_id"`
	Mongo   SysAssetsMongo `json:"-" gorm:"foreignKey:MongoId;references:ID"`

	KafkaId uint           `json:"kafkaId" gorm:"comment:kafka_id"`
	Kafka   SysAssetsKafka `json:"-" gorm:"foreignKey:KafkaId;references:ID"`

	HostId uint            `json:"hostId" gorm:"comment:主机id"`
	Host   SysAssetsServer `json:"host" gorm:"foreignKey:HostId;references:ID"`
}
