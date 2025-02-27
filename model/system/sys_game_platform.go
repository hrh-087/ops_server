package system

type SysGamePlatform struct {
	ProjectModel
	PlatformName     string `json:"platformName" gorm:"comment:游戏平台名称"`
	PlatformCode     string `json:"platformCode" gorm:"comment:游戏平台代码"`
	PlatformDescribe string `json:"platformDescribe" gorm:"comment:平台描述"`
	ImageTag         string `json:"imageTag" gorm:"镜像tag"`
	ImageUri         string `json:"imageUri" gorm:"镜像uri"`
}

func (s *SysGamePlatform) TableName() string {
	return "sys_game_platforms"
}
