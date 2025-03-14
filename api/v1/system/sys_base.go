package system

import (
	"crypto/md5"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"net/http"
	"ops-server/global"
	"ops-server/model/common/request"
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

func (*BaseApi) GenerateExcel(c *gin.Context) {

	var params request.ExcelTypeParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取参数失败"})
		return
	}

	if err := utils.Verify(params, utils.ExcelTypeVerify); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "参数验证失败"})
		return
	}

	projectId := c.GetString("projectId")
	if projectId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取项目ID失败"})
		return
	}

	// 创建 Excel
	f := excelize.NewFile()
	annotationMap := make(map[string]string)

	itemData, err := global.OPS_REDIS.HGetAll(c, "item_2").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取数据失败"})
		return
	}

	//excelType := "item"
	switch params.ExcelType {
	case "item":
		sheetIndex := 0
		for key, value := range itemData {
			switch key {
			case "item":
				annotationMap = map[string]string{
					"itemId":   "道具id",
					"itemName": "道具名称",
				}

			case "rank":
				annotationMap = map[string]string{
					"rankId":   "榜单id",
					"rankName": "榜单名称",
					"rankType": "榜单类型",
				}
			}

			var itemList []interface{}
			if err = json.Unmarshal([]byte(value), &itemList); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "解析数据失败"})
				return
			}
			err = utils.CreateSheet(f, key, annotationMap, itemList)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建 Sheet 失败"})
				return
			}

			if sheetIndex == 0 {
				index, err := f.NewSheet(key)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建 Sheet 失败"})
					return
				}
				f.SetActiveSheet(index)
			}

			sheetIndex++
		}
	}
	f.DeleteSheet("Sheet1")

	// 返回 Excel 文件
	c.Writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=export.xlsx")
	_ = f.Write(c.Writer)

}
