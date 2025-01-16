package request

import (
	"ops-server/model/common/request"
	"ops-server/model/system"
)

type SearchCloudProduceParams struct {
	system.SysCloudProduce
	request.PageInfo
}
