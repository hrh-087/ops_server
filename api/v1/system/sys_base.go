package system

import (
	"crypto/md5"
	"github.com/gin-gonic/gin"
	"ops-server/global"
	"ops-server/model/common/response"
	"ops-server/utils"
	"path/filepath"
	"time"
)

type BaseApi struct {
}

const FileSize = 1024 * 1024 * 5

func (*BaseApi) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("解析文件失败", c)
		return
	}

	if file.Size > FileSize {
		response.FailWithMessage("文件过大", c)
		return
	}

	var fileType string
	ext := filepath.Ext(file.Filename)

	switch ext {
	case ".zip":
		fileType = "hotUpdate"
	case ".xlsx":
		fileType = "excel"
	default:
		response.FailWithMessage("文件类型错误", c)
		return
	}
	filePath := filepath.Join(global.OPS_CONFIG.Local.Path, fileType, time.Now().Format("2006-01-02"))
	//filePath := global.OPS_CONFIG.Local.Path + fileType + "/" + time.Now().Format("2006-01-02") + "/"
	filename := md5.Sum([]byte(file.Filename))

	if err := c.SaveUploadedFile(file, filepath.Join(filePath, utils.MD5V(filename[:])+ext)); err != nil {
		response.FailWithMessage("上传文件失败", c)
		return
	}
	response.OkWithDetailed(gin.H{
		"filePath":       filepath.Join(filePath, utils.MD5V(filename[:])+ext),
		"sourceFileName": file.Filename,
	}, "上传成功", c)
}
