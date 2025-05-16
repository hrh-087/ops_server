package request

import (
	"ops-server/model/common/request"
	"ops-server/model/system"
)

type SearchGameParams struct {
	system.SysGameServer
	request.PageInfo
}

type SearchGameTypeParams struct {
	system.SysGameType
	request.PageInfo
}

type UpdateGameConfigParams struct {
	UpdateType int8   `json:"updateType"` // 1 更新所有已安装游戏服  2 更新gameIds中的游戏服
	GameIds    []int8 `json:"ids"`
}

type GameTaskParams struct {
	GameServerIds []uint `json:"gameServerIds"`
	TaskType      int8   `json:"taskType"`
}

type CopyGameTypeParams struct {
	ProjectId   uint  `json:"projectId"`
	GameTypeIds []int `json:"gameTypeIds"`
}
