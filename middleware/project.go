package middleware

import (
	"github.com/gin-gonic/gin"
	"ops-server/model/common/response"
	"ops-server/service"
	"ops-server/utils"
	"strconv"
)

var projectService = service.ServiceGroupApp.SystemServiceGroup.ProjectService

func ProjectAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		waitUse, _ := utils.GetClaims(c)

		roleId := strconv.Itoa(int(waitUse.AuthorityId))
		projectId := c.GetHeader("X-Project-Id")

		if projectId == "" {
			response.FailWithMessage("未设置项目ID", c)
			c.Abort()
			return
		}

		// 检测 projectId 是否是纯数字
		if _, err := strconv.ParseInt(projectId, 10, 64); err != nil {
			response.FailWithMessage("无法解析项目ID", c)
			c.Abort()
			return
		}

		err := projectService.CheckProject(roleId, projectId)

		if err != nil {
			response.FailWithMessage("项目权限不足", c)
			c.Abort()
			return
		}

		c.Set("projectId", projectId)
		c.Next()
	}
}
