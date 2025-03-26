package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
)

type JobTaskRouter struct {
}

func (s *JobTaskRouter) InitJobTaskRouter(Router *gin.RouterGroup) {
	//router := Router.Group("job")
	routerWithoutRecord := Router.Group("job")

	routerApi := v1.ApiGroupApp.SystemApiGroup.JobTaskApi
	{

	}
	{
		routerWithoutRecord.POST("task/result/", routerApi.GetJobTaskResult)
		routerWithoutRecord.GET("task/", routerApi.GetJobTaskList)
	}
}
