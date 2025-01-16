package system

type GameUpdate struct {
	ProjectModel
	Name          string `json:"name" gorm:"comment:名称"`
	Description   string `json:"description" gorm:"type:text;comment:描述"`
	UpdateType    int8   `json:"updateType" gorm:"comment:更新类型 1 正常更新 2 热更游戏文件 3 更新游戏配置"`
	SourceHotFile string `json:"sourceHotFile" gorm:"comment:源热更新文件"`
	HotFile       string `json:"hotFile" gorm:"comment:热更新文件"`
	Step          int8   `json:"step" gorm:"comment:步骤"`
	StepName      string `json:"stepName" gorm:"comment:步骤名称"`
	ServerType    int8   `json:"hotServerType" gorm:"comment:热更类型 1:游戏服 2: 游戏类型"`
	ServerList    string `json:"hotServerList" gorm:"comment:热更列表 多个以逗号分隔"`
	GameVersion   string `json:"gameVersion" gorm:"comment:游戏版本"`
	UpdateParams  string `json:"updateParams" gorm:"type:text;comment:更新参数"`
	TotalStep     int8   `json:"totalStep" gorm:"comment:总步骤数量"`
}

func (GameUpdate) TableName() string {
	return "sys_game_update"
}
