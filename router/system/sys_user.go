package system

import (
	"github.com/gin-gonic/gin"
	v1 "ops-server/api/v1"
	"ops-server/middleware"
)

type UserRouter struct {
}

// InitUserRouter 初始化用户路由
func (u *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	//userRouter := Router.Group("user")
	userRouter := Router.Group("user").Use(middleware.OperationRecord())
	userRouterWithoutRecord := Router.Group("user")

	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi

	{
		userRouter.POST("register/", baseApi.Register)                     // 管理员注册账号
		userRouter.POST("changePassword/", baseApi.ChangePassword)         // 用户修改密码
		userRouter.POST("setUserAuthority/", baseApi.SetUserAuthority)     // 设置用户权限
		userRouter.POST("setUserProject/", baseApi.SetUserProject)         // 设置用户权限
		userRouter.DELETE("deleteUser/", baseApi.DeleteUser)               // 删除用户
		userRouter.PUT("setUserInfo/", baseApi.SetUserInfo)                // 设置用户信息
		userRouter.PUT("setSelfInfo/", baseApi.SetSelfInfo)                // 设置自身信息
		userRouter.POST("setUserAuthorities/", baseApi.SetUserAuthorities) // 设置用户权限组
		userRouter.POST("resetPassword/", baseApi.ResetPassword)           // 设置用户权限组

	}
	{
		userRouterWithoutRecord.POST("getUserList/", baseApi.GetUserList) // 分页获取用户列表
		userRouterWithoutRecord.GET("getUserinfo/", baseApi.GetUserInfo)  // 获取自身信息
	}
}
