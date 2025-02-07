package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
)

type OperationRecordRouter struct{}

func (s *OperationRecordRouter) InitSysOperationRecordRouter(Router *gin.RouterGroup) {
	operationRecordRouter := Router.Group("system")
	authorityMenuApi := v1.ApiGroupApp.SystemApiGroup.OperationRecordApi
	{
		//operationRecordRouter.POST("record/", authorityMenuApi.CreateSysOperationRecord)        // 新建SysOperationRecord
		operationRecordRouter.DELETE("record/", authorityMenuApi.DeleteSysOperationRecord)            // 删除SysOperationRecord
		operationRecordRouter.DELETE("record/batch/", authorityMenuApi.DeleteSysOperationRecordByIds) // 批量删除SysOperationRecord
		//operationRecordRouter.GET("record/", authorityMenuApi.FindSysOperationRecord)           // 根据ID获取SysOperationRecord
		operationRecordRouter.GET("record/", authorityMenuApi.GetSysOperationRecordList) // 获取SysOperationRecord列表
	}
}
