package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type CronTaskRouter struct {
}

func (CronTaskRouter) InitCronTaskRouter(Router *gin.RouterGroup) {
	router := Router.Group("job").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("job")

	routerApi := v1.ApiGroupApp.SystemApiGroup.CronTaskApi

	{
		router.POST("cron/", routerApi.CreateCronTask)
		router.PUT("cron/", routerApi.UpdateCronTask)
		router.DELETE("cron/", routerApi.DeleteCronTask)
		//router.POST("cron/exec/", routerApi.ExecCronTask)
	}

	{
		routerWithoutRecord.GET("cron/:id/", routerApi.GetCronTaskById)
		routerWithoutRecord.GET("cron/", routerApi.GetCronTaskList)
	}
}
