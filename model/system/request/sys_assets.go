package request

import (
	"ops-server/model/common/request"
)

type SearchAssetsServerParams struct {
	request.NameAndPlatformSearch
	request.PageInfo
}
