package response

import "ops-server/model/system"

type SysAPIResponse struct {
	Data system.SysApi `json:"data"`
}

type SysAPIListResponse struct {
	Data []system.SysApi `json:"data"`
}

type SysSyncApis struct {
	NewApis    []system.SysApi `json:"newApis"`
	DeleteApis []system.SysApi `json:"deleteApis"`
}
