package request

type GameConfigFile struct {
	PlatformCode   string `json:"platformCode"`   // 渠道代号
	Vmid           int64  `json:"vmid"`           // 游戏服的vmid
	Name           string `json:"name"`           // 游戏服的名称
	PubIp          string `json:"pubIp"`          // 游戏服的公网ip
	InnerIp        string `json:"innerIp"`        // 游戏服的内网ip
	TcpPort        int64  `json:"tcpPort"`        // 游戏服的tcp端口
	HttpPort       int64  `json:"httpPort"`       // 游戏服的http端口
	GrpcPort       int64  `json:"grpcPort"`       // 游戏服的grpc端口
	MongoUri       string `json:"mongoUri"`       // mongo的uri
	MongoAuth      string `json:"mongoAuth"`      // mongo的认证
	DbName         string `json:"dbName"`         // mongo的db名
	KafkaUri       string `json:"kafkaUri"`       // kafka的uri
	RedisUri       string `json:"redisUri"`       // redis的uri
	RedisPass      string `json:"redisPass"`      // redis的密码
	RedisPort      int64  `json:"redisPort"`      // redis的端口
	RedisMeshUri   string `json:"redisMeshUri"`   // redis的mesh的uri
	RedisMeshPass  string `json:"redisMeshPass"`  // redis的mesh的密码
	RedisMeshPort  int64  `json:"meshPort"`       // redis的mesh的端口
	GatewayUri     string `json:"gatewayUri"`     // 网关的uri
	FightType      string `json:"fightType"`      // 战斗类型
	LtsGroupId     string `json:"ltsGroupId"`     // lts的groupid
	LtsStreamId    string `json:"ltsStreamId"`    // lts的streamid
	AccessKey      string `json:"accessKey"`      // lts的accesskey
	SecretKey      string `json:"secretKey"`      // lts的secretkey
	CloudProjectId string `json:"cloudProjectId"` // lts的云项目id
	CloudRegionId  string `json:"cloudRegionId"`  // lts的云区域id

	FilterUrl   string `json:"filterUrl"`   // 敏感词检测地址
	FilterToken string `json:"filterToken"` // 敏感词检测token
}

type DockerComposeFile struct {
	ImageService     string `json:"imageService"`
	ImageTag         string `json:"imageTag"`
	JsonConfigVolume string `json:"jsonConfigVolume"`
	ImageName        string `json:"imageName"`
	ImageUri         string `json:"imageUri"`
}
