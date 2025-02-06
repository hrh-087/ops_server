package system

type SysAssetsLb struct {
	ProjectModel
	Name         string `json:"name" gorm:"comment:负载均衡名称"`
	InstanceId   string `json:"instanceId" gorm:"unique;comment:实例id"`
	PubIp        string `json:"pubIp" gorm:"comment:公网ip"`
	Level        int8   `json:"level" gorm:"default:1;comment:等级"`
	PrivateIp    string `json:"privateIp" gorm:"comment:私网ip"`
	SubnetCidrId string `json:"subnetCidrId" gorm:"comment:子网id"`

	Listener       []SysAssetsListener `json:"listener" gorm:"foreignKey:LbId;references:ID"`
	PlatformId     uint                `json:"platformId" gorm:"comment:平台id"`
	Platform       SysGamePlatform     `json:"platform" gorm:"foreignKey:PlatformId;references:ID"`
	CloudProduceId uint                `json:"cloudProduceId" gorm:"index;comment:云产品id"`
	CloudProduce   SysCloudProduce     `json:"cloudProduce" gorm:"foreignKey:CloudProduceId;references:ID"`
}

func (SysAssetsLb) TableName() string {
	return "sys_assets_lb"
}

type SysAssetsListener struct {
	ProjectModel
	Name        string `json:"name" gorm:"index:idx_name_lb_id,unique;comment:监听器名称"`
	InstanceId  string `json:"instanceId" gorm:"comment:实例id"`
	Port        int    `json:"port" gorm:"comment:监听器端口"`
	Protocol    string `json:"protocol" gorm:"comment:监听器协议"`
	BackendIp   string `json:"backendIp" gorm:"comment:后端服务器ip"`
	BackendPort int    `json:"backendPort" gorm:"comment:后端服务器端口"`

	HostId uint            `json:"hostId" gorm:"comment:主机id"`
	Host   SysAssetsServer `json:"-" gorm:"foreignKey:HostId;references:ID"`
	LbId   uint            `json:"lbId" gorm:"index:idx_name_lb_id;comment:负载均衡id"`
	Lb     SysAssetsLb     `json:"-" gorm:"foreignKey:LbId;references:ID"`
}

func (SysAssetsListener) TableName() string {
	return "sys_assets_listener"
}
