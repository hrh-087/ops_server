package clent_json

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
)

type RankList struct {
	Id       int    `json:"_id"`
	RankName string `json:"_name"`
	RankType int    `json:"_type"`
	RankDes  string `json:"_des"`
}

func ReadRankListJson(ctx context.Context, fileName string) (result map[string]RankList, err error) {
	configDir, err := GetConfigDir(ctx)
	if err != nil {
		return
	}

	file, err := os.Open(filepath.Join(configDir, fileName))
	if err != nil {
		return
	}
	defer file.Close()

	// 创建 JSON 解析器
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
