package request

import (
	"ops-server/model/common/request"
	"ops-server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
