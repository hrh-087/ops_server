package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type SysTaskRouter struct {
}

func (s *SysTaskRouter) InitSysTaskRouter(Router *gin.RouterGroup) {
	router := Router.Group("job").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("job")

	routerApi := v1.ApiGroupApp.SystemApiGroup.SysTaskApi

	{
		router.POST("taskManage/", routerApi.CreateTask)
		router.PUT("taskManage/", routerApi.UpdateTask)
		router.DELETE("taskManage/", routerApi.DeleteTask)
		router.POST("taskManage/exec/", routerApi.ExecTask)
	}
	{
		routerWithoutRecord.GET("taskManage/", routerApi.GetTaskList)
		routerWithoutRecord.GET("taskManage/:id/", routerApi.GetTaskById)
	}
}
