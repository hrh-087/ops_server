package request

import (
	"ops-server/model/system"
)

type AddMenuAuthorityInfo struct {
	Menus       []system.SysBaseMenu `json:"menus"`
	AuthorityId uint                 `json:"authorityId"` // 角色ID
}

func DefaultMenu() []system.SysBaseMenu {
	//return []system.SysBaseMenu{{
	//	OpsModel:  global.OpsModel{ID: 999},
	//	ParentId:  0,
	//	Path:      "dashboard",
	//	Name:      "dashboard",
	//	Component: "view/dashboard/index.vue",
	//	Sort:      1,
	//	Meta: system.Meta{
	//		Title: "仪表盘",
	//		Icon:  "setting",
	//	},
	//}}
	return []system.SysBaseMenu{}
}
