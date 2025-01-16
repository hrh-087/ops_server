package system

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type Job struct {
	ProjectId  uint       `json:"projectId" gorm:"comment:项目id"`
	SysProject SysProject `json:"-" gorm:"foreignKey:ProjectId;references:ID"`

	Tasks []JobTask `json:"tasks" gorm:"foreignKey:JobId;references:JobId"`

	JobId    uuid.UUID `json:"jobId" gorm:"primary_key"`
	Name     string    `json:"name" gorm:"comment:任务名称"`
	Type     string    `json:"type" gorm:"comment:任务类型"`
	Status   int       `json:"status" gorm:"comment:任务状态 0 未开始 1 运行中 2已完成 3 运行失败"`
	Creator  string    `json:"creator" gorm:"comment:创建者"`
	ExecTime float64   `json:"execTime" gorm:"comment:总耗时"`
	CreateAt time.Time `json:"createAt" gorm:"comment:创建时间"`
}

func (j *Job) TableName() string {
	return "jobs"
}

type JobTask struct {
	JobId uuid.UUID `json:"jobId" gorm:"index;comment:任务id"`

	AsynqId string    `json:"asynqId" gorm:"index;comment:asynq任务id"`
	TaskId  uuid.UUID `json:"taskId" gorm:"primary_key"`

	HostName string    `json:"hostName" gorm:"comment:服务器名称"`
	HostIp   string    `json:"hostIp" gorm:"comment:服务器ip"`
	ExecTime float64   `json:"execTime" gorm:"comment:耗时"`
	Status   string    `json:"status" gorm:"comment:任务状态"`
	CreateAt time.Time `json:"createAt" gorm:"comment:创建时间"`
}

func (j *JobTask) TableName() string {
	return "job_tasks"
}

type JobCommand struct {
	ProjectModel
	Name        string `json:"name" gorm:"comment:命令名称"`
	Command     string `json:"command" gorm:"type:text;comment:命令"`
	CommandType int8   `json:"commandType" gorm:"default:1;comment:命令类型"`
	Description string `json:"description" gorm:"type:text;comment:描述"`
	UseBatch    bool   `json:"useBatch" gorm:"comment:是否允许批量执行"`
}

func (c *JobCommand) TableName() string {
	return "job_commands"
}
