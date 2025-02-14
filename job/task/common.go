package task

import (
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	"github.com/hibiken/asynq"
	"ops-server/global"
)

type CommonTaskParams struct {
	TaskId    uuid.UUID
	ProjectId uint
	HostId    uint
}

func NewCommonTask(taskType string, params CommonTaskParams) (*asynq.TaskInfo, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	task := NewTask(taskType, payload)
	return global.AsynqClinet.Enqueue(task)
}
