package initialize

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"ops-server/global"
	"ops-server/model/system"
	"os"
	"path/filepath"
	"reflect"
)

type exportData interface {
	TableName() string
}

func ExportData() (err error) {
	var jsonDir string
	if global.OPS_CONFIG.Local.JsonDir == "" {
		jsonDir = "./jsonDir"
	} else {
		jsonDir = global.OPS_CONFIG.Local.JsonDir
	}

	if global.OPS_DB == nil {
		panic("db is nil")
	}

	var tables = []exportData{
		system.SysApi{},
		system.SysBaseMenu{},
	}

	for _, table := range tables {
		var data interface{}

		switch reflect.TypeOf(table) {
		case reflect.TypeOf(system.SysBaseMenu{}):
			var menuData []system.SysBaseMenu
			err = global.OPS_DB.Model(table).Order("id").Omit("created_at", "updated_at", "deleted_at").Find(&menuData).Error
			if err != nil {
				return err
			}

			data = menuData
		case reflect.TypeOf(system.SysApi{}):
			var apiData []system.SysApi
			err = global.OPS_DB.Model(table).Omit("id", "created_at", "updated_at", "deleted_at").Find(&apiData).Error
			if err != nil {
				return err
			}
			data = apiData
		default:
			return errors.New("unknown table type")
		}

		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}

		// 使用 filepath.Join 拼接路径
		filePath := filepath.Join(jsonDir, fmt.Sprintf("%s.json", table.TableName()))
		file, err := os.Create(filePath)
		if err != nil {
			color.Error.Println("Error creating file:", err)
			return err
		}

		_, err = file.Write(jsonData)
		if err != nil {
			file.Close() // 写入失败时立即关闭文件
			color.Error.Printf("导出%s失败: %v\n", table.TableName(), err)
			return err
		}
		file.Close() // 写入成功后关闭文件
		color.Info.Printf("导出%s成功\n", table.TableName())
	}

	return nil
}
