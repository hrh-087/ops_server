package request

import (
	"ops-server/model/common/request"
	"ops-server/model/system"
)

type CronTaskParams struct {
	request.PageInfo
	system.CronTask
}

type ExecCronTaskParams struct {
	ExecType int    `json:"execType"` // 1 开启 2 关闭
	CronRule string `json:"cronRule"` // 周期性定时任务规则
	ExecTime string `json:"execTime"` // 一次性定时任务执行时间
}
