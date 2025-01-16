package system

import (
	"ops-server/global"
)

type SysProject struct {
	global.OpsModel
	ProjectName string         `json:"projectName" gorm:"comment:项目名称"`
	Authorities []SysAuthority `json:"authorities" gorm:"many2many:sys_project_authority;"`
}

func (SysProject) TableName() string {
	return "sys_projects"
}
