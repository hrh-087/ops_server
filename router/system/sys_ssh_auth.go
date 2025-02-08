package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type SshAuthRouter struct{}

func (s *SshAuthRouter) InitSshAuthRouter(Router *gin.RouterGroup) {
	router := Router.Group("system").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("system")

	sshAuthApi := v1.ApiGroupApp.SystemApiGroup.SshAuthApi

	{
		router.POST("key/", sshAuthApi.CreateSshAuth)
		router.DELETE("key/", sshAuthApi.DeleteSshAuth)
		router.PUT("key/", sshAuthApi.UpdateSshAuth)
	}
	{
		routerWithoutRecord.GET("key/", sshAuthApi.GetSshAuthList)
		routerWithoutRecord.GET("key/:id/", sshAuthApi.GetSshAuthById)
	}
}
