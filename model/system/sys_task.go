package system

import "ops-server/global"

type SysTask struct {
	global.OpsModel

	Name     string `json:"name" gorm:"comment:任务名称"`
	TaskType string `json:"taskType" gorm:"comment:任务类型"`
}

func (t *SysTask) TableName() string {
	return "sys_tasks"
}
