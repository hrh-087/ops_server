package response

// gm开关
type GmSwitch struct {
	ServerId int    `json:"serverId"`
	TypeKey  string `json:"typeKey"`
	TypeName string `json:"typeName"`
	State    bool   `json:"state"`
}

type RankOpenConfig struct {
	// 排行榜关闭时间，格式2025-03-04T16:20:00+08:00
	EndTime string `json:"endTime,omitempty"`
	// id，排行榜开启配置主键id，自增
	ID int `json:"id,omitempty"`
	// 排行榜id，对应策划配置表List的主键id
	RankID int `json:"rankId,omitempty"`
	// 排行榜展示数量
	ShowCount int `json:"showCount,omitempty"`
	// 排行榜开启时间，格式2025-03-04T16:20:00+08:00
	StartTime string `json:"startTime,omitempty"`

	CloseTime string `json:"closeTime,omitempty"`
}

type RankRewardConfig struct {
	// id，排行榜奖励配置主键id，自增
	ID int `json:"id"`
	// 排名
	Rank int `json:"rank"`
	// 排行榜id，对应策划配置表List的主键id
	OpenId int `json:"openId"`
	// 奖励，rank排名时获得的奖励，格式map，key是道具id，value是道具数量
	Rewards map[string]int `json:"rewards"`
}

// 活动管理
type GameServerCron struct {
	ServerId int    `json:"serverId"`
	Key      string `json:"key"`
	Cron     string `json:"cron"`
}
