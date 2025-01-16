package system

type SysProjectAuthority struct {
	SysProjectId            uint `json:"sys_project_id" gorm:"comment:项目ID;column:sys_project_id"`
	SysAuthorityAuthorityId uint `json:"sys_authority_authority_id" gorm:"comment:角色ID;column:sys_authority_authority_id"`
}

func (s SysProjectAuthority) TableName() string {
	return "sys_project_authority"
}
