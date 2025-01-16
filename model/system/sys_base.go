package system

import "ops-server/global"

type ProjectModel struct {
	global.OpsModel
	ProjectId  uint       `json:"projectId" gorm:"comment:项目id"`
	SysProject SysProject `json:"-" gorm:"foreignKey:ProjectId;references:ID"`
}
