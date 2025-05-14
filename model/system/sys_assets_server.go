package system

import (
	"github.com/gofrs/uuid/v5"
	"ops-server/global"
)

type SysAssetsServer struct {
	ProjectModel
	UUID               uuid.UUID `json:"uuid" gorm:"comment:唯一标识"`
	InstanceId         string    `json:"instanceId" gorm:"comment:示例id·"`
	ServerName         string    `json:"serverName" gorm:"comment:实例名称"`
	PubIp              string    `json:"pubIp" gorm:"comment:外网ip"`
	PrivateIp          string    `json:"privateIp" gorm:"comment:内网ip"`
	InstanceStates     int64     `json:"instanceStates" gorm:"comment:状态"`
	SSHPort            string    `json:"sshPort" gorm:"comment:ssh端口"`
	VpcId              string    `json:"vpcId" gorm:"comment:vpcId"`
	SubVpcId           string    `json:"subVpcId" gorm:"comment:子网id"`
	HostType           int8      `json:"hostType" gorm:"default:1; comment:服务器类型 1 本地服务器 2 云服务器 "`
	SysAssetDeviceInfo `json:"deviceInfo" gorm:"embedded;comment:附加属性"`
	Status             int64                 `json:"status" gorm:"comment:状态 0:待初始化 1 正常运行 2待回收 3 已回收;default:1"`
	CloudProduceId     uint                  `json:"cloudProduceId" gorm:"default:0;comment:云产商id"`
	Cloud              SysCloudProduce       `json:"cloudProduce" gorm:"foreignKey:CloudProduceId;references:ID"`
	PlatformId         uint                  `json:"platformId" gorm:"comment:渠道id"`
	Platform           SysGamePlatform       `json:"platform" gorm:"foreignKey:PlatformId;references:ID"`
	Ports              []SysAssetsServerPort `json:"ports" gorm:"foreignKey:ServerId"`
	ServerType         int8                  `json:"serverType" gorm:"default:1; comment:服务器类型 1 游戏服务 2 rim服务 3 运维后台"`
}

func (a *SysAssetsServer) TableName() string {
	return "sys_assets_servers"
}

type SysAssetsServerPort struct {
	global.OpsModel
	ServerId uint  `json:"serverId" gorm:"comment:服务器id"`
	Port     int64 `json:"port" gorm:"comment:端口"`
}

func (a *SysAssetsServerPort) TableName() string {
	return "sys_assets_server_ports"
}

type SysAssetsRedis struct {
	ProjectModel
	PlatformId         uint            `json:"platformId" gorm:"comment:渠道id"`
	Platform           SysGamePlatform `json:"platform" gorm:"foreignKey:PlatformId;references:ID"`
	SysAssetDeviceInfo `json:"deviceInfo" gorm:"embedded;comment:附加属性"`

	Name      string `json:"name" gorm:"comment:redis名称"`
	Host      string `json:"host" gorm:"comment:连接地址"`
	Port      int64  `json:"port" gorm:"comment:连接端口"`
	Password  string `json:"password" gorm:"comment:连接密码"`
	IsCluster bool   `json:"isCluster" gorm:"comment:是否集群"`
}

func (a *SysAssetsRedis) TableName() string {
	return "sys_assets_redis"
}

type SysAssetsMongo struct {
	ProjectModel
	PlatformId         uint            `json:"platformId" gorm:"comment:渠道id"`
	Platform           SysGamePlatform `json:"platform" gorm:"foreignKey:PlatformId;references:ID"`
	SysAssetDeviceInfo `json:"deviceInfo" gorm:"embedded;comment:附加属性"`

	Name string `json:"name" gorm:"comment:名称"`
	Host string `json:"host" gorm:"comment:连接地址"`
	Auth string `json:"auth" gorm:"comment:认证信息:authSource=admin&replicaSet=replica"`
}

func (a *SysAssetsMongo) TableName() string {
	return "sys_assets_mongo"
}

type SysAssetsMysql struct {
	ProjectModel
	PlatformId         uint            `json:"platformId" gorm:"comment:渠道id"`
	Platform           SysGamePlatform `json:"platform" gorm:"foreignKey:PlatformId;references:ID"`
	SysAssetDeviceInfo `json:"deviceInfo" gorm:"embedded;comment:附加属性"`

	Name string `json:"name" gorm:"comment:名称"`
	Host string `json:"host" gorm:"comment:连接地址"`
	Port int64  `json:"port" gorm:"comment:连接端口"`
	User string `json:"user" gorm:"comment:用户名"`
	Pass string `json:"pass" gorm:"comment:密码"`
}

func (a *SysAssetsMysql) TableName() string {
	return "sys_assets_mysql"
}

type SysAssetsKafka struct {
	ProjectModel
	PlatformId         uint            `json:"platformId" gorm:"comment:渠道id"`
	Platform           SysGamePlatform `json:"platform" gorm:"foreignKey:PlatformId;references:ID"`
	SysAssetDeviceInfo `json:"deviceInfo" gorm:"embedded;comment:附加属性"`

	Name string `json:"name" gorm:"comment:名称"`
	Host string `json:"host" gorm:"comment:连接地址"`
	User string `json:"user" gorm:"comment:用户名"`
	Pass string `json:"pass" gorm:"comment:密码"`
}

func (a *SysAssetsKafka) TableName() string {
	return "sys_assets_kafka"
}

type SysAssetDeviceInfo struct {
	System     string `json:"systemOs" gorm:"comment:操作系统"`    // 操作系统
	CpuNum     string `json:"cpuNum" gorm:"comment:cpu数量"`     // cpu数量
	CpuCoreNum string `json:"cpuCoreNum" gorm:"comment:cpu核数"` // 单cpu核数
	CpuTotal   string `json:"cpuTotal" gorm:"comment:cpu总核数"`  // 总核数
	Mem        string `json:"mem" gorm:"comment:内存大小"`         // 内存大小/G
	Disk       string `json:"disk" gorm:"comment:磁盘容量"`        // 磁盘容量/G
}
