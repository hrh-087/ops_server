package inspector

import (
	"fmt"
	"github.com/hibiken/asynq"
	"ops-server/global"
)

func NewAsynqInspector() *asynq.Inspector {
	return asynq.NewInspector(asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", global.OPS_CONFIG.Asynq.Addr, global.OPS_CONFIG.Asynq.Port),
		Password: global.OPS_CONFIG.Asynq.Pass,
		DB:       global.OPS_CONFIG.Asynq.Db,
	})
}
