package request

import (
	"ops-server/model/common/request"
)

type JobTaskParams struct {
	request.PageInfo
	JobId string `json:"jobId" form:"jobId"`
}
