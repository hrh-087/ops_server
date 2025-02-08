package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type CloudProduceRouter struct {
}

func (p *CloudProduceRouter) InitCloudProduceRouter(Router *gin.RouterGroup) {
	cloudRouter := Router.Group("cloud").Use(middleware.OperationRecord())
	cloudRouterWithoutRecord := Router.Group("cloud")

	cloudRouterApi := v1.ApiGroupApp.SystemApiGroup.CloudProduceApi

	{
		cloudRouter.POST("cloudProduce/", cloudRouterApi.CreateCloudProduce)
		cloudRouter.GET("cloudProduce/:id/", cloudRouterApi.GetCloudProduceById)
		cloudRouter.PUT("cloudProduce/", cloudRouterApi.UpdateCloudProduce)
		cloudRouter.DELETE("cloudProduce/", cloudRouterApi.DeleteCloudProduce)
	}
	{
		cloudRouterWithoutRecord.GET("cloudProduce/", cloudRouterApi.GetCloudProduceList)
		cloudRouterWithoutRecord.GET("cloudProduce/all/", cloudRouterApi.GetCloudProduceAll)

	}

}
