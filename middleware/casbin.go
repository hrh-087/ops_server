package middleware

import (
	"github.com/gin-gonic/gin"
	"ops-server/global"
	"ops-server/model/common/response"
	"ops-server/service"
	"ops-server/utils"
	"strconv"
	"strings"
)

var casbinService = service.ServiceGroupApp.SystemServiceGroup.CasbinService

func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		waitUse, _ := utils.GetClaims(c)

		// 获取请求的URI
		path := c.Request.URL.Path

		obj := strings.TrimPrefix(path, global.OPS_CONFIG.System.RouterPrefix)
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := strconv.Itoa(int(waitUse.AuthorityId))
		e := casbinService.Casbin()

		success, _ := e.Enforce(sub, obj, act)
		//fmt.Println("sub: ", sub, "obj: ", obj, "act: ", act)
		if !success {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
