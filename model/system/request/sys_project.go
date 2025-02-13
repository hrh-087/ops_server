package request

import (
	"ops-server/model/common/request"
	"ops-server/model/system"
)

type SearchProjectParams struct {
	system.SysProject
	request.PageInfo
}

type AddAuthorityProject struct {
	AuthorityId uint                `json:"authorityId"`
	ProjectIds  []system.SysProject `json:"projectIds:q"`
}
