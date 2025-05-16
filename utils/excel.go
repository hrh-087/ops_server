package utils

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"ops-server/global"
	"reflect"
	"strconv"
)

func CreateSheet(f *excelize.File, sheetName string, annotationHeaders map[string]string, items []interface{}) error {
	sheetIndex, err := f.NewSheet(sheetName)
	if err != nil {
		global.OPS_LOG.Error("创建工作表失败", zap.Error(err))
		return err
	}

	// **提取字段名**
	headers, err := getHeaders(items[0], annotationHeaders)
	if err != nil {
		global.OPS_LOG.Error("提取字段名失败", zap.Error(err))
		return err
	}

	// **写入表头**
	for i, header := range headers {
		column := string(rune('A' + i))
		for j, headerName := range header {
			err = f.SetCellValue(sheetName, column+strconv.Itoa(j+1), headerName)
			if err != nil {
				global.OPS_LOG.Error("写入表头失败", zap.Error(err))
				return err
			}
		}
	}

	// **写入数据**
	for row, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		for col, header := range headers {
			column := string(rune('A' + col))
			cell := fmt.Sprintf("%s%d", column, row+4)
			err = f.SetCellValue(sheetName, cell, itemMap[header[0]])
			if err != nil {
				global.OPS_LOG.Error("写入数据失败", zap.Error(err))
				return err
			}
		}
	}

	// **设置默认 Sheet**
	f.SetActiveSheet(sheetIndex)

	return err
}

// getHeaders
// **📌 提取字段名**
func getHeaders(item interface{}, annotationHeaders map[string]string) ([][]string, error) {

	data := make([][]string, 0)

	if itemMap, ok := item.(map[string]interface{}); ok {
		for key := range itemMap {
			headers := make([]string, 0, len(itemMap))
			annotationName, ok := annotationHeaders[key]
			if !ok {
				return nil, errors.New("not found")
			}
			headers = append(headers, key, annotationName, reflect.TypeOf(itemMap[key]).String())
			data = append(data, headers)
		}
	}
	return data, nil

}
