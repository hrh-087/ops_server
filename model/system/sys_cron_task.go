package system

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type CronTask struct {
	CreatedAt time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"comment:更新时间"`

	ProjectId  uint       `json:"projectId" gorm:"comment:项目id"`
	SysProject SysProject `json:"-" gorm:"foreignKey:ProjectId;references:ID"`

	Tasks []JobTask `json:"tasks" gorm:"foreignKey:JobId;references:CronTaskId"`

	CronTaskId   uuid.UUID `json:"cronTaskId" form:"cronTaskId" gorm:"primary_key"`
	TaskId       string    `json:"taskId" gorm:"index;comment:任务id"`
	Name         string    `json:"name" gorm:"comment:任务名称"`
	TaskTypeName string    `json:"taskTypeName" gorm:"comment:任务类型名称"`
	Type         int       `json:"type" gorm:"comment:任务类型, 1 一次性任务 2 周期性任务"`
	ExecTime     time.Time `json:"execTime" gorm:"default:null;comment:执行时间"` // 一次性执行任务直接指定执行时间
	CronRule     string    `json:"cronRule" gorm:"comment:cron规则"`            // 周期性任务需要设置cron规则
	Status       int       `json:"status" gorm:"default:2;comment:任务状态 1 开启 2 关闭"`
	Describe     string    `json:"describe" gorm:"comment:任务描述"`
	Creator      string    `json:"creator" gorm:"comment:创建者"`
}

func (*CronTask) TableName() string {
	return "sys_cron_tasks"
}
