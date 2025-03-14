package request

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页大小
	Keyword  string `json:"keyword" form:"keyword"`   //关键字
}

type NameAndPlatformSearch struct {
	Name       string `json:"name" form:"name"`
	PlatformId uint   `json:"platformId" form:"platformId"`
}

// GetById Find by id structure
type GetById struct {
	ID int `json:"id" form:"id"` // 主键ID
}

func (r *GetById) Uint() uint {
	return uint(r.ID)
}

type GetAuthorityId struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` //角色ID
}

type IdsReq struct {
	Ids []int `json:"ids" form:"ids"`
}

type UploadConfigParams struct {
	Data map[string][]interface{} `json:"data"`
}

type GmItemConfigParams struct {
	ItemType string `json:"itemType"`
}
