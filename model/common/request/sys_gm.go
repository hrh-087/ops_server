package request

type GmSwitchParams struct {
	ServerId int    `json:"serverId" form:"serverId"` // 渠道id
	TypeKey  string `json:"typeKey" form:"typeKey"`   // 开关类型
	State    bool   `json:"state" form:"state"`       // 开关状态
}

// 排行榜

type GmRankOpenParams struct {
	ServerId int `json:"serverId" form:"serverId"` // 渠道id
}

type GmRankRewardParams struct {
	ServerId int `json:"serverId" form:"serverId"` // 渠道id
	Id       int `json:"id" form:"id"`             // 榜单配置的唯一id
}

type Reward struct {
	// 榜单奖励
	RewardId  string `json:"rewardId"`  // 奖励id
	RewardNum int    `json:"rewardNum"` // 奖励数量
}

type GmRankReward struct {
	Id      int      `json:"id"`      // 榜单奖励唯一id
	OpenId  int      `json:"openId"`  // 榜单id
	Rank    string   `json:"rank"`    // 榜单排名
	Rewards []Reward `json:"rewards"` // 奖励列表
}

type GmRankConfig struct {
	Id         int            `json:"id"`         // 榜单唯一id
	RankId     int            `json:"rankId"`     // 榜单id
	ShowCount  int            `json:"showCount"`  // 展示数量
	StartTime  string         `json:"startTime"`  // 开始时间
	EndTime    string         `json:"endTime"`    // 结束时间
	CloseTime  string         `json:"closeTime"`  // 榜单关闭时间
	RewardList []GmRankReward `json:"rewardList"` // 奖励列表
}

type GmRankConfigParams struct {
	ServerId   int            `json:"serverId" form:"serverId"`     // 渠道id
	RankConfig []GmRankConfig `json:"rankConfig" form:"rankConfig"` // 榜单列表
}

type Rank struct {
	RankId   int    `json:"rankId"`
	RankName string `json:"rankName"`
	RankType int    `json:"rankType"`
}

type Item struct {
	ItemId   int         `json:"itemId"`
	ItemName interface{} `json:"itemName"`
}

// 维度推送
type DimensionPushParams struct {
	ServerList []string `json:"platformList" form:"platformList"` // 推送的服务器列表
}

// 导号
type ImportPlayerDataParams struct {
	PlatformCode   string `json:"platformCode" form:"platformCode"`
	Account        string `json:"account" form:"account"`               // 账号
	ImportId       string `json:"importId" form:"importId"`             // 导入的account表id
	OutputPlayerId string `json:"outputPlayerId" form:"outputPlayerId"` // 导出的玩家id
}

// 游戏服活动管理
type GameServerActivityParams struct {
	ServerId      int    `json:"serverId" form:"serverId"`           // 渠道id
	Key           string `json:"key" form:"key"`                     // 活动key
	Cron          string `json:"cron" form:"cron"`                   // 定时任务rule * * * * * *
	ActivityExtra string `json:"activityExtra" form:"activityExtra"` // 活动额外参数
}
