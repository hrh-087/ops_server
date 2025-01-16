package request

import (
	"ops-server/model/common/request"
)

type SearchGameParams struct {
	request.NameAndPlatformSearch
	request.PageInfo
}
