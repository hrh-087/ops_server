package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
)

type JobRouter struct {
}

func (s *JobRouter) InitJobRouter(Router *gin.RouterGroup) {
	router := Router.Group("job")
	routerWithoutRecord := Router.Group("job")

	routerApi := v1.ApiGroupApp.SystemApiGroup.JobApi

	{
		router.POST("info/", routerApi.GetJobById)
	}
	{
		routerWithoutRecord.GET("/", routerApi.GetJobList)
	}
}
