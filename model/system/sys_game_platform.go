package system

type SysGamePlatform struct {
	ProjectModel
	PlatformName     string `json:"platformName" gorm:"comment:游戏平台名称"`
	PlatformCode     string `json:"platformCode" gorm:"comment:游戏平台代码"`
	PlatformDescribe string `json:"platformDescribe" gorm:"comment:平台描述"`
	ImageTag         string `json:"imageTag" gorm:"镜像tag"`
	ImageUri         string `json:"imageUri" gorm:"镜像uri"`
	// lts配置
	LtsLogGroupId  string `json:"ltsLogGroupId" gorm:"comment:日志组id"`
	LtsLogStreamId string `json:"ltsLogStreamId" gorm:"comment:日志流id"`
	CloudRegionId  string `json:"cloudRegionId" gorm:"comment:云区域id"`
	CloudProjectId string `json:"cloudProjectId" gorm:"comment:云项目id"`
	CloudSecretId  string `json:"cloudSecretId" gorm:"comment:云secretId"`
	CloudSecretKey string `json:"cloudSecretKey" gorm:"comment:云secretKey"`

	// 敏感词配置
	FilterUrl   string `json:"filterUrl" gorm:"comment:敏感词地址"`
	FilterToken string `json:"filterToken" gorm:"type:text;comment:敏感词token"`
	// gm地址
	GmUrl      string `json:"gmUrl" gorm:"comment:gm地址"`
	GatewayUrl string `json:"gatewayUrl" gorm:"comment:网关地址"`
}

func (s *SysGamePlatform) TableName() string {
	return "sys_game_platforms"
}
