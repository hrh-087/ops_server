package request

type GameConfigFile struct {
	PlatformCode  string `json:"platformCode"`
	Vmid          int64  `json:"vmid"`
	Name          string `json:"name"`
	PubIp         string `json:"pubIp"`
	InnerIp       string `json:"innerIp"`
	TcpPort       int64  `json:"tcpPort"`
	HttpPort      int64  `json:"httpPort"`
	GrpcPort      int64  `json:"grpcPort"`
	MongoUri      string `json:"mongoUri"`
	MongoAuth     string `json:"mongoAuth"`
	DbName        string `json:"dbName"`
	KafkaUri      string `json:"kafkaUri"`
	RedisUri      string `json:"redisUri"`
	RedisPass     string `json:"redisPass"`
	RedisPort     int64  `json:"redisPort"`
	RedisMeshUri  string `json:"redisMeshUri"`
	RedisMeshPass string `json:"redisMeshPass"`
	RedisMeshPort int64  `json:"meshPort"`
	GatewayUri    string `json:"gatewayUri"`
	FightType     string `json:"fightType"`
}

type DockerComposeFile struct {
	ImageService     string `json:"imageService"`
	ImageTag         string `json:"imageTag"`
	JsonConfigVolume string `json:"jsonConfigVolume"`
	ImageName        string `json:"imageName"`
	ImageUri         string `json:"imageUri"`
}
