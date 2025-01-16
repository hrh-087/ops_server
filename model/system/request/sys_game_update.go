package request

import (
	"github.com/gofrs/uuid/v5"
	"ops-server/model/system"
)

type GameUpdateTaskParams struct {
	TaskTypeName string
	StepName     string
	JobId        uuid.UUID
	Command      string
	Params       string
}

type HotUpdateParams struct {
	ServerType int8  `json:"serverType"`
	ServerList []int `json:"serverList"`
}

type GameUpdateParams struct {
	system.GameUpdate
	HotUpdateParams
}
