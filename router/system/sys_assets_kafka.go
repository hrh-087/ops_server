package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type AssetsKafkaRouter struct {
}

func (i *AssetsKafkaRouter) InitAssetsKafkaRouter(Router *gin.RouterGroup) {
	router := Router.Group("assets").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("assets")

	routerApi := v1.ApiGroupApp.SystemApiGroup.AssetsKafkaApi

	{
		router.GET("kafka/:id/", routerApi.GetAssetsKafkaById)
		router.POST("kafka/", routerApi.CreateKafka)
		router.PUT("kafka/", routerApi.UpdateKafka)
		router.DELETE("kafka/", routerApi.DeleteKafka)
	}
	{
		routerWithoutRecord.GET("kafka/", routerApi.GetAssetsKafkaList)
		routerWithoutRecord.GET("kafka/all/", routerApi.GetAssetsKafkaAll)
	}
}
