package system

import (
	"ops-server/global"
)

type SysProject struct {
	global.OpsModel
	ProjectName     string         `json:"projectName" gorm:"unique;comment:项目名称"`
	Authorities     []SysAuthority `json:"authorities" gorm:"many2many:sys_project_authority;"`
	ConfigDir       string         `json:"configDir" gorm:"comment:项目配置文件目录"`
	SvnUrl          string         `json:"svnUrl" gorm:"comment:svn地址"`
	GmUrl           string         `json:"gmUrl" gorm:"comment:gm地址"`
	GatewayUrl      string         `json:"gatewayUrl" gorm:"comment:网关地址"`
	ClientConfigDir string         `json:"clientConfigDir" gorm:"comment:客户端配置文件目录"`
	ClientSvnUrl    string         `json:"clientSvnUrl" gorm:"comment:客户端svn地址"`
	Status          int            `json:"status" gorm:"comment:项目状态 0 未初始化 1 已初始化 2 初始化失败"`
	IsTest          bool           `json:"isTest" gorm:"comment:是否测试项目"`
	WebHook         string         `json:"webHook" gorm:"comment:webhook 测试项目用于通知钉钉"`
}

func (SysProject) TableName() string {
	return "sys_projects"
}
