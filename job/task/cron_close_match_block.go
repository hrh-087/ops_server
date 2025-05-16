package task

import (
	"context"
	"github.com/hibiken/asynq"
	"log"
	"ops-server/global"
)

func HandleCloseMatchBlock(ctx context.Context, t *asynq.Task) error {
	log.Println("cron close match block1")
	global.OPS_LOG.Info("cron close match block")
	return nil
}
