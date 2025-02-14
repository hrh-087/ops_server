package system

type SysTask struct {
	ProjectModel

	Name     string `json:"name" gorm:"comment:任务名称"`
	TaskType string `json:"taskType" gorm:"unique;comment:任务类型"`
}

func (t *SysTask) TableName() string {
	return "sys_tasks"
}
