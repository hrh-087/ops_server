package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
)

type AssetsMysqlRouter struct {
}

func (i *AssetsMysqlRouter) InitAssetsMysqlRouter(Router *gin.RouterGroup) {
	router := Router.Group("assets")
	routerWithoutRecord := Router.Group("assets")

	routerApi := v1.ApiGroupApp.SystemApiGroup.AssetsMysqlApi

	{
		router.GET("mysql/:id/", routerApi.GetAssetsMysqlById)
		router.POST("mysql", routerApi.CreateMysql)
		router.PUT("mysql", routerApi.UpdateMysql)
		router.DELETE("mysql", routerApi.DeleteMysql)
	}
	{
		routerWithoutRecord.GET("mysql", routerApi.GetAssetsMysqlList)
		routerWithoutRecord.GET("mysql/all", routerApi.GetAssetsMysqlAll)
	}
}
