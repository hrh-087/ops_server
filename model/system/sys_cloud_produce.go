package system

type SysCloudProduce struct {
	ProjectModel
	CloudName  string `json:"cloudName" gorm:"comment:云平台名称"`
	RegionId   string `json:"regionId" gorm:"comment:区域ID"`
	RegionName string `json:"regionName" gorm:"comment:区域名称"`
	SecretId   string `json:"secretId" gorm:"comment:SecretId"`
	SecretKey  string `json:"secretKey" gorm:"comment:SecretKey"`
	IsActive   bool   `json:"isActive" gorm:"comment:是否激活"`
	IsCloud    bool   `json:"isCloud" gorm:"comment:是否为云平台"`
}

func (SysCloudProduce) TableName() string {
	return "sys_cloud_produces"
}
