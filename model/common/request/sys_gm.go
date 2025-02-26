package request

type GmSwitchParams struct {
	ServerId int    `json:"serverId" form:"serverId"`
	TypeKey  string `json:"typeKey" form:"typeKey"`
	State    bool   `json:"state" form:"state"`
}
