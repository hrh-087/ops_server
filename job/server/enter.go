package server

import (
	"fmt"
	"github.com/hibiken/asynq"
	"ops-server/global"
)

func NewAsynqClinet() *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", global.OPS_CONFIG.Asynq.Addr, global.OPS_CONFIG.Asynq.Port),
		Password: global.OPS_CONFIG.Asynq.Pass,
		DB:       global.OPS_CONFIG.Asynq.Db,
	})
}
