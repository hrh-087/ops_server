package system

import (
	"ops-server/global"
)

type SysProject struct {
	global.OpsModel
	ProjectName string         `json:"projectName" gorm:"unique;comment:项目名称"`
	Authorities []SysAuthority `json:"authorities" gorm:"many2many:sys_project_authority;"`
	ConfigDir   string         `json:"configDir" gorm:"comment:项目配置文件目录"`
	SvnUrl      string         `json:"svnUrl" gorm:"comment:svn地址"`
	GmUrl       string         `json:"gmUrl" gorm:"comment:gm地址"`
	GatewayUrl  string         `json:"gatewayUrl" gorm:"comment:网关地址"`
	Status      int            `json:"status" gorm:"comment:项目状态 0 未初始化 1 已初始化 2 初始化失败"`
}

func (SysProject) TableName() string {
	return "sys_projects"
}
