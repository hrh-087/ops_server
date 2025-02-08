package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type AssetsMongoRouter struct {
}

func (i *AssetsMongoRouter) InitAssetsMongoRouter(Router *gin.RouterGroup) {
	router := Router.Group("assets").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("assets")

	routerApi := v1.ApiGroupApp.SystemApiGroup.AssetsMongoApi

	{
		router.GET("mongo/:id/", routerApi.GetAssetsMongoById)
		router.POST("mongo/", routerApi.CreateMongo)
		router.PUT("mongo/", routerApi.UpdateMongo)
		router.DELETE("mongo/", routerApi.DeleteMongo)
	}
	{
		routerWithoutRecord.GET("mongo/", routerApi.GetAssetsMongoList)
		routerWithoutRecord.GET("mongo/all/", routerApi.GetAssetsMongoAll)
	}
}
