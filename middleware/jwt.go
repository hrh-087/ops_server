package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"ops-server/model/common/response"
	"ops-server/service"
	"ops-server/utils"
)

var jwtService = service.ServiceGroupApp.SystemServiceGroup.JwtService

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := utils.GetToken(c)
		if token == "" {
			response.NoAuth("未登录或非法访问", c)
			c.Abort()
			return
		}
		// 多点登录检测
		/*if jwtService.IsBlacklist(token) {
			response.NoAuth("您的帐户异地登陆或令牌失效", c)
			utils.ClearToken(c)
			c.Abort()
			return
		}*/

		j := utils.NewJWT()
		// parseToKEN 解析token包含的信息

		claims, err := j.ParseToken(token)

		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				response.NoAuth("授权已过期", c)
				utils.ClearToken(c)
				c.Abort()
				return
			}
			response.NoAuth(err.Error(), c)
			utils.ClearToken(c)
			c.Abort()
			return
		}

		c.Set("claims", claims)

		c.Next()
	}
}
