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
	RedisMeshUri  string `json:"redisMeshUri"`
	RedisMeshPass string `json:"redisMeshPass"`
}
