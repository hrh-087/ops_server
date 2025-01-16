package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
)

type JobCommandRouter struct {
}

func (r *JobCommandRouter) InitJobCommandRouter(Router *gin.RouterGroup) {
	router := Router.Group("job")
	routerWithoutRecord := Router.Group("job")

	routerApi := v1.ApiGroupApp.SystemApiGroup.JobCommandApi

	{
		router.POST("command/", routerApi.CreateCommand)
		router.PUT("command/", routerApi.UpdateCommand)
		router.DELETE("command/", routerApi.DeleteCommand)
		router.POST("command/batchCommand/", routerApi.ExecBatchCommand)
	}

	{
		routerWithoutRecord.GET("command/:id/", routerApi.GetCommandById)
		routerWithoutRecord.GET("command/", routerApi.GetCommandList)
		routerWithoutRecord.GET("command/all/", routerApi.GetCommandAll)
	}
}
