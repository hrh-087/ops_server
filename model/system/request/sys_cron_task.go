package request

import (
	"ops-server/model/common/request"
	"ops-server/model/system"
)

type CronTaskParams struct {
	request.PageInfo
	system.CronTask
}
