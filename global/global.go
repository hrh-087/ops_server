package global

import (
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"ops-server/config"
)

var (
	OPS_DB     *gorm.DB
	OPS_CONFIG config.Server
	OPS_REDIS  redis.UniversalClient
	OPS_VP     *viper.Viper

	OPS_ROUTERS gin.RoutesInfo

	// OPS_LOG
	OPS_LOG        *zap.Logger
	AsynqClient    *asynq.Client
	AsynqInspect   *asynq.Inspector
	AsynqScheduler *asynq.Scheduler
)
